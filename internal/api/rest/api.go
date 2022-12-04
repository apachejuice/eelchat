package rest

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/apachejuice/eelchat/internal/api/spec"
	"github.com/apachejuice/eelchat/internal/config"
	. "github.com/apachejuice/eelchat/internal/config/keys"
	"github.com/apachejuice/eelchat/internal/logs"
	"github.com/ogen-go/ogen/middleware"
)

// Main API object, containing handlers to be connected.
// Does some logging, maybe rate limiting in the future.
type API struct {
	logger      logs.Logger // the logger
	middlewares []spec.Middleware
	ghost       *apiHandler
	actions     ActionIDSource
	server      *spec.Server
	tlsConf     *tls.Config
	tlsOnce     sync.Once
	triedTls    bool
}

func (a *API) logMiddleware(req middleware.Request, next middleware.Next) (middleware.Response, error) {
	actionId := a.actions.Append(req.OperationID)
	a.ghost.lastActionId = actionId

	log := a.logger.New("ip", req.Raw.RemoteAddr, "operation", req.OperationID, "actionId", actionId)
	log.Debug("Incoming request")

	resp, err := next(req)
	if err != nil {
		log.Error("Error handling request", "error", err.Error())
	} else {
		var fields []any
		if tresp, ok := resp.Type.(interface{ GetStatusCode() int }); ok {
			status := tresp.GetStatusCode()
			statusStr := strconv.Itoa(status)
			fields = append(fields, "status")
			fields = append(fields, statusStr)

			if status >= 500 {
				log.Error("Failure handling request", "status", statusStr)
				return resp, err
			}
		}

		log.Info("Request handled successfully", fields...)
	}

	return resp, err
}

// Creates a new API.
func NewAPI(
	createUserFunc CreateUserFunc,
) *API {
	a := &API{logger: logs.NewLogger("api")}
	a.AddMiddleware(a.logMiddleware)
	// add the handlers
	a.ghost = &apiHandler{
		createUserFunc: createUserFunc,
	}

	s, err := spec.NewServer(a.ghost, spec.WithMiddleware(a.middlewares...))
	if err != nil {
		a.logger.Fatal("Unable to create server instance", "error", err.Error())
		return nil
	}

	a.server = s
	return a
}

func (a *API) ConfigureTLS() {
	a.tlsOnce.Do(func() {
		a.logger.Info("Configuring TLS without host information, patched later.")
		a.triedTls = true

		certPath, keyPath := ConfigKeyTlsCert.Get(), ConfigKeyTlsKey.Get()
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			a.logger.Error("Unable to load TLS certificate", "error", err.Error(), "certificate", certPath, "key", keyPath)
			return
		}

		cnf := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		a.tlsConf = cnf
	})
}

func (a *API) Run() {
	addr := ConfigKeyServeHost.Get()
	port := ConfigKeyServePort.Get()
	env := config.GetEnv()
	log := a.logger.New("env", env, "address", addr, "port", port)

	hs := &http.Server{
		Addr:      fmt.Sprintf("%s:%s", addr, port),
		Handler:   a.server,
		TLSConfig: a.tlsConf, // may be nil, as is also by default
	}

	if a.tlsConf == nil {
		if a.triedTls {
			a.logger.Warn("Attempted TLS configuration failed, continuing to serve HTTP")
		}

		if env == config.EnvProduction {
			log.Warn("Running without TLS in a production environment")
		}

		log.Info("Serving HTTP")
		hs.ListenAndServe()
	} else {
		a.tlsConf.ServerName = strings.Split(addr, ":")[0] // no port

		log.Info("Serving HTTPS")
		hs.ListenAndServeTLS("", "") // cert and key in hs.TLSConfig
	}
}

func (a *API) AddMiddleware(middleware spec.Middleware) {
	a.middlewares = append(a.middlewares, middleware)
}

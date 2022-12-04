package rest

import (
	"fmt"
	"log"
	"time"

	"github.com/teris-io/shortid"
)

// An action id identifies a specific request pipeline and helps with debugging.
type ActionID struct {
	OperationID   string
	Discriminator string // a short identifier
}

func (a ActionID) String() string {
	return fmt.Sprintf("%s:%s", a.OperationID, a.Discriminator)
}

// Generates an action id for debugging.
// It is composed of a format operationId-salt.
func generateActionId(operationId string) ActionID {
	end, err := shortid.Generate()
	if err != nil {
		log.Fatal(err)
	}

	return ActionID{
		OperationID:   operationId,
		Discriminator: end,
	}
}

// A struct to store information about action IDs.
type ActionIDSource struct {
	ids map[int64]ActionID

	StartTime int64
	EndTime   int64
}

func (a ActionIDSource) Get(t time.Time) (ActionID, error) {
	if val, ok := a.ids[t.UnixNano()]; ok {
		return val, nil
	}

	return ActionID{"", ""}, fmt.Errorf("no action happened at %s", t.String())
}

func (a *ActionIDSource) Append(operationId string) ActionID {
	actionId := generateActionId(operationId)
	nanos := time.Now().UTC().UnixNano()

	if a.ids == nil {
		a.ids = make(map[int64]ActionID)
		a.StartTime = nanos
	}

	a.EndTime = nanos
	a.ids[nanos] = actionId
	return actionId
}

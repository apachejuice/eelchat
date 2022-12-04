package api

// Component strings identifying the component that failed
const (
	ComponentAPI      string = "api"      // the api layer, before any data processing has occurred
	ComponentImpl     string = "impl"     // implementation of api calls
	ComponentDatabase string = "database" // db operations
)

package app_errors

type IdentifiableError struct {
	HTTPStatusCode int
	CustomMSG      string
}

package e

type ErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

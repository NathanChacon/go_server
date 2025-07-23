package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
)

type APIError struct {
	ErrorID     string `json:"errorId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

var (
	ErrNotAuthenticated = errors.New("NOT AUTHENTICATED")
	ErrNotAuthorized    = errors.New("NOT AUTHORIZED")
	ErrNotFound         = errors.New("NOT FOUND")
	ErrBadRequest       = errors.New("BAD REQUEST")
)

func sendJson(response http.ResponseWriter, apiError APIError) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(apiError.Status)
	json.NewEncoder(response).Encode(apiError)
}

func generateApiError(err error) APIError {
	var apiError APIError

	switch {
	case errors.Is(err, ErrNotAuthenticated):
		apiError = APIError{
			ErrorID:     "not_authenticated",
			Title:       "Oops, not authenticated",
			Description: "You need to authenticate to proceed.",
			Status:      http.StatusUnauthorized,
		}

	case errors.Is(err, ErrNotAuthorized):
		apiError = APIError{
			ErrorID:     "not_authorized",
			Title:       "Access Denied",
			Description: "You are not allowed to perform this action.",
			Status:      http.StatusForbidden,
		}

	case errors.Is(err, ErrNotFound):
		apiError = APIError{
			ErrorID:     "not_found",
			Title:       "Not Found",
			Description: "The requested resource was not found.",
			Status:      http.StatusNotFound,
		}

	case errors.Is(err, ErrBadRequest):
		apiError = APIError{
			ErrorID:     "bad_request",
			Title:       "Bad Request",
			Description: "wrong payload",
			Status:      http.StatusBadRequest,
		}

	default:
		apiError = APIError{
			ErrorID:     "internal_error",
			Title:       "Something went wrong",
			Description: "Please try again later.",
			Status:      http.StatusInternalServerError,
		}
	}

	return apiError
}

func HandleError(err error, response http.ResponseWriter) {
	apiError := generateApiError(err)

	sendJson(response, apiError)
}

func HandleErrorWithCustomDescription(err error, response http.ResponseWriter, customDescription string) {
	apiError := generateApiError(err)
	apiError.Description = customDescription
	sendJson(response, apiError)
}

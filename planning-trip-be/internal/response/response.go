package response

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// APIResponse wraps every JSON response that leaves the server.
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// APIError represents an error that should be returned to the client.
type APIError struct {
	Status  int
	Message string
	Err     error
}

// Error implements the error interface to allow wrapping.
func (e APIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Write writes a JSON response with the given status, data, and optional message.
func Write(w http.ResponseWriter, status int, data interface{}, message string) {
	resp := APIResponse{
		Success:   status >= http.StatusOK && status < http.StatusMultipleChoices,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("response encode error: %v", err)
	}
}

// WriteError writes an error response using the provided APIError.
func WriteError(w http.ResponseWriter, apiErr APIError) {
	if apiErr.Status == 0 {
		apiErr.Status = http.StatusInternalServerError
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(apiErr.Status)
	}
	if apiErr.Err != nil {
		log.Printf("api error: %v", apiErr.Err)
	}
	Write(w, apiErr.Status, nil, apiErr.Message)
}

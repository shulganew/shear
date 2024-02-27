// Package entities represent a Model layer of app model. It has main shorterer entities and DTO objects for JSON marshal.
package entities

// DTO object for batch request.
type BatchRequest struct {
	SessionID string `json:"correlation_id"`
	Origin    string `json:"original_url"`
}

// DTO object for batch response.
type BatchResponse struct {
	SessionID string `json:"correlation_id"`
	Answer    string `json:"short_url"`
}

package entities

type BatchRequest struct {
	SessionID string `json:"correlation_id"`
	Origin    string `json:"original_url"`
}

type BatchResponse struct {
	SessionID string `json:"correlation_id"`
	Answer    string `json:"short_url"`
}

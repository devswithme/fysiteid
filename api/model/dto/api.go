package dto

type ApiMeta struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type APIResponse[T any] struct {
	Data T       `json:"data"`
	Meta ApiMeta `json:"meta"`
}

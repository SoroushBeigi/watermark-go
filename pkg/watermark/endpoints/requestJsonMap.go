package endpoints

import "github.com/SoroushBeigi/watermark-go/internal"

type GetRequest struct {
	Filters []internal.Filter `json:"filters,omitempty"`
}

type GetResponse struct {
	Documents []internal.Document `json:"documents"`
	Error     string              `json:"err,omitempty"`
}

type StatusRequest struct {
	TicketId string `json:"ticketId"`
}

type StatusResponse struct {
	Status internal.Status `json:"status"`
	Error  string          `json:"err,omitempty"`
}

type WatermarkRequest struct {
	TicketId string `json:"ticketId"`
	Mark     string `json:"mark"`
}

type WatermarkResponse struct {
	Code  int    `json:"code"`
	Error string `json:"err"`
}

type AddDocumentRequest struct {
	Document *internal.Document `json:"document"`
}

type AddDocumentResponse struct {
	TicketId string `json:"ticketId"`
	Error    string `json:"err,omitempty"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Code  int    `json:"status"`
	Error string `json:"err,omitempty"`
}

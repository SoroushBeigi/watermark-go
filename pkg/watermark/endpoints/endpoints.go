package endpoints

import (
	"context"
	"errors"
	"os"

	"github.com/SoroushBeigi/watermark-go.git/internal"
	"github.com/SoroushBeigi/watermark-go.git/pkg/watermark"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	GetEndpoint           endpoint.Endpoint
	StatusEndpoint        endpoint.Endpoint
	WatermarkEndpoint     endpoint.Endpoint
	AddDocumentEndpoint   endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func NewEndpointSet(service watermark.Service) Set {
	return Set{
		GetEndpoint:           MakeGetEndpoint(service),
		StatusEndpoint:        MakeStatusEndpoint(service),
		WatermarkEndpoint:     MakeWatermarkEndpoint(service),
		AddDocumentEndpoint:   MakeAddDocumentEndpoint(service),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(service),
	}
}

func MakeGetEndpoint(service watermark.Service) endpoint.Endpoint {
	return func(context context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		docs, err := service.Get(context, req.Filters...)
		if err != nil {
			return GetResponse{Documents: docs, Error: err.Error()}, nil
		}
		return GetResponse{docs, ""}, nil
	}
}

func MakeAddDocumentEndpoint(service watermark.Service) endpoint.Endpoint {
	return func(context context.Context, request interface{}) (interface{}, error) {
		req := request.(AddDocumentRequest)
		ticketId, err := service.AddDocument(context, req.Document)
		if err != nil {
			return AddDocumentResponse{TicketId: ticketId, Error: err.Error()}, nil
		}
		return AddDocumentResponse{TicketId: ticketId, Error: ""}, nil
	}
}

func MakeStatusEndpoint(service watermark.Service) endpoint.Endpoint {
	return func(context context.Context, request interface{}) (interface{}, error) {
		req := request.(StatusRequest)
		status, err := service.Status(context, req.TicketId)
		if err != nil {
			return StatusResponse{Status: status, Error: err.Error()}, nil
		}
		return StatusResponse{Status: status, Error: ""}, nil
	}
}

func MakeServiceStatusEndpoint(service watermark.Service) endpoint.Endpoint {
	return func(context context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := service.ServiceStatus(context)
		if err != nil {
			return ServiceStatusResponse{Code: code, Error: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Error: ""}, nil
	}
}

func MakeWatermarkEndpoint(service watermark.Service) endpoint.Endpoint {
	return func(context context.Context, request interface{}) (interface{}, error) {
		req := request.(WatermarkRequest)
		code, err := service.Watermark(context, req.TicketId, req.Mark)
		if err != nil {
			return WatermarkResponse{Code: code, Error: err.Error()}, nil
		}
		return WatermarkResponse{Code: code, Error: ""}, nil

	}
}

func (s *Set) Get(context context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	resp, err := s.GetEndpoint(context, GetRequest{Filters: filters})
	if err != nil {
		return []internal.Document{}, err
	}
	getResp := resp.(GetResponse)
	if getResp.Error != "" {
		return []internal.Document{}, errors.New(getResp.Error)
	}
	return getResp.Documents, nil
}

func (s *Set) Status(context context.Context, ticketId string) ([]internal.Status, error) {
	resp, err := s.GetEndpoint(context, StatusRequest{TicketId: ticketId})
	if err != nil {
		return internal.Failed, err
	}
	statusResp := resp.(StatusResponse)
	if statusResp.Error != "" {
		return internal.Failed{}, errors.New(statusResp.Error)
	}
	return statusResp.Status, nil
}

func (s *Set) Watermark(context context.Context, ticketId, mark string) (int, error) {
	resp, err := s.WatermarkEndpoint(context, WatermarkRequest{TicketId: ticketId, Mark: mark})
	watermarkResp := resp.(WatermarkResponse)
	if err != nil {
		return watermarkResp.Code, err
	}
	if watermarkResp.Err != "" {
		return watermarkResp.Code, errors.New(watermarkResp.Error)
	}
	return watermarkResp.Code, nil
}

func (s *Set) AddDocument(context context.Context, doc *internal.Document) (string, error) {
	resp, err := s.AddDocumentEndpoint(context, AddDocumentRequest{Document: doc})
	if err != nil {
		return "", err
	}
	addResp := resp.(AddDocumentResponse)
	if addResp.Error != "" {
		return "", errors.New(addResp.Error)
	}
	return addResp.TicketId, nil
}

func (s *Set) ServiceStatus(context context.Context) (int, error) {
	resp, err := s.ServiceStatusEndpoint(context, ServiceStatusRequest{})
	serviceStatusResp := resp.(ServiceStatusResponse)
	if err != nil {
		return serviceStatusResp.Code, err
	}
	if serviceStatusResp.Error != "" {
		return serviceStatusResp.Code, errors.New(serviceStatusResp.Error)
	}
	return serviceStatusResp.Code, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

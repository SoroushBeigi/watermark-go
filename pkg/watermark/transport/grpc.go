package transport

import (
    "context"
	"github.com/velotiotech/watermark-service/api/v1/pb/watermark"
    grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get           grpctransport.Handler
	status        grpctransport.Handler
	addDocument   grpctransport.Handler
	watermark     grpctransport.Handler
	serviceStatus grpctransport.Handler
}

fucn newGRPCServer(ep endpoints.Set) watermark.WatermarkServer{
	return &grpcServer{
		get: grpctransport.NewServer(
			ep.GetEndpoint,
			decodeGRPCGetRequest,
			decodeGRPCGetResponse,
		),
		status: grpctransport.NewServer(
			ep.StatusEndpoint,
			decodeGRPCStatusRequest,
            decodeGRPCStatusResponse,
		),
		addDocument: grpctransport.NewServer(
            ep.AddDocumentEndpoint,
            decodeGRPCAddDocumentRequest,
            decodeGRPCAddDocumentResponse,
        ),
        watermark: grpctransport.NewServer(
            ep.WatermarkEndpoint,
            decodeGRPCWatermarkRequest,
            decodeGRPCWatermarkResponse,
        ),
        serviceStatus: grpctransport.NewServer(
            ep.ServiceStatusEndpoint,
            decodeGRPCServiceStatusRequest,
            decodeGRPCServiceStatusResponse,
        ),
	}
}

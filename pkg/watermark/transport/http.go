package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SoroushBeigi/watermark-go.git/pkg/watermark/endpoints"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/get", httptransport.NewServer(
		ep.GetEndpoint,
		decodeHTTPGetRequest,
		encodeResponse,
	))
	mux.Handle("/status", httptransport.NewServer(
		ep.StatusEndpoint,
		decodeHTTPStatusRequest,
		encodeRespose,
	))
	mux.Handle("/watermark", httptransport.NewServer(
		ep.WatermarkEndpoint,
		decodeHTTPWatermarkRequest,
		encodeResponse,
	))
	mux.Handle("/addDocument", httptransport.NewServer(
		ep.AddDocumentEndpoint,
		decodeHTTPAddDocumentRequest,
		encodeRespose,
	))

	mux.Handle("/healthz", httptransport.NewServer(
		ep.ServiceStatusEndpoint,
		decodeHTTPServiceStatusRequest,
		encodeRespose,
	))

	return mux
}

func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.GetRequest
	if r.ContentLength == 0 {
		logger.Log("Get request with no body")
		return req, nil
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.StatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPWatermarkRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.WatermarkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPAddDocumentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.AddDocumentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeHTTPServiceStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ServiceStatusRequest
	return req, nil
}

func encodeResponse(context context.Context, writer http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		encodeError(context, e, writer)
		return nil
	}
	return json.NewEncoder(writer).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case util.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case util.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

package proxy

import (
	"context"
	"net/http"

	"simple-http-proxy/internal/domain/proxy"
	"simple-http-proxy/pkg/log"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *Service) DoClientRequest(ctx context.Context, clientReqDTO proxy.Request) (targetResponse proxy.Response, err error) {
	requestID := middleware.GetReqID(ctx)
	logger := log.LoggerFromContext(ctx).With().Str("DoClientRequest", requestID).Logger()

	clientReq := &http.Request{
		Method: clientReqDTO.Method,
		URL:    clientReqDTO.ParsedURL,
		Header: clientReqDTO.ParsedHeaders,
	}

	response, err := http.DefaultClient.Do(clientReq)
	if err != nil {
		logger.Error().Err(err).Msg("failed to do request")
		return
	}

	targetResponse = buildTargetResponse(requestID, response)

	requestResponsePair := &proxy.Entity{
		ID:       requestID,
		Request:  clientReqDTO,
		Response: targetResponse,
	}

	s.cache.Set(ctx, requestID, *requestResponsePair)

	return
}

func buildTargetResponse(requestID string, response *http.Response) (targetResponse proxy.Response) {
	headers := make(map[string]string)
	for header, values := range response.Header {
		headers[header] = values[0]
	}

	length := int64(len(requestID) + len(response.Status) + len(headers))

	targetResponse = proxy.Response{
		RequestID:  requestID,
		StatusCode: response.StatusCode,
		Headers:    headers,
		Length:     length,
	}

	return
}

package http

import (
	"net/http"

	"simple-http-proxy/internal/domain/proxy"
	proxyService "simple-http-proxy/internal/service/proxy"
	"simple-http-proxy/pkg/server/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ProxyHandler struct {
	ProxyService *proxyService.Service
}

func NewProxyHandler(ps *proxyService.Service) *ProxyHandler {
	return &ProxyHandler{
		ProxyService: ps,
	}
}

func (h *ProxyHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.HandleFunc("/", h.handleClientRequest)

	return r
}

// @Summary Proxy client request
// @Description Proxy client request
// @Tags Proxy
// @Accept json
// @Produce json
// @Param request body proxy.Request true "Request"
// @Success 200 {object} proxy.Response
// @Failure 500 {string} string
// @Router /proxy [post]
func (h *ProxyHandler) handleClientRequest(w http.ResponseWriter, r *http.Request) {
	var clientReqDTO proxy.Request
	err := render.Bind(r, &clientReqDTO)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	targetResponse, err := h.ProxyService.DoClientRequest(r.Context(), clientReqDTO)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, targetResponse)
}

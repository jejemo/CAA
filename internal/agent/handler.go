package agent

import (
	"go-caa/internal/server/resp"
	"net/http"

	"github.com/rs/zerolog/log"
)

type handler struct {
	svc *Service
}

func NewHandler(svc *Service) handler {
	return handler{
		svc: svc,
	}
}

func (h *handler) GetAgents(w http.ResponseWriter, r *http.Request) {
	h.svc.GetAvailableAgents()

	resp.WriteJSON(w, http.StatusOK, "OK")
}

func (h *handler) AssignAgent(w http.ResponseWriter, r *http.Request) {
	err := h.svc.AssignAgentToRoom(265117341, 169606)
	if err != nil {
		log.Info().Msgf("err: %v", err)
	}

	resp.WriteJSON(w, http.StatusOK, "OK")
}

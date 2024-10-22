package caa

import (
	"encoding/json"
	"fmt"
	"go-caa/internal/server/resp"
	"io"
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

func (h *handler) WebhookCAA(w http.ResponseWriter, r *http.Request) {
	var caa QismoCAAWebhook
	if err := h.parseJSONBody(r, &caa); err != nil {
			h.handleError(w, "Unable to parse request body", err, http.StatusBadRequest)
			return
	}

	if err := h.svc.ProcessAllocateAgent(caa); err != nil {
			h.handleError(w, "Unable to Process Allocate Agent", err, http.StatusBadRequest)
			return
	}

	h.respondOK(w)
}

func (h *handler) WebhookMarkAsResolve(w http.ResponseWriter, r *http.Request) {
	var caa QismoMarkAsResolvedWebhook
	if err := h.parseJSONBody(r, &caa); err != nil {
			h.handleError(w, "Unable to parse request body", err, http.StatusBadRequest)
			return
	}

	if err := h.svc.ResolveRoom(caa); err != nil {
			h.handleError(w, "Unable to Resolve Room", err, http.StatusBadRequest)
			return
	}

	if err := h.svc.AllocateAgentsToRooms(); err != nil {
			h.handleError(w, "Unable to Allocate Agents to Rooms", err, http.StatusBadRequest)
			return
	}

	h.respondOK(w)
}

// Helper methods

func (h *handler) parseJSONBody(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
			return fmt.Errorf("reading request body: %w", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("parsing JSON body: %w", err)
	}

	return nil
}

func (h *handler) handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	log.Printf("%s: %v", message, err) // Log the error
	http.Error(w, message, statusCode)
}

func (h *handler) respondOK(w http.ResponseWriter) {
	resp.WriteJSON(w, http.StatusOK, "OK")
}

func (h *handler) GetPendingQueue(w http.ResponseWriter, r *http.Request) {
	queues, err := h.svc.repo.GetPendingQueues()
	if err != nil {
		log.Info().Msgf("err: %v", err)
	}

	for _, queue := range queues {
		log.Info().Msgf("\nqueue %v, %v, %v", queue.RoomID, queue.AgentID, queue.CreateAt)
	}

	resp.WriteJSON(w, http.StatusOK, "OK")
}

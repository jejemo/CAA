package server

import (
	"encoding/json"
	"go-caa/internal/agent"
	"go-caa/internal/api/client"
	"go-caa/internal/caa"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-retryablehttp"
)

func newHTTPClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Logger = nil
	retryClient.HTTPClient.Timeout = 20 * time.Second
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	return retryClient.StandardClient()
}

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// qismo
	qismo := client.NewQismoClient(os.Getenv("OMNICHANNEL_URL"), os.Getenv("OMNICHANNEL_APP_ID"), os.Getenv("OMNICHANNEL_SECRET_KEY"), newHTTPClient())

	// agent
	agentRepo := agent.NewRepository(s.db, qismo)
	agentSvc := agent.NewService(agentRepo)
	agentHandler := agent.NewHandler(agentSvc)

	// caa
	caaRepo := caa.NewQueueRepository(s.db)
	caaSvc := caa.NewService(caaRepo, agentSvc)
	caaHandler := caa.NewHandler(caaSvc)

	r.HandleFunc("/agent", agentHandler.GetAgents).Methods("GET")
	r.HandleFunc("/agentassign", agentHandler.AssignAgent).Methods("GET")
	r.HandleFunc("/wh/caa", caaHandler.WebhookCAA).Methods("POST")
	// r.HandleFunc("/wh/caa", caaHandler.).Methods("GET")
	r.HandleFunc("/", s.HelloWorldHandler)
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

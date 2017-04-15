package server

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"procurementlistenerservice/model"
	"strconv"
)

// Server is the main struct for the backend service.
type Server struct {
	port    int
	service model.PartnerBackendService
}

// CreateServer creates a new Server instance for serving incoming requests at the given port.
func CreateServer(serverPort int, service model.PartnerBackendService) (*Server, error) {
	return &Server{
		port:    serverPort,
		service: service,
	}, nil
}

// Start initiates the http service and starts listening to incoming connections.
func (s *Server) Start() {
	router := mux.NewRouter()

	http.Handle("/", handlers.LoggingHandler(os.Stderr, router))
	s.registerDispatchers(router)
	log.Printf("Starting server on port '%d'\n", s.port)
	http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
}

func (s *Server) registerDispatchers(router *mux.Router) {
	log.Print("Registering dispatcher at /entitlementEvents")
	router.HandleFunc("/entitlementEvents", s.onEntitlementEvent).Methods("POST")
}

func (s *Server) onEntitlementEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read body: '%v'\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var notification model.EntitlementEvent
	err = json.Unmarshal(body, &notification)
	if err != nil {
		log.Printf("Unable to parse body: '%v'\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = model.ValidateEntitlementEvent(notification)
	if err != nil {
		log.Printf("Invalid entitlement event received: '%v'\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := s.service.OnEntitlementEvent(notification)
	if err != nil {
		log.Printf("Error handling entitlement event: '%v'\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling response: '%v'\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch response.Status {
	case model.RESPONSESTATUS_INVALIDREQUEST:
		w.WriteHeader(http.StatusBadRequest)
		break
	case model.RESPONSESTATUS_ACCEPTED:
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes)
		break
	case model.RESPONSESTATUS_ASYNC:
		w.WriteHeader(http.StatusAccepted)
		break
	case model.RESPONSESTATUS_REJECTED:
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBytes)
		break

	default:
		log.Printf("Unknown response status: '%d' response: '%v'\n", response.Status, response)

	}
}

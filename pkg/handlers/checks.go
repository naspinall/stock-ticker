package handlers

import (
	"net/http"
)

// Kubernetes liveness handler, used to determine if the service is alive
type LivenessHandler struct{}

func NewLivenessHandler() *LivenessHandler {
	return &LivenessHandler{}
}

func (lh *LivenessHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Kubernetes readiness handler, used to determine if the service is ready to receive traffic.
type ReadinessHandler struct{}

func NewReadinessHandler() *ReadinessHandler {
	return &ReadinessHandler{}
}

func (lh *ReadinessHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

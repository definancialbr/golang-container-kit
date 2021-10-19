package probes

import "net/http"

type ProbeService interface {
	LivenessHandler() http.HandlerFunc
	ReadinessHandler() http.HandlerFunc
}

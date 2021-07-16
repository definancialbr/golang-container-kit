package probes

import "net/http"

type ProbeService interface {
	LivenessHandler() http.Handler
	ReadinessHandler() http.Handler
}

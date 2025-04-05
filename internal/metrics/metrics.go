package metrics

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	register sync.Once

	websocketConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "websocket_connections",
			Help: "Number of active WebSocket connections",
		},
		[]string{"endpoint"},
	)

	MessagesReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_received_total",
			Help: "Total number of messages received via WebSocket",
		},
		[]string{"endpoint"},
	)

	MessagesPublished = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_messages_published_total",
			Help: "Total number of messages successfully published",
		},
		[]string{"endpoint"},
	)

	PublishErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "websocket_publish_errors_total",
			Help: "Total number of publish errors",
		},
		[]string{"endpoint"},
	)

	WebsocketLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "websocket_message_latency_seconds",
			Help:    "Latency of WebSocket message processing",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"endpoint"},
	)
)

func InitMetrics() {
	register.Do(func() {
		prometheus.MustRegister(websocketConnections)
		prometheus.MustRegister(MessagesReceived)
		prometheus.MustRegister(MessagesPublished)
		prometheus.MustRegister(PublishErrors)
		prometheus.MustRegister(WebsocketLatency)
	})
}

func PrometheusMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path

		isWebSocket := websocket.IsWebSocketUpgrade(r)
		if isWebSocket {
			websocketConnections.WithLabelValues(endpoint).Inc()
			next.ServeHTTP(w, r)
			websocketConnections.WithLabelValues(endpoint).Dec()
			return
		}

		next.ServeHTTP(w, r)
	}
}

func PromHandler() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}

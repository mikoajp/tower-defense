package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	WsConnections      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "td_ws_connections", Help: "Number of active WS connections"})
	TicksTotal         = prometheus.NewCounter(prometheus.CounterOpts{Name: "td_engine_ticks_total", Help: "Total engine ticks"})
	EngineEnemies      = prometheus.NewGauge(prometheus.GaugeOpts{Name: "td_engine_enemies", Help: "Current number of enemies"})
	EngineProjectiles  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "td_engine_projectiles", Help: "Current number of projectiles"})
	EngineTowers       = prometheus.NewGauge(prometheus.GaugeOpts{Name: "td_engine_towers", Help: "Current number of towers"})
	EngineTickSeconds  = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "td_engine_tick_seconds",
		Help:    "Engine tick delta time in seconds",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 12),
	})
)

func init() {
	prometheus.MustRegister(WsConnections, TicksTotal, EngineEnemies, EngineProjectiles, EngineTowers, EngineTickSeconds)
}

func MountMetrics(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

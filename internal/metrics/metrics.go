package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RPCCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "fractal_rpc_calls_total",
		Help: "Total RPC calls by method",
	}, []string{"method"})

	BlockStoreLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "fractal_block_store_seconds",
		Help: "Time to store a block",
	})
)

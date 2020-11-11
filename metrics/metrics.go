package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	goroutinesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "website_ssl_earliest_cert_expiry",
			Help: "websie ssl cert expiry timestamp",
		},
		[]string{"website"},
	)
)

func init() {
	prometheus.MustRegister(goroutinesGauge)
}

func Observer(url string,timestamp int64) {
	for {
		goroutinesGauge.With(prometheus.Labels{"website": url}).Set(float64(timestamp))
		time.Sleep(1 * time.Second)
	}
}


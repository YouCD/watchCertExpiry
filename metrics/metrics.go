package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"time"
)

var (
	goroutinesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "website_ssl_earliest_cert_expiry",
			Help: "website ssl cert expiry timestamp",
		},
		[]string{"website"},
	)
)

func init() {
	prometheus.MustRegister(goroutinesGauge)
}
func GetExpiryTimestamp(url string) (timestamp int64) {
	//var tr = &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	client := &http.Client{}
	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer response.Body.Close()
	if response.TLS != nil && len(response.TLS.PeerCertificates) != 0 {
		value := response.TLS.PeerCertificates[0].NotAfter
		timestamp = value.Unix()
		return timestamp
	} else {
		log.Println("Site does not use HTTPS certificates.")
		return 0
	}

}
func Observer(url string,) {
	for {
		timestamp:=GetExpiryTimestamp(url)
		goroutinesGauge.With(prometheus.Labels{"website": url}).Set(float64(timestamp))
		time.Sleep(10 * time.Second)
	}
}


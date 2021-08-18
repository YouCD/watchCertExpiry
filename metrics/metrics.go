package metrics

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	webCertGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "website_cert_expiry_timestamp",
			Help: "website cert expiry timestamp",
		},
		[]string{"website"},
	)

	websiteIsDownGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "website_is_down",
			Help: "website is down",
		},
		[]string{"website"},
	)

	websiteCertConfigrationFailed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "website_configuration_failed",
			Help: "website configuration failed",
		},
		[]string{"website"},
	)
	WesiteIsDownErr                 = errors.New("website is down")
	WesiteCertConfigrationFailedErr = errors.New("website cert configration failed")
)

func init() {
	prometheus.MustRegister(webCertGauge)
	prometheus.MustRegister(websiteIsDownGauge)
	prometheus.MustRegister(websiteCertConfigrationFailed)
}
func GetExpiryTimestamp(url string) (timestamp int64, err error) {
	//var tr = &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	client := &http.Client{}
	response, err := client.Get(url)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "connection refused"):
			return 0, WesiteIsDownErr
		case strings.Contains(err.Error(), "certificate signed by unknown authority"):
			return 0, WesiteCertConfigrationFailedErr
		}
		log.Println(err)
		return 0, err
	}
	defer response.Body.Close()
	if response.TLS != nil && len(response.TLS.PeerCertificates) != 0 {
		value := response.TLS.PeerCertificates[0].NotAfter
		timestamp = value.Unix()
		return timestamp, nil
	}
	log.Println("Site does not use HTTPS certificates.")
	return 0, err

}

func Observer(url string) {
	for {
		timestamp, err := GetExpiryTimestamp(url)
		if err != nil && errors.Is(err, WesiteIsDownErr) {
			websiteIsDownGauge.WithLabelValues(url).Set(float64(0))
			log.Printf("website %s is down", url)
		} else if errors.Is(err, WesiteCertConfigrationFailedErr) {
			websiteCertConfigrationFailed.WithLabelValues(url).Set(float64(0))
			log.Printf("website %s configration failed", url)
		} else {
			websiteIsDownGauge.WithLabelValues(url).Set(float64(1))
			websiteCertConfigrationFailed.WithLabelValues(url).Set(float64(1))
		}

		webCertGauge.With(prometheus.Labels{"website": url}).Set(float64(timestamp))
		time.Sleep(10 * time.Second)
	}
}

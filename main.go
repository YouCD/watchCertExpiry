package main

import (
	"crypto/tls"
	"errors"
	"github.com/YouCD/watchCertExpiry/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"os"
	"strings"

	"net/http"
)

func GetUrlFromEnv() (urls []string,err error) {
	str := os.Getenv("CERT_MONITOR_URL_LIST")
	if len(str) > 0 {
		urlList := strings.Split(str, ",")
		for _, v := range urlList {
			v = "https://" + v
			urls = append(urls, v)
		}
		return
	} else {

		return nil,errors.New("You must set the CERT_MONITOR_URL_LIST environment variable\nExample: \n    export CERT_MONITOR_URL_LIST=github.com")
	}
}

func GetExpiryTimestamp(url string) (timestamp int64) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return 0
	}
	if response.TLS != nil && len(response.TLS.PeerCertificates) != 0 {
		value := response.TLS.PeerCertificates[0].NotAfter
		timestamp = value.Unix()

		return timestamp
	} else {
		log.Println("Site does not use HTTPS certificates.")
		return 0
	}

}

func main() {
	urls ,err:= GetUrlFromEnv()
	if err !=nil{
		log.Panic(err)
	}
	for _,url:=range urls{
		timestamp:=GetExpiryTimestamp(url)
		log.Printf("monitor site %s",url)
		go metrics.Observer(url, timestamp)
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"errors"
	"fmt"
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



func main() {
	urls ,err:= GetUrlFromEnv()
	if err !=nil{
		fmt.Println(err)
		os.Exit(1)
	}
	for _,url:=range urls{
		log.Printf("monitor site %s",url)
		go metrics.Observer(url)
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

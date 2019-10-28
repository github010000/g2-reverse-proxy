package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/cnpst/g2-reverse-proxy/common/config"
	"github.com/cnpst/g2-reverse-proxy/common/log"
	"github.com/pkg/errors"
)

func setUp() error {
	if err := config.LoadConfig(); err != nil {
		return errors.Wrap(err, "failed to load configs")
	}

	if err := log.SetUp(); err != nil {
		return errors.Wrap(err, "failed to initialize log")
	}

	return nil
}

// Get the port to listen on
func getListenAddress() string {
	port := config.ListenPort()
	return ":" + port
}

// Get the url for a given proxy condition
func getProxyUrl() string {
	proxyUrl := config.OpsgenieURL()
	return proxyUrl
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	//Input request body Logging
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Errorf("Error reading body:[", err, "]")
			panic(err)
		}
		logBody := ioutil.NopCloser(bytes.NewBuffer(body))
		runBody := ioutil.NopCloser(bytes.NewBuffer(body))
		logBodyBuf, err := ioutil.ReadAll(logBody)
		log.Info("Input Request Body:[", string(logBodyBuf), "]")
		req.Body = runBody
	}
	// proxy url
	proxyUrl := getProxyUrl()
	log.Info("Proxy to URL:[", proxyUrl, "]")
	// parse the url
	url, _ := url.Parse(proxyUrl)
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host
	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func main() {
	if err := setUp(); err != nil {
		fmt.Printf("failed to set up, %s\n", err.Error())
		os.Exit(1)
		return
	}
	log.Info("g2 reverse proxy started")
	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		log.Errorf("failed to g2 reverse proxy, %s", err.Error())
	}
	log.Info("g2 reverse proxy finished")
}

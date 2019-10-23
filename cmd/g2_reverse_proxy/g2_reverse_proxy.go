package main

import (
		"bytes"
		"encoding/json"
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
	url := config.OpsgenieURL()
	return url
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	log.Info("Input Request Header Host: %s\n", req.Header.Get("Host"))
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host
	log.Info("Proxy Request Header Host: %s\n", req.Header.Get("Host"))

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}


func requestBodyToByte(request *http.Request) []byte {
	// Read body to buffer
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Errorf("Error reading body: %v", err)
		panic(err)
	}
	return body
}
// Get a json decoder for a given requests body
func requestBodyDecoder(request *http.Request) *json.Decoder {
	// Read body to buffer
	body:= requestBodyToByte(request)
	// Because go lang is a pain in the ass if you read the body then any susequent calls
	// are unable to read the body again....
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	log.Info("Input Request Body: %s\n", string(requestBodyToByte(req)))
	url := getProxyUrl()
	log.Info("Redirecting to URL: %s\n", url)
	serveReverseProxy(url, res, req)
	log.Info("Proxy Request Body: %s\n", string(requestBodyToByte(req)))
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

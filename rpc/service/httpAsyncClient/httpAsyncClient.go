package httpAsyncClient

import (
	"github.com/WeBankBlockchain/WeCross-Go-SDK/common"
	"github.com/WeBankBlockchain/WeCross-Go-SDK/rpc/service/connections"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	HTTP_CLIENT_TIME_OUT = 100000 * time.Millisecond
	MAX_HOLD_CONNECTIONS = 100
)

type AsyncHttpClient struct {
	httpClient *http.Client
}

func NewAsyncHttpClient(conn *connections.Connection) (*AsyncHttpClient, *common.WeCrossSDKError) {
	dialer := &net.Dialer{
		Timeout:   HTTP_CLIENT_TIME_OUT,
		KeepAlive: 15 * time.Second,
	}

	var transport *http.Transport

	if conn.GetSslSwitch() != common.SSL_OFF {
		// for ssl
		caCert, err := os.ReadFile(conn.GetCaCert())
		if err != nil {
			return nil, common.NewWeCrossSDKFromString(common.INTERNAL_ERROR, "Init http client error: "+err.Error())
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)
		clientCert, err := tls.LoadX509KeyPair(conn.GetSslCert(), conn.GetSslKey())
		if err != nil {
			return nil, common.NewWeCrossSDKFromString(common.INTERNAL_ERROR, "Init http client error: "+err.Error())
		}
		tlsConfig := &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{clientCert},
		}
		if conn.GetSslSwitch() == common.SSL_ON_CLIENT_AUTH {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}

		transport = &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			TLSClientConfig:       tlsConfig,
			TLSHandshakeTimeout:   HTTP_CLIENT_TIME_OUT,
			DisableKeepAlives:     false,
			ResponseHeaderTimeout: HTTP_CLIENT_TIME_OUT,
			MaxIdleConns:          MAX_HOLD_CONNECTIONS,
			IdleConnTimeout:       HTTP_CLIENT_TIME_OUT,
		}
	} else {
		transport = &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			TLSHandshakeTimeout:   HTTP_CLIENT_TIME_OUT,
			DisableKeepAlives:     false,
			ResponseHeaderTimeout: HTTP_CLIENT_TIME_OUT,
			MaxIdleConns:          MAX_HOLD_CONNECTIONS,
			IdleConnTimeout:       HTTP_CLIENT_TIME_OUT,
			ForceAttemptHTTP2:     true,
		}
	}

	httpClient := &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       HTTP_CLIENT_TIME_OUT,
	}
	return &AsyncHttpClient{httpClient: httpClient}, nil
}

func (ac *AsyncHttpClient) Prepare(method, url string, body []byte) (*http.Request, *common.WeCrossSDKError) {
	body_reader := bytes.NewReader(body)
	request, err := http.NewRequest(method, url, body_reader)
	if err != nil {
		return nil, common.NewWeCrossSDKFromString(common.RPC_ERROR, "fail in AsyncHttpClient.Prepare:"+err.Error())
	}
	return request, nil
}

func (ac *AsyncHttpClient) SendRequest(request *http.Request) (*http.Response, *common.WeCrossSDKError) {
	response, err := ac.httpClient.Do(request)
	if err != nil {
		return nil, common.NewWeCrossSDKFromError(common.RPC_ERROR, err)
	}
	return response, nil
}

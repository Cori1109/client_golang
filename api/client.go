// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package api provides clients the HTTP API's.
package api

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

// CancelableTransport is like net.Transport but provides
// per-request cancelation functionality.
type CancelableTransport interface {
	http.RoundTripper
	CancelRequest(req *http.Request)
}

// DefaultTransport is used if no Transport is set in Config.
var DefaultTransport CancelableTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
}

// Config defines configuration parameters for a new client.
type Config struct {
	// The address of the Prometheus to connect to.
	Address string

	// Transport is used by the Client to drive HTTP requests. If not
	// provided, DefaultTransport will be used.
	Transport CancelableTransport
}

func (cfg *Config) transport() CancelableTransport {
	if cfg.Transport == nil {
		return DefaultTransport
	}
	return cfg.Transport
}

// Client is the interface for an API client.
type Client interface {
	URL(ep string, args map[string]string) *url.URL
	Do(context.Context, *http.Request) (*http.Response, []byte, error)
}

// New returns a new Client.
//
// It is safe to use the returned Client from multiple goroutines.
func New(cfg Config) (Client, error) {
	u, err := url.Parse(cfg.Address)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimRight(u.Path, "/")

	return &httpClient{
		endpoint:  u,
		transport: cfg.transport(),
	}, nil
}

type httpClient struct {
	endpoint  *url.URL
	transport CancelableTransport
}

func (c *httpClient) URL(ep string, args map[string]string) *url.URL {
	p := path.Join(c.endpoint.Path, ep)

	for arg, val := range args {
		arg = ":" + arg
		p = strings.Replace(p, arg, val, -1)
	}

	u := *c.endpoint
	u.Path = p

	return &u
}

func (c *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	resp, err := ctxhttp.Do(ctx, &http.Client{Transport: c.transport}, req)

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return nil, nil, err
	}

	var body []byte
	done := make(chan struct{})
	go func() {
		body, err = ioutil.ReadAll(resp.Body)
		close(done)
	}()

	select {
	case <-ctx.Done():
		err = resp.Body.Close()
		<-done
		if err == nil {
			err = ctx.Err()
		}
	case <-done:
	}

	return resp, body, err
}

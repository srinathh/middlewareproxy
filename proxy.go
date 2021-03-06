// Package middlewareproxy is a debugging proxy server library based on httputil.ReverseProxy
// that allows you to interpose your own Handlers as middleware package.
package middlewareproxy

import (
	"github.com/justinas/alice"
	"github.com/srinathh/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Config defines the configuration required for running the proxy server. The struct
// members contain tags for github.com/artyom/autoflags if you'd like to embed this struct
// into your own program's command line options
type Config struct {
	Addr     string `flag:"addr,the address on which to serve the proxy"`
	Scheme   string `flag:"scheme,the scheme for the remote server"`
	Host     string `flag:"host,the host of the remote proxy server"`
	BasePath string `flag:"basepath,the base path of the remote proxy server"`
}

// Run starts the proxy server specified by the config and interposing the provided middleware
func Run(config Config, mw ...middleware.MiddleWare) error {
	chain := alice.New()
	for _, m := range mw {
		chain = chain.Append(alice.Constructor(m))
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: config.Scheme, Host: config.Host, Path: config.BasePath})
	return http.ListenAndServe(config.Addr, chain.Then(proxy))
}

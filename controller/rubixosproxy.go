package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
)

func (inst *Controller) ROSProxy(c *gin.Context) {
	remote, err := Builder("0.0.0.0", 1660)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}

	for key := range c.Writer.Header() {
		c.Writer.Header().Del(key)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxy_path")
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

package ssl

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"wechat/pkg/zlog"
	"strconv"
)

func TlsHandler(host string, port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开发模式（默认 debug）不强制 HTTPS 跳转，便于本地 http://host:port 调 API
		dev := gin.Mode() != gin.ReleaseMode
		secureMiddleware := secure.New(secure.Options{
			IsDevelopment: dev,
			SSLRedirect:   !dev,
			SSLHost:       host + ":" + strconv.Itoa(port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			zlog.Fatal(err.Error())
			return
		}

		c.Next()
	}
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	app "github.com/lianmc123/app-frame"
	"github.com/lianmc123/app-frame/transport/http"
	"log"
	stdhttp "net/http"
)

func main() {
	httpSrv := http.NewServer(":9090",
		func(router gin.IRouter) {
			/*router.Use(func(c *gin.Context) {
				if tr, ok := transport.FromServerContext(c.Request.Context()); ok {
					fmt.Println("operation before:::", tr.Operation())
				}
				c.Next()
				if tr, ok := transport.FromServerContext(c.Request.Context()); ok {
					fmt.Println("operation after:~~:", tr.Operation())
				}
			})*/
			router.GET("/hello/:name", func(c *gin.Context) {
				name := c.Param("name")
				if name == "error" {
					c.JSON(stdhttp.StatusBadRequest, gin.H{"name": "fuck"})
					c.Error(fmt.Errorf("name error fuck"))
					return
				} else if name == "error1" {
					c.JSON(stdhttp.StatusInternalServerError, gin.H{"name": "StatusInternalServerError"})
					//c.Error(fmt.Errorf("%s error", name))
					c.Error(fmt.Errorf("name error ~~~~~"))
					return
				}
				c.JSON(stdhttp.StatusOK, map[string]interface{}{
					"hello": name,
				})
			})
			router.GET("/fuck", func(c *gin.Context) {
				c.JSON(stdhttp.StatusOK, map[string]interface{}{
					"hello": "fuck",
				})
			})
		},
		false,
		/*http.Middleware(func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
				if tr, ok := transport.FromServerContext(ctx); ok {
					fmt.Println(tr.Endpoint())
					fmt.Println(tr.RequestHeader())
					fmt.Println(tr.Kind())
					fmt.Println("operation:!:", tr.Operation())
				}
				reply, err = handler(ctx, req)
				//fmt.Println(reply, err)
				return
			}
		}),*/
		http.CustomRecover(func(c *gin.Context, err interface{}) {
			c.JSON(stdhttp.StatusInternalServerError, map[string]interface{}{
				"error": "panic",
			})
		}),
	)

	application := app.New(
		app.Version("1.0.0"),
		app.Name("gin"),
		app.Service(httpSrv),
	)
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}

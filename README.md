## gii

### start

```go
package main

import (
	"log"
	"time"

	"github.com/lukedever/gii"
)

func onlyForV1() gii.HandlerFunc {
	return func(c *gii.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v1", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gii.New()
	r.Use(gii.Logger())
	r.GET("/hello", func(c *gii.Context) {
		c.JSON(200, gii.H{
			"code": "hello",
		})
	})

	v1 := r.Group("v1")
	v1.Use(onlyForV1())
	v1.POST("/world", func(c *gii.Context) {
		c.JSON(200, gii.H{
			"code": "world",
		})
	})

	v2 := r.Group("v2")
	v2.GET("/panic", func(c *gii.Context) {
		arr := []int{1, 2, 3}
		c.JSON(200, gii.H{
			"data": arr[4],
		})
	})

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
```

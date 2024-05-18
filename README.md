Gin Handler Wrapper

```go
package main

import (
    "github.com/twoojoo/ghw"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    handler := func (c *gin.Context) error {
        baz := c.Param("baz")

        if baz == "" {
            return ghw.ErrBadRequest("baz can't be empty")
        }

        c.JSON(200, gin.H{"baz": baz})
        return nil
    }

    r.GET("/foo/bar/:baz", ghw.Wrap(handler))

    err = r.Run(":3000")
	if err != nil {
		panic(err)
	}
}
```
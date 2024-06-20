package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

func main() {
	r := gin.Default()

	r.POST("/gitee-hook", func(c *gin.Context) {
		// 打印请求头部
		fmt.Println("收到请求，打印header:" + time.Now().String())
		headers := c.Request.Header
		for key, values := range headers {
			for _, value := range values {
				fmt.Printf("Header: %s = %s\n", key, value)
			}
		}
		fmt.Println("收到请求，打印body" + time.Now().String())

		// 读取并打印请求主体
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(500, "Failed to read request body")
			return
		}
		fmt.Printf("Body: %s\n", string(body))

		// 返回成功响应
		c.String(200, "Headers and body printed to console")
	})

	r.Run(":8888") // 监听并在 0.0.0.0:8080 上启动服务
}

package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func main() {
	r := gin.Default()

	// 创建一个中间件用于打印请求体
	r.Use(func(c *gin.Context) {
		// 读取请求体
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// 打印请求体
		fmt.Println(string(bodyBytes))

		// 把读取过的请求体重新赋值给 c.Request.Body 以便后续处理
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// 继续处理请求
		c.Next()
	})

	// 处理 POST 请求
	r.POST("/gitee-hook", func(c *gin.Context) {
		var jsonData map[string]interface{}
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Request received"})
	})

	// 启动服务并监听 8080 端口
	r.Run(":8888")
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	log.Println("Server Start")
	r := gin.Default()
	boot := r.Group("/boot")
	{
		boot.GET("/hello-go", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "Hello World!",
			})
		})
		boot.POST("/alice", actAlice)
		boot.POST("/test", func(c *gin.Context) {
			buf := make([]byte, 8192)
			n, _ := c.Request.Body.Read(buf)
			b := string(buf[0:n])
			fmt.Println(b)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
	log.Fatal(r.Run())
}

func actAlice(c *gin.Context) {
	var pushTime time.Duration = 800 //ミリ秒

	// gpio処理開始
	err := rpio.Open()
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error Failed to open GPIO")
		return
	}

	pin_boot := rpio.Pin(21) // GPIO21<-GPIO番号であることに注意
	pin_boot.Output()

	// pushTimeミリ秒出力（aliceの電源スイッチピンをショート）
	fmt.Println("Start GPIO operating")
	pin_boot.High()
	time.Sleep(pushTime * time.Millisecond)
	pin_boot.Low()

	//他のピンがdefault Inputなので戻しておく
	pin_boot.Input()

	//gpi処理終わり
	rpio.Close()

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

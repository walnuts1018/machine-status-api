package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	log.Println("start server...")
	r := gin.Default()
	r.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	r.PUT("/boot/alice", boot_alice)
	log.Fatal(r.Run())
}

func boot_alice(c *gin.Context) {

	var push_time time.Duration = 800 //ミリ秒

	// gpio処理開始
	err := rpio.Open()

	if err != nil {
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	pin_boot := rpio.Pin(40) // GPIO40*ピン*<-ピンであることに注意
	pin_boot.Output()

	// push_timeミリ秒出力（aliceの電源スイッチピンをショート）
	pin_boot.High()
	time.Sleep(push_time * time.Millisecond)
	pin_boot.Low()

	//gpi処理終わり
	rpio.Close()

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

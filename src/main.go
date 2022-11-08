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
	sw := r.Group("/switch")
	{
		sw.POST("/alice", switchAlice)
	}

	start := r.Group("/boot")
	{
		start.POST("/alice", switchAlice)
	}

	shutdown := r.Group("/shutdown")
	{
		shutdown.POST("/alice", switchAlice)
	}

	status := r.Group("/status")
	{
		status.GET("/hello-go", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "Hello World!",
			})
		})
	}

	log.Fatal(r.Run(":80"))
}

func switchAlice(c *gin.Context, pushTime time.Duration) {
	// gpio処理開始
	err := rpio.Open()
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "Server Error Failed to open GPIO")
		return
	}

	pin_boot := rpio.Pin(21) // GPIO21<-GPIO番号であることに注意
	pin_boot.Output()

	// push_timeミリ秒出力（aliceの電源スイッチピンをショート）
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

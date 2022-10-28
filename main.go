package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("start server...")
	r := gin.Default()
	r.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	r.PUT("/somePut", boot_machine)
	log.Fatal(r.Run())
}

func boot_machine(c *gin.Context, name string) {
	
    if err != nil{
        c.String(http.StatusInternalServerError, "Server Error")
        return
    }
    c.JSON(http.StatusCreated, gin.H{
        "status": "ok",
    })
}
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luthermonson/go-proxmox"
	"github.com/stianeikeland/go-rpio/v4"
)

var PVE_API_URL, PVE_API_TOKEN_ID, PVE_API_SECRET string

var pveClient *proxmox.Client

type MachineStatus int

func (s MachineStatus) ToString() string {
	switch s {
	case 1:
		return "Healthy"
	case 0:
		return "Unknown"
	case -1:
		return "Unhealthy"
	case -2:
		return "Inactive"
	default:
		return "Error"
	}
}

const (
	Healthy   MachineStatus = 1  // Power ON and can access
	Unknown   MachineStatus = 0  // Status Unknown
	Unhealthy MachineStatus = -1 //Power ON, but can not access
	Inactive  MachineStatus = -2 // Power OFF, but can not access
)

func main() {
	loadEnv()
	err := pveLogin()
	if err != nil {
		log.Println(err)
		log.Printf("Retry on API access...")
	}
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/status", func(context *gin.Context) {
			time.Sleep(1000 * time.Second)
			context.JSON(200, gin.H{
				"message": "api server is running",
			})
		})

		v1.POST("/start/:machineName", machineStart)
		v1.POST("/stop/:machineName", machineStop)
		v1.POST("/status/:machineName", machineStatus)
	}
	log.Printf("Server Started")
	log.Fatal(r.Run())
}

func loadEnv() {
	err := godotenv.Load("./.config/.env")

	if err != nil {
		log.Printf("failed to open env: %s", err)
	}

	PVE_API_URL = os.Getenv("PVE_API_URL")
	PVE_API_TOKEN_ID = os.Getenv("PVE_API_TOKEN_ID")
	PVE_API_SECRET = os.Getenv("PVE_API_SECRET")
}

func pveLogin() error {
	pveClient = proxmox.NewClient(PVE_API_URL,
		proxmox.WithAPIToken(PVE_API_TOKEN_ID, PVE_API_SECRET),
	)

	version, err := pveClient.Version()
	if err != nil {
		return fmt.Errorf("failed to access pve server: %w", err)
	}

	log.Printf("Successfully logged in to pve server (Version:%s)", version.Release)
	return nil
}

func machineStart(c *gin.Context) {
	targetMachineName := c.Param("machineName")
	if targetMachineName == "alice" {
		startAlice(c)
		return
	}

	status, err := getAliceStatus()
	if status != Healthy {
		message := "alice is not in healthy status: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
		log.Println(message)
		return
	}

	node, err := pveClient.Node("alice")
	if err != nil {
		err := pveLogin()
		message := "Failed to find alice: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
		log.Println(message)
		return
	}
	vms, err := node.VirtualMachines()
	if err != nil {
		err := pveLogin()
		message := "Failed to find vm list: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
		log.Println(message)
		return
	}

	var message string
	for _, vm := range vms {
		if vm.Name == targetMachineName {
			isrunning := vm.IsRunning()
			if isrunning {
				err := fmt.Errorf(vm.Name + " is already running")
				message = "Failed to stop " + targetMachineName + ": " + err.Error()
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": message,
				})
				log.Printf("Failed to start %s: %s", targetMachineName, err)
				return
			} else {
				_, err := vm.Start()
				if err != nil {
					message = "Failed to start " + targetMachineName + ": " + err.Error()
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": message,
					})
					log.Printf("Failed to start %s: %s", targetMachineName, err)
					return
				}
				message = "Starting vm: " + targetMachineName + "..."
				c.JSON(http.StatusOK, gin.H{
					"message": message,
				})
				log.Println(message)
				return
			}
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "Failed to start " + targetMachineName + ": " + "no such vm",
	})
}

func machineStop(c *gin.Context) {
	targetMachineName := c.Param("machineName")
	if targetMachineName == "alice" {
		stopAlice(c)
		return
	}

	status, err := getAliceStatus()
	if status != Healthy {
		message := "alice is not in healthy status: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
		log.Println(message)
		return
	}

	node, err := pveClient.Node("alice")
	if err != nil {
		message := "Failed to find alice: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
		log.Println(message)
		return
	}
	vms, err := node.VirtualMachines()
	if err != nil {
		message := "Failed to find vm list: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
		log.Println(message)
		return
	}

	var message string
	for _, vm := range vms {
		if vm.Name == targetMachineName {
			isrunning := vm.IsRunning()
			if isrunning {
				_, err := vm.Shutdown()
				if err != nil {
					_, err = vm.Stop()
					if err != nil {
						message = "Failed to stop " + targetMachineName + ": " + err.Error()
						c.JSON(http.StatusInternalServerError, gin.H{
							"message": message,
						})
						log.Printf("Failed to stop %s: %s", targetMachineName, err)
						return
					}
				}
				message = "Stopping vm: " + targetMachineName + "..."
				c.JSON(http.StatusOK, gin.H{
					"message": message,
				})
				log.Println(message)
				return
			} else {
				err := fmt.Errorf(vm.Name + " is already stopped")
				message = "Failed to stop " + targetMachineName + ": " + err.Error()
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": message,
				})
				log.Printf("Failed to stop %s: %s", targetMachineName, err)
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Failed to stop " + targetMachineName + ": " + "no such vm",
	})
}

func machineStatus(c *gin.Context) {
	targetMachineName := c.Param("machineName")
	if targetMachineName == "alice" {
		status, _ := getAliceStatus()
		c.JSON(http.StatusOK, gin.H{
			"status": status.ToString(),
		})
		log.Println(status.ToString())
		return
	}

	aliceStatus, err := getAliceStatus()
	if aliceStatus != Healthy {
		message := "alice is not in healthy status: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Unknown",
			"message": message,
		})
		log.Println(message)
		return
	}

	node, err := pveClient.Node("alice")
	if err != nil {
		message := "Failed to find alice: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
		log.Println(message)
		return
	}
	vms, err := node.VirtualMachines()
	if err != nil {
		message := "Failed to find vm list: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": message,
		})
		log.Println(message)
		return
	}

	for _, vm := range vms {
		if vm.Name == targetMachineName {
			log.Println(vm)
			isrunning := vm.IsRunning()
			if isrunning {
				var vmstatus = Healthy
				c.JSON(http.StatusOK, gin.H{
					"status": Healthy.ToString(),
				})
				log.Println(vm.Name + ": " + vmstatus.ToString())
			} else {
				var vmstatus = Inactive
				c.JSON(http.StatusOK, gin.H{
					"status": vmstatus.ToString(),
				})
				log.Println(vm.Name + ": " + vmstatus.ToString())
			}
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"message": "Failed to get status " + targetMachineName + ": " + "no such vm",
	})
}

func startAlice(c *gin.Context) {
	var pushTime time.Duration = 800 * time.Millisecond

	status, _ := getAliceStatus()
	if status != Inactive {
		err := fmt.Errorf("failed to start alice: alice is already running")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		log.Println(err)
		return
	}
	err := pressAliceSwitch(pushTime)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		log.Printf("Failed to start alice: %s", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Starting Alice...",
	})
}

func stopAlice(c *gin.Context) {
	var pushTime time.Duration = 800 * time.Millisecond

	status, _ := getAliceStatus()
	if status == Inactive {
		err := fmt.Errorf("failed to stop alice: alice is already stopping")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		log.Println(err)
		return
	}

	err := pressAliceSwitch(pushTime)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		log.Printf("Failed to stop alice: %s", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Stopping Alice...",
	})
}

func pressAliceSwitch(pushTime time.Duration) error {
	// gpio処理開始
	err := rpio.Open()
	if err != nil {
		return fmt.Errorf("failed to open gpio: %w", err)
	}

	pwSwPin := rpio.Pin(21) // GPIO21<-GPIO番号であることに注意
	pwSwPin.Output()

	// pushTimeミリ秒出力（aliceの電源スイッチピンをショート）
	log.Printf("Start GPIO operating: 21")
	pwSwPin.High()
	time.Sleep(pushTime)
	pwSwPin.Low()

	//他のピンがdefault Inputなので戻しておく
	pwSwPin.Input()

	//gpi処理終わり
	err = rpio.Close()
	if err != nil {
		return fmt.Errorf("failed to close gpio: %w", err)
	}
	return nil
}

func getAliceStatus() (MachineStatus, error) {
	pin, err := readAlicePWLed()

	if err != nil {
		return Unknown, err
	}

	if pin == rpio.Low {
		return Inactive, fmt.Errorf("alice is not running: %w", err)
	} else if pin == rpio.High {
		_, err := pveClient.Version()
		if err != nil {
			err := pveLogin()
			if err != nil {
				return Unhealthy, fmt.Errorf("can not access alice: %w", err)
			}
		}
		return Healthy, nil
	}
	return Unknown, err
}

func readAlicePWLed() (rpio.State, error) {
	// gpio処理開始
	err := rpio.Open()

	if err != nil {
		return 0, fmt.Errorf("failed to open gpio: %w", err)
	}

	pwLedPin := rpio.Pin(16) // GPIO16<-GPIO番号であることに注意
	pwLedPin.Input()

	log.Printf("Start GPIO operating: 16")
	status := pwLedPin.Read()

	//gpi処理終わり
	err = rpio.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to close gpio: %w", err)
	}
	return status, nil
}

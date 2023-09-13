package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/walnuts1018/machine-status-api/usecase"
)

var machineUsecase *usecase.MachineUsecase

func NewHandler(mu *usecase.MachineUsecase) *gin.Engine {
	machineUsecase = mu
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.GET("/status", GetMachineStatus)
		v1.POST("/machines/start/:machineName", StartMachine)
		v1.POST("/machines/stop/:machineName", StopMachine)
	}
	return r
}

func StartMachine(ctx *gin.Context) {
	machineName := ctx.Param("machineName")
	task := usecase.NewTasks(func() error {
		err := machineUsecase.StartMachine(machineName)
		if err != nil {
			return fmt.Errorf("failed to start machine: %w", err)
		}
		return nil
	})
	t, err := json.Marshal(task.RegisteredAt)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": fmt.Sprintf("failed to marshal time: %v", err),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"id":            task.ID,
		"status":        task.Status,
		"registered_at": t,
	})
}

func StopMachine(ctx *gin.Context) {
	machineName := ctx.Param("machineName")
	task := usecase.NewTasks(func() error {
		err := machineUsecase.StopMachine(machineName)
		if err != nil {
			return fmt.Errorf("failed to stop machine: %w", err)
		}
		return nil
	})
	t, err := json.Marshal(task.RegisteredAt)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": fmt.Sprintf("failed to marshal time: %v", err),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"id":            task.ID,
		"status":        task.Status,
		"registered_at": t,
	})
}

func GetMachineStatus(ctx *gin.Context) {
	machineName := ctx.Param("machineName")
	status, err := machineUsecase.GetMachineStatus(machineName)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": fmt.Sprintf("failed to get status: %v", err),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"name":   machineName,
		"status": status,
	})
}

func GetTaskStatus(ctx *gin.Context) {
	taskID := ctx.Param("taskID")
	task := usecase.FindTaskByID(taskID)
	if task == nil {
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("task not found: %s", taskID),
		})
		return
	}

	tr, err := json.Marshal(task.RegisteredAt)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": fmt.Sprintf("failed to marshal time: %v", err),
		})
		return
	}

	ts, err := json.Marshal(task.StartedAt)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": fmt.Sprintf("failed to marshal time: %v", err),
		})
		return
	}

	tf, err := json.Marshal(task.FinishedAt)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": fmt.Sprintf("failed to marshal time: %v", err),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"id":            task.ID,
		"status":        task.Status,
		"registered_at": tr,
		"started_at":    ts,
		"finished_at":   tf,
	})
}

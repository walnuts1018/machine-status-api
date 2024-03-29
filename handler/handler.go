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
		machines := v1.Group("/machines")
		{
			machines.GET("/status/:machineName", GetMachineStatus)
			start := machines.Group("/start")
			{
				start.POST("/:machineName", StartMachine)
				start.POST("/:machineName/automated", StartMachineAutomated)
			}
			stop := machines.Group("/stop")
			{
				stop.POST("/:machineName", StopMachine)
				stop.POST("/:machineName/automated", StopMachineAutomated)

			}
		}
		tasks := v1.Group("/tasks")
		{
			tasks.GET("/:taskID", GetTaskStatus)
		}
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
		"registered_at": string(t),
	})
}

func StartMachineAutomated(ctx *gin.Context) {
	machineName := ctx.Param("machineName")
	task := usecase.NewTasks(func() error {
		err := machineUsecase.StartMachineAutomated(machineName)
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
		"registered_at": string(t),
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
		"registered_at": string(t),
	})
}

func StopMachineAutomated(ctx *gin.Context) {
	machineName := ctx.Param("machineName")
	task := usecase.NewTasks(func() error {
		err := machineUsecase.StopMachineAutomated(machineName)
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
		"registered_at": string(t),
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
		"registered_at": string(tr),
		"started_at":    string(ts),
		"finished_at":   string(tf),
	})
}

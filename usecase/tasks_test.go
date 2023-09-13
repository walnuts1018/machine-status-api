package usecase

import (
	"context"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		ctx context.Context
	}
	results := make([]int, 0)
	tasks := make([]*tasks, 0)
	for i := 0; i < 10; i++ {
		tmp := i
		task := NewTasks(func() error {
			results = append(results, tmp)
			return nil
		})
		tasks = append(tasks, task)
	}
	tasks = append(tasks, NewTasks(func() error {
		cancel()
		return fmt.Errorf("error")
	}))

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				ctx: ctx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.ctx)
			if len(results) != len(tasks)-1 {
				t.Errorf("result len error, want %d, got %d", len(tasks), len(results))
			}
			for i, result := range results {
				if i != result {
					t.Errorf("task result error, want %d, got %d", i, result)
				}
			}
			for _, task := range tasks[:len(tasks)-2] {
				discoveredTask := FindTaskByID(task.ID)
				if discoveredTask.Status != Success {
					t.Errorf("task status error, want %s, got %s", Success, discoveredTask.Status)
				}
			}
			if tasks[len(tasks)-1].Status != Failure {
				t.Errorf("task status error, want %s, got %s", Failure, tasks[len(tasks)-1].Status)
			}

		})
	}
}

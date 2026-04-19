package domain

import "time"

type Statistics struct {
	TasksCreated               int
	TasksCompleted             int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(
	tasksCreated,
	tasksCompleted int,
	taskCompletionRate *float64,
	tasksAverageCompletionTime *time.Duration,
) *Statistics {
	return &Statistics{
		tasksCreated,
		tasksCompleted,
		taskCompletionRate,
		tasksAverageCompletionTime,
	}
}

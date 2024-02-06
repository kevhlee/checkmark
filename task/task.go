package task

import (
	"slices"
	"strings"
)

type Task struct {
	Done     bool     `json:"done"`
	Name     string   `json:"name"`
	Priority Priority `json:"priority"`
}

func CompareTasks(a, b Task) int {
	diff := b.Priority - a.Priority
	if diff == 0 {
		return strings.Compare(a.Name, b.Name)
	}
	return int(diff)
}

func SortTasks(tasks []Task) {
	slices.SortStableFunc(tasks, CompareTasks)
}

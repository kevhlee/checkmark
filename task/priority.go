package task

import "fmt"

type Priority int

func (p Priority) Name() string {
	return priorityNameMap[p]
}

func (p Priority) Symbol() string {
	return prioritySymbolMap[p]
}

func (p Priority) String() string {
	return fmt.Sprintf("%s %s", prioritySymbolMap[p], priorityNameMap[p])
}

const (
	LowPriority Priority = iota
	HighPriority
	FirePriority
)

var (
	Priorities = []Priority{
		LowPriority,
		HighPriority,
		FirePriority,
	}

	priorityNameMap = map[Priority]string{
		LowPriority:  "Low",
		HighPriority: "High",
		FirePriority: "Fire",
	}

	prioritySymbolMap = map[Priority]string{
		LowPriority:  "🎯",
		HighPriority: "❗",
		FirePriority: "🔥",
	}
)

package task

type Priority int

func (p Priority) String() string {
	return priorityStringMap[p]
}

func (p Priority) Symbol() string {
	return prioritySymbolMap[p]
}

const (
	LowPriority Priority = iota
	HighPriority
	FirePriority
)

var (
	priorityStringMap = map[Priority]string{
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

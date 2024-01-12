package schedule

type TypeStrategy string

const (
	TypeStrategyStateful TypeStrategy = "stateful"
	TypeStrategyTick     TypeStrategy = "tick"
)

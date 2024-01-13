package invocation

type TypeStrategy string

const (
	TypeStrategyHistoricWindow TypeStrategy = "historic_window"
	TypeStrategyTick           TypeStrategy = "tick"
)

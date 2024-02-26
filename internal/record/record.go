package record

type Record interface {
	Map() map[string]any
}

type Result interface {
	Records() []Record
}

package record

type Record interface {
	Bytes() ([]byte, error)
}

type Result interface {
	Records() []Record
}

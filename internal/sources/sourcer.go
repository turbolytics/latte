package sources

type Type string

const (
	TypePostgres Type = "postgres"
)

type Sourcer interface {
}

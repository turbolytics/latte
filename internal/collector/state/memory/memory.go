package memory

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/turbolytics/collector/internal/collector/state"
	"github.com/turbolytics/collector/internal/timeseries"
	"time"
)

type Store struct {
	db *sql.DB
}

func (m *Store) Close() error {
	m.Close()
	return nil
}

func (m *Store) MostRecentInvocation(ctx context.Context, collector string) (*state.Invocation, error) {
	row := m.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT 
	collector_name,
	time, 
	window_start,
	window_end
FROM
    invocations
WHERE
    collector_name = '%s'
ORDER BY 
	time DESC
LIMIT 1
`, collector))

	w := &timeseries.Window{}
	i := &state.Invocation{}
	err := row.Scan(&i.CollectorName, &i.Time, &w.Start, &w.End)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	i.Window = w
	return i, nil
}

func (m *Store) SaveInvocation(invocation *state.Invocation) error {
	_, err := m.db.Exec(
		fmt.Sprintf(
			"INSERT INTO invocations VALUES('%s', '%s', '%s', '%s')",
			invocation.CollectorName,
			invocation.Time.Format(time.RFC3339),
			invocation.Window.Start.Format(time.RFC3339),
			invocation.Window.End.Format(time.RFC3339),
		),
	)
	return err
}

func NewFromGenericConfig(m map[string]any) (*Store, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
CREATE TABLE invocations (
    collector_name VARCHAR,
    time TIMESTAMP,
    window_start TIMESTAMP,
    window_end TIMESTAMP
)
`)
	if err != nil {
		return nil, err
	}

	s := &Store{
		db: db,
	}

	return s, nil
}

package partition

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPartitioner_Render(t *testing.T) {
	p, err := New("year={{.Year}}/month={{.Month}}/day={{.Day}}/hour={{.Hour}}/minute={{.Minute}}/test")
	assert.NoError(t, err)

	d := time.Date(2024, 03, 01, 00, 00, 00, 0, time.UTC)
	s, err := p.Render(d)
	assert.NoError(t, err)
	assert.Equal(t,
		"year=2024/month=3/day=1/hour=0/minute=0/test",
		s,
	)
}

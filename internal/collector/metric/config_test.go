package metric

import (
	"github.com/stretchr/testify/assert"
	"github.com/turbolytics/collector/internal/collector/source"
	"testing"
)

func TestInvocation_SetDefaults_Schedule(t *testing.T) {
	conf := &Config{}
	defaults(conf)
	assert.Equal(t, source.TypeStrategyTick, conf.Source.Strategy)
}

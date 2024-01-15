package sources

import "github.com/turbolytics/collector/internal/timeseries"

type Context struct {
	Window timeseries.Bucket
}

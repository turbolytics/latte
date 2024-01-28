package timeseries

import "time"

type WindowConfig struct {
	Window string `mapstructure:"window"`

	// mapstructure does not parse to native go types
	// maybe there is a way to implement an unmarshaller
	duration time.Duration
}

func (wc *WindowConfig) Init() error {
	d, err := time.ParseDuration(wc.Window)
	if err != nil {
		return err
	}
	wc.duration = d
	return nil
}

func (wc WindowConfig) Duration() *time.Duration {
	return &wc.duration
}

package invoker

import (
	"context"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/sink"
	"github.com/turbolytics/latte/internal/source"
	"time"
)

type TestTransformer struct{}

func (tt TestTransformer) Transform(r record.Result) error {
	tr := r.(TestResult)

	for _, rec := range tr.records {
		tm := make(map[string]any)
		for k, v := range rec.Map() {
			tm[k+"_transformed"] = v.(string) + "_transformed"
		}
		rec.m = tm
	}

	return nil
}

type TestRecord struct {
	m map[string]any
}

func (tr TestRecord) Map() map[string]any {
	return tr.m
}

type TestResult struct {
	records []*TestRecord
}

func (tr TestResult) Records() []record.Record {
	var rs []record.Record
	for _, r := range tr.records {
		rs = append(rs, r)
	}
	return rs
}

type TestSink struct {
	closes int
	writes []record.Record
}

func (ts TestSink) Type() sink.Type {
	return "tester"
}

func (ts *TestSink) Write(r record.Record) (int, error) {
	ts.writes = append(ts.writes, r)
	return 0, nil
}

func (ts *TestSink) Flush() error {
	return nil
}

func (ts *TestSink) Close() error {
	ts.closes++
	return nil
}

type TestSourcer struct {
	w  *time.Duration
	t  source.Type
	tr TestResult
}

func (ts TestSourcer) Window() *time.Duration {
	return ts.w
}

func (ts TestSourcer) Type() source.Type {
	return ts.t
}

func (ts TestSourcer) Source(ctx context.Context) (record.Result, error) {
	return ts.tr, nil
}

type TestConfig struct {
	invocationStrategy TypeStrategy
	name               string
	sourcer            TestSourcer
	sinks              []*TestSink
	transformer        TestTransformer
}

func (t TestConfig) Transformer() Transformer {
	return t.transformer
}

func (t TestConfig) InvocationStrategy() TypeStrategy {
	return t.invocationStrategy
}

func (t TestConfig) Name() string {
	return t.name
}

func (t TestConfig) Sinks() []Sinker {
	var sinks []Sinker
	for _, sink := range t.sinks {
		sinks = append(sinks, sink)
	}
	return sinks
}

func (t TestConfig) Schedule() Schedule {
	return nil
}

func (t TestConfig) Sourcer() Sourcer {
	return t.sourcer
}

func (t TestConfig) Storer() Storer {
	return nil
}

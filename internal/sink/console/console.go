package console

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/encoding"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/sink"
	"io"
	"os"
)

type config struct {
	Encoding encoding.Config
}

type Console struct {
	encoder encoding.Encoder
	w       io.Writer
}

func (c *Console) Close() error {
	return nil
}

func (c *Console) Flush() error {
	return nil
}

func (c *Console) Type() sink.Type {
	return sink.TypeConsole
}

func (c *Console) Write(r record.Record) (int, error) {
	buf := &bytes.Buffer{}
	if err := c.encoder.Init(buf); err != nil {
		return 0, nil
	}

	if err := c.encoder.Write(r.Map()); err != nil {
		return 0, err
	}

	bs := buf.Bytes()
	fmt.Println(string(bs))
	return 0, nil
}

func NewFromGenericConfig(m map[string]any) (*Console, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(os.Stdout)

	e, err := encoding.NewEncoder(conf.Encoding)

	if err != nil {
		return nil, err
	}

	return &Console{
		encoder: e,
		w:       writer,
	}, nil
}

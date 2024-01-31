package console

import (
	"bufio"
	"fmt"
	"github.com/turbolytics/latte/internal/sinks"
	"io"
	"os"
)

type Console struct {
	w io.Writer
}

func (c *Console) Close() error {
	return nil
}

func (c *Console) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))
	return 0, nil
}

func (c *Console) Type() sinks.Type {
	return sinks.TypeConsole
}

func NewFromGenericConfig(m map[string]any) (*Console, error) {
	writer := bufio.NewWriter(os.Stdout)
	return &Console{
		w: writer,
	}, nil
}

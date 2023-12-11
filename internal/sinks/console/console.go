package console

import (
	"bufio"
	"io"
	"os"
)

type Console struct {
	w io.Writer
}

func (c *Console) Write(bs []byte) (int, error) {
	return c.w.Write(bs)
}

func NewFromGenericConfig(m map[string]any) (*Console, error) {
	writer := bufio.NewWriter(os.Stdout)
	return &Console{
		w: writer,
	}, nil
}

package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Console struct {
	w io.Writer
}

func (c *Console) Write(bs []byte) (int, error) {
	fmt.Println("here1")
	fmt.Println(string(bs))
	fmt.Println("here2")
	return 0, nil
}

func NewFromGenericConfig(m map[string]any) (*Console, error) {
	writer := bufio.NewWriter(os.Stdout)
	return &Console{
		w: writer,
	}, nil
}

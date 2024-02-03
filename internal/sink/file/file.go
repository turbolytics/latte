package file

import (
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/sink/type"
	"os"
)

type config struct {
	Path string
}

type File struct {
	config config
	f      *os.File
}

func (fs *File) Close() error {
	return fs.f.Close()
}

func (fs *File) Type() _type.Type {
	return _type.TypeHTTP
}

func (fs *File) Write(bs []byte) (int, error) {
	n, err := fs.f.Write(bs)
	if err != nil {
		return n, err
	}
	_, err = fs.f.Write([]byte("\n"))
	return n, err
}

func NewFromGenericConfig(m map[string]any, validate bool) (*File, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	var f *os.File
	var err error
	if !validate {
		f, err = os.OpenFile(conf.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
	}

	return &File{
		config: conf,
		f:      f,
	}, nil
}

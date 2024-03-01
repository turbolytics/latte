package file

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"github.com/turbolytics/latte/internal/encoding"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/sink"
	"os"
)

type config struct {
	Encoding encoding.Config
	Path     string
}

type File struct {
	config  config
	encoder encoding.Encoder
	f       *os.File
}

func (fs *File) Close() error {
	return fs.f.Close()
}

func (fs *File) Flush() error {
	return nil
}

func (fs *File) Type() sink.Type {
	return sink.TypeHTTP
}

func (fs *File) Write(r record.Record) (int, error) {
	buf := &bytes.Buffer{}
	if err := fs.encoder.Init(buf); err != nil {
		return 0, nil
	}

	if err := fs.encoder.Write(r.Map()); err != nil {
		return 0, err
	}

	bs := buf.Bytes()
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

	e, err := encoding.NewEncoder(conf.Encoding)

	if err != nil {
		return nil, err
	}

	return &File{
		config:  conf,
		encoder: e,
		f:       f,
	}, nil
}

package template

import (
	"bytes"
	"os"
	"text/template"
)

func Parse(bs []byte) ([]byte, error) {
	funcMap := template.FuncMap{
		"getEnv": func(key string) string {
			return os.Getenv(key)
		},
		"getEnvOrDefault": func(key string, d string) string {
			envVal := os.Getenv(key)
			if envVal == "" {
				return d
			}

			return envVal
		},
	}
	t, err := template.New("config").Funcs(funcMap).Parse(string(bs))
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := t.Execute(&out, nil); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

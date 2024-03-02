package partition

import (
	"bytes"
	"text/template"
	"time"
)

type tmplData struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
}

type Partitioner struct {
	template *template.Template
}

func (p *Partitioner) Render(t time.Time) (string, error) {
	td := tmplData{
		Year:   t.Year(),
		Month:  int(t.Month()),
		Day:    t.Day(),
		Hour:   t.Hour(),
		Minute: t.Minute(),
	}
	var out bytes.Buffer
	err := p.template.Execute(&out, td)
	return out.String(), err
}

func New(partition string) (*Partitioner, error) {
	t, err := template.New("config").Parse(string(partition))
	if err != nil {
		return nil, err
	}
	p := &Partitioner{
		template: t,
	}

	return p, nil
}

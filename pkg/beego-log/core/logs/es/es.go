package es

import (
	logs2 "agile-proxy/pkg/beego-log/core/logs"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
)

// NewES returns a LoggerInterface
func NewES() logs2.Logger {
	cw := &esLogger{
		Level:       logs2.LevelDebug,
		indexNaming: indexNaming,
	}
	return cw
}

// esLogger will log msg into ES
// before you using this implementation,
// please import this package
// usually means that you can import this package in your main package
// for example, anonymous:
// import _ "agile-proxy/pkg/core/logs/es"
type esLogger struct {
	*elasticsearch.Client
	DSN       string `json:"dsn"`
	Level     int    `json:"level"`
	formatter logs2.LogFormatter
	Formatter string `json:"formatter"`

	indexNaming IndexNaming
}

func (el *esLogger) Format(lm *logs2.LogMsg) string {
	msg := lm.OldStyleFormat()
	idx := LogDocument{
		Timestamp: lm.When.Format(time.RFC3339),
		Msg:       msg,
	}
	body, err := json.Marshal(idx)
	if err != nil {
		return msg
	}
	return string(body)
}

func (el *esLogger) SetFormatter(f logs2.LogFormatter) {
	el.formatter = f
}

// {"dsn":"http://localhost:9200/","level":1}
func (el *esLogger) Init(config string) error {
	err := json.Unmarshal([]byte(config), el)
	if err != nil {
		return err
	}
	if el.DSN == "" {
		return errors.New("empty dsn")
	} else if u, err := url.Parse(el.DSN); err != nil {
		return err
	} else if u.Path == "" {
		return errors.New("missing prefix")
	} else {
		conn, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{el.DSN},
		})
		if err != nil {
			return err
		}
		el.Client = conn
	}
	if len(el.Formatter) > 0 {
		fmtr, ok := logs2.GetFormatter(el.Formatter)
		if !ok {
			return errors.New(fmt.Sprintf("the formatter with name: %s not found", el.Formatter))
		}
		el.formatter = fmtr
	}
	return nil
}

// WriteMsg writes the msg and level into es
func (el *esLogger) WriteMsg(lm *logs2.LogMsg) error {
	if lm.Level > el.Level {
		return nil
	}

	msg := el.formatter.Format(lm)

	req := esapi.IndexRequest{
		Index:        indexNaming.IndexName(lm),
		DocumentType: "logs",
		Body:         strings.NewReader(msg),
	}
	_, err := req.Do(context.Background(), el.Client)
	return err
}

// Destroy is a empty method
func (el *esLogger) Destroy() {
}

// Flush is a empty method
func (el *esLogger) Flush() {
}

type LogDocument struct {
	Timestamp string `json:"timestamp"`
	Msg       string `json:"msg"`
}

func init() {
	logs2.Register(logs2.AdapterEs, NewES)
}

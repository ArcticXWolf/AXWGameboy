package internal

import (
	"container/ring"
	"fmt"
	"log"
)

type RingLog struct {
	name    string
	content *ring.Ring
}

func NewRingLog(name string, size int) *RingLog {
	return &RingLog{
		name:    name,
		content: ring.New(size),
	}
}

func (r *RingLog) Printf(format string, a ...interface{}) {
	r.content.Value = fmt.Sprintf(format, a...)
	r.content = r.content.Next()
}

func (r *RingLog) GetMessages() (result []string) {
	r.content.Do(func(p interface{}) {
		if p != nil {
			result = append(result, p.(string))
		}
	})
	return
}

type RingLogger struct {
	Loggers     []*RingLog
	defaultSize int
}

func NewRingLogger(defaultSize int) *RingLogger {
	return &RingLogger{
		Loggers:     make([]*RingLog, 0),
		defaultSize: defaultSize,
	}
}

func (r *RingLogger) Printf(category string, format string, a ...interface{}) {
	found := false
	for _, v := range r.Loggers {
		if v.name == category {
			v.Printf(format, a...)
			found = true
		}
	}

	if !found {
		newLogger := NewRingLog(category, r.defaultSize)
		r.Loggers = append(r.Loggers, newLogger)
		newLogger.Printf(format, a...)
		log.Printf("Allocated new logger %s with length %d", newLogger.name, newLogger.content.Len())
	}
}

func (r *RingLogger) GetAllMessages() (result map[string][]string) {
	result = make(map[string][]string)
	for _, v := range r.Loggers {
		result[v.name] = v.GetMessages()
	}
	return
}

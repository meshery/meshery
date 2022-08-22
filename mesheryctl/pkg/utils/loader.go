package utils

import (
	"fmt"
	"sync"
	"time"
)

var symbols = []string{"|", "/", "-", "\\"}

type Loader struct {
	Message string
	runChan chan struct{}
	stop    sync.Once
}

func NewLoader(message string) *Loader {
	l := &Loader{
		Message: message,
		runChan: make(chan struct{}),
	}
	return l.Start()
}

func (l *Loader) Start() *Loader {
	go l.writer()
	return l
}

func (l *Loader) Stop() {
	l.stop.Do(func() {
		close(l.runChan)
		l.clear()
	})
	fmt.Printf("\033[2K")
}

func (l *Loader) animate() {
	var m string
	for i := 0; i < len(symbols); i++ {
		m = l.Message + " " + symbols[i]
		Log.Info(m)
		time.Sleep(time.Millisecond * 200)
		l.clear()
	}
}

func (l *Loader) writer() {
	l.animate()
	for {
		select {
		case <-l.runChan:
			return
		default:
			l.animate()
		}
	}
}

func (l *Loader) clear() {
	// \033[ is the ANSI escape
	// to erase the line before and after the cursor
	fmt.Printf("\033[2K")
	// to shift the cursor up by one line
	fmt.Printf("\033[1A")
}

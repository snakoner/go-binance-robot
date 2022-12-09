package logger

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Logger struct {
	Filename string
	File     *os.File
}

func New(filename string) *Logger {
	return &Logger{
		Filename: filename,
	}
}

func (logger *Logger) Open() error {
	f, err := os.OpenFile(logger.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
		return err
	}

	logger.File = f
	return nil
}

func formatTimeNow() string {
	t := time.Now()
	h := strconv.Itoa(t.Hour())
	m := strconv.Itoa(t.Minute())
	s := strconv.Itoa(t.Second())
	if len(h) == 1 {
		h = "0" + h
	}

	if len(m) == 1 {
		m = "0" + m
	}

	if len(s) == 1 {
		s = "0" + s
	}

	result := fmt.Sprintf("%s:%s:%s", h, m, s)
	return result
}

func (logger *Logger) Write(s string) error {
	_, err := logger.File.WriteString(fmt.Sprintf("%s : %s\n", formatTimeNow(), s))
	return err
}

func (logger *Logger) Close() error {
	err := logger.File.Close()
	return err
}

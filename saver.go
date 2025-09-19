package main

import (
	"io"
	"log"
	"os"
)

type Saver struct {
	file *os.File
}

func NewSaver() Saver {
	saver := new(Saver)
	return *saver
}

func (saver *Saver) Start() error {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	saver.file = logFile

	multiWritter := io.MultiWriter(os.Stdout, saver.file)
	log.SetOutput(multiWritter)

	return nil
}

func (saver *Saver) Stop() error {
	err := saver.file.Close()
	log.SetOutput(os.Stdout)

	return err
}

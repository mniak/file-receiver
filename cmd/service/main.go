package main

import (
	"log"

	"receivefiles"

	"github.com/kardianos/service"
)

var logger service.Logger

type Service struct{}

func (p *Service) Start(s service.Service) error {
	return app.Start()
}

func (p *Service) Stop(s service.Service) error {
	return app.Stop()
}

var app receivefiles.Service

func main() {
	app = receivefiles.NewApp(receivefiles.AppParams{
		HttpParams: receivefiles.HttpParams{
			Port:             10777,
			ReceivedFilesDir: `C:\Users\Win\Desktop\Arquivos Recebidos`,
		},
	})

	svcConfig := &service.Config{
		Name:        "receivefiles",
		DisplayName: "File Receiver",
		Description: "To receive files from a web server",
	}

	prg := &Service{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

package main

import (
	"context"

	"github.com/kardianos/service"
)

type Program struct {
	HTTP       *HTTPServer
	context    context.Context
	cancelFunc func()
}

func NewProgram() *Program {
	ctx, cancel := context.WithCancel(context.Background())

	prog := Program{
		HTTP:       &HTTPServer{},
		context:    ctx,
		cancelFunc: cancel,
	}
	return &prog
}

func (p *Program) Start(s service.Service) error {
	p.HTTP.Start()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	return p.HTTP.Stop(p.context)
}

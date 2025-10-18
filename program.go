package main

import (
	"github.com/kardianos/service"
	"go.uber.org/multierr"
)

type Program struct {
	HTTP *HTTPServer
}

func NewProgram() *Program {
	prog := Program{
		HTTP: &HTTPServer{},
	}
	return &prog
}

func (p *Program) Start(s service.Service) error {
	var allErrors error
	err := p.HTTP.Start()
	multierr.AppendInto(&allErrors, err)
	return allErrors
}

func (p *Program) Stop(s service.Service) error {
	var allErrors error
	err := p.HTTP.Stop()
	multierr.AppendInto(&allErrors, err)
	return allErrors
}

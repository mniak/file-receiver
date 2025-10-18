package main

import (
	"runtime"

	"github.com/kardianos/service"
)

type NotAService struct{}

func (d *NotAService) Install() error {
	return nil
}

func (d *NotAService) Logger(errs chan<- error) (service.Logger, error) {
	return nil, nil
}

func (d *NotAService) Platform() string {
	return runtime.GOOS
}

func (d *NotAService) Restart() error {
	return nil
}

func (d *NotAService) Run() error {
	return nil
}

func (d *NotAService) Start() error {
	return nil
}

func (d *NotAService) Status() (service.Status, error) {
	return service.StatusUnknown, nil
}

func (d *NotAService) Stop() error {
	return nil
}

func (d *NotAService) String() string {
	return ""
}

func (d *NotAService) SystemLogger(errs chan<- error) (service.Logger, error) {
	return nil, nil
}

func (d *NotAService) Uninstall() error {
	return nil
}

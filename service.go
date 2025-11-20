package receivefiles

import (
	"sync"

	"go.uber.org/multierr"
)

type Service interface {
	Start() error
	Wait() error
	Stop() error
}

type MultiService []Service

func (multi *MultiService) Start() error {
	var wg sync.WaitGroup
	var errs error
	var errsMut sync.Mutex
	for _, s := range *multi {
		wg.Go(func() {
			err := s.Start()
			if err != nil {
				errsMut.Lock()
				multierr.AppendInto(&errs, err)
				errsMut.Unlock()
			}
		})
	}
	wg.Wait()
	return errs
}

func (multi *MultiService) Wait() error {
	var wg sync.WaitGroup
	var errs error
	var errsMut sync.Mutex
	for _, s := range *multi {
		wg.Go(func() {
			err := s.Wait()
			if err != nil {
				errsMut.Lock()
				multierr.AppendInto(&errs, err)
				errsMut.Unlock()
			}
		})
	}
	wg.Wait()
	return errs
}

func (multi *MultiService) Stop() error {
	var wg sync.WaitGroup
	var errs error
	var errsMut sync.Mutex
	for _, s := range *multi {
		wg.Go(func() {
			err := s.Stop()
			if err != nil {
				errsMut.Lock()
				multierr.AppendInto(&errs, err)
				errsMut.Unlock()
			}
		})
	}
	wg.Wait()
	return errs
}

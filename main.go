package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

func main() {
	program := NewProgram()

	cmd := cobra.Command{
		Use: "file-receiver",
		Run: func(cmd *cobra.Command, args []string) {
			if service.Interactive() {
				cobra.CheckErr(RunInteractive(program))
			} else {
				cobra.CheckErr(RunAsService(program))
			}
		},
	}

	cmd.Flags().IntVarP(&program.HTTP.Port, "port", "p", 10777, "The server of the http server")
	cmd.Execute()
}

func RunInteractive(program *Program) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT /*ctrl+c*/, syscall.SIGTERM /*kill*/)
	defer stop()

	svc := NotAService{}
	err := program.Start(&svc)
	if err != nil {
		return err
	}

	go func() {
		log.Println("Application is running. Press Ctrl+C or use VS Code stop button to shut down.")
		s := <-ctx.Done()
		log.Println("Shutdown signal received.", s)

		program.Stop(&svc)
	}()

	<-ctx.Done()
	log.Println("Exiting application.")
	return nil
}

func RunAsService(program *Program) error {
	svcConfig := service.Config{}
	svc, err := service.New(program, &svcConfig)
	if err != nil {
		return err
	}
	return svc.Run()
}

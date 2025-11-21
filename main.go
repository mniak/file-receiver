package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

func main() {
	var p AppParams
	var app Service
	cmd := cobra.Command{
		Use: "receivefiles",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			p.ReceivedFilesDir, err = filepath.Abs(p.ReceivedFilesDir)
			if err != nil {
				return err
			}

			app = NewApp(p)

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			app := NewApp(p)

			url := "https://flores7.mniak.dev"
			fmt.Printf("Access %s to send files\nor scan the QR Code below\n", url)
			showQRCode(url)
			fmt.Printf("The files will be stored on %s\n", p.ReceivedFilesDir)

			cobra.CheckErr(app.Start())
			cobra.CheckErr(app.Wait())
		},
	}
	cmd.PersistentFlags().IntVar(&p.Port, "port", 10777, "HTTP port of the server")
	cmd.PersistentFlags().StringVar(&p.ReceivedFilesDir, "save-to", "./uploads", "Where to save the received files")

	svcCmd := cobra.Command{
		Use: "service",
		Run: func(cmd *cobra.Command, args []string) {
			var logger service.Logger
			svcConfig := &service.Config{
				Name:        "receivefiles",
				DisplayName: "File Receiver",
				Description: "Receive files from a web page",
			}

			svc, err := service.New(NewSystemService(app), svcConfig)
			if err != nil {
				log.Fatal(err)
			}
			logger, err = svc.Logger(nil)
			if err != nil {
				log.Fatal(err)
			}
			err = svc.Run()
			if err != nil {
				logger.Error(err)
			}
		},
	}

	cmd.AddCommand(&svcCmd)
	cobra.CheckErr(cmd.Execute())
}

func showQRCode(text string) {
	config := qrterminal.Config{
		HalfBlocks: true,
		Level:      qrterminal.M,
		Writer:     os.Stdout,
	}
	qrterminal.GenerateWithConfig(text, config)
}

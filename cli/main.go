package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	cli "github.com/urfave/cli/v2"
	"ksadebiekorf.be/mailing/env"
)

func printLightRed(msg string) {
	if strings.Contains(os.Getenv("TERM"), "256color") {
		fmt.Printf("\u001b[38;5;202m%s\u001b[0m\n", msg)
	} else {
		fmt.Printf("\u001b[91m%s\u001b[0m\n", msg)
	}
}

// Mailing utility `mails`
func main() {
	if env.IsTST() {
		printLightRed("Running in test mode")
	} else if env.IsGTU() {
		printLightRed("Running in GTU mode")
	} // unspecified == PRD

	// Initialize logging
	logFile, err := os.OpenFile("mails-util-log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	if os.Getenv("LOG_TO_STDOUT") == "1" {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	} else {
		log.SetOutput(logFile)
	}

	log.Println("=== Invoking mails ===")

	// CLI
	app := &cli.App{
		Name:  "mails",
		Usage: "Mailinglijsten utility",
		Commands: []*cli.Command{
			{
				Name:  "send",
				Usage: "send an email to a ban immediately",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "from",
						Usage: "Afzender van de email. e.g. Marie Jo <marie.jo@ksadebiekorf.be>",
					},
					&cli.StringFlag{
						Name:  "reply-to",
						Usage: "Reply to naam + address",
					},
					&cli.IntSliceFlag{
						Name:    "ban",
						Aliases: []string{"b"},
						Usage:   "Definieer een ban(nen) waarnaar de email verzonden moet worden. e.g. mails send -b 1,5 verstuurt naar ban 1 en 5",
					},
					&cli.StringFlag{
						Name:    "subject",
						Aliases: []string{"s"},
						Usage:   "Het onderwerp van de mail",
					},
					&cli.StringFlag{
						Name:    "message",
						Aliases: []string{"m"},
						Usage:   "Het email bericht. In HTML vorm",
					},
					&cli.StringFlag{
						Name:    "content-type",
						Aliases: []string{"t"},
						Usage:   "Content type van het bericht (i.e. \"text/html\" of \"text/plain\")",
						Value:   "text/html",
					},
				},
				Action: func(ctx *cli.Context) error {
					// fmt.Printf("ban: %d\n", ctx.IntSlice("ban"))
					// fmt.Printf("subject: %s\n", ctx.String("subject"))
					// fmt.Printf("message: %s\n", ctx.String("message"))
					// fmt.Printf("from: %s\n", ctx.String("from"))
					// fmt.Printf("reply-to: %s\n", ctx.String("reply-to"))

					return sendEmail(
						ctx.String("from"),
						ctx.String("reply-to"),
						ctx.IntSlice("ban"),
						ctx.String("subject"),
						ctx.String("message"),
						ctx.String("content-type"),
					)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

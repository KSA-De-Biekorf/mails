package main

import (
	// "fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

// Mailing utility `mails`
func main() {
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
					)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

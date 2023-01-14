package main

import (
	"fmt"
	"log"

	"ksadebiekorf.be/mailing/db"
	"ksadebiekorf.be/mailing/env"
	"ksadebiekorf.be/mailing/mail"
)

func sendEmail(
	from,
	replyTo string,
	bannen []int,
	subject,
	message string,
) error {
	var err error
	var emails []db.EmailEntry
	if env.IsTest() {
		emails = []db.EmailEntry{
			{
				Email:     "john.appleseed@ksadebiekorf.be",
				FirstName: "John",
				LastName:  "Appleseed",
			}, {
				Email:     "john.doe@ksadebiekorf.be",
				FirstName: "John",
				LastName:  "Doe",
			},
		}
	} else {
		mailing_db, err := db.DbConnect()
		if err != nil {
			return err
		}

		// Fetch emails voor bannen
		for _, ban := range bannen {
			res, err := db.FetchBan(mailing_db, ban)
			if err != nil {
				return err
			}
			emails = append(emails, res...)
		}
	}

	// Send email
	cfg := mail.MailConfig{
		// From:        mail.ParseEmail(from),
		// ReplyTo:     mail.ParseEmail(replyTo),
		Subject:     subject,
		HTMLContent: message,
	}
	if replyTo != "" {
		cfg.ReplyTo = mail.ParseEmail(replyTo)
	}

	if from != "" {
		cfg.From = mail.ParseEmail(from)
	} else {
		var name string
		if cfg.ReplyTo == nil {
			name = "KSA De Biekorf"
		} else {
			name = cfg.ReplyTo.Name
		}
		cfg.From = mail.NewEmail(name, "no_reply@ksadebiekorf.be")
	}

	if env.IsPRD() || env.IsGTU() {
		for _, email := range emails {
			log.Printf("Sending email %v\n", email)
			cfg.Tos = []db.EmailEntry{email}
			req := mail.CreateMailRequest(&cfg)
			err = mail.SendEmail(req)
			if err != nil {
				return nil
			}
		}
	} else {
		fmt.Println("Sending email")
		fmt.Printf("From: %v\n", *cfg.From)
		fmt.Printf("ReplyTo: %v\n", *cfg.ReplyTo)
		fmt.Printf("%#v\n", cfg)
		fmt.Println(emails)
	}

	log.Println("> Emails sent successfully")

	return err
}

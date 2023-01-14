package main

import (
	"ksadebiekorf.be/mailing/db"
	"ksadebiekorf.be/mailing/mail"
	"log"
)

func sendEmail(
	from,
	replyTo string,
	bannen []int,
	subject,
	message string,
) error {
	var err error
	mailing_db, err := db.DbConnect()
	if err != nil {
		return err
	}

	// Fetch emails voor bannen
	var emails []db.EmailEntry
	for _, ban := range bannen {
		res, err := db.FetchBan(mailing_db, ban)
		if err != nil {
			return err
		}
		emails = append(emails, res...)
	}

	// // tmp
	// emails := []db.EmailEntry{
	// 	{
	// 		Email:     "everaert.jonas@outlook.com",
	// 		FirstName: "Jonas",
	// 		LastName:  "Everaert",
	// 	}, {
	// 		Email:     "info@jonaseveraert.be",
	// 		FirstName: "Jonas",
	// 		LastName:  "Everaert",
	// 	}, {
	// 		Email:     "jonas.vbs4@gmail.com",
	// 		FirstName: "Jonas",
	// 		LastName:  "Everaert",
	// 	},
	// }

	// Send email
	cfg := mail.MailConfig{
		From:        mail.ParseEmail(from),
		ReplyTo:     mail.ParseEmail(replyTo),
		Subject:     subject,
		HTMLContent: message,
	}
	for _, email := range emails {
		log.Printf("Sending email %v\n", email)
		cfg.Tos = []db.EmailEntry{email}
		req := mail.CreateMailRequest(&cfg)
		err = mail.SendEmail(req)
		if err != nil {
			return nil
		}
	}

	log.Println("> Emails sent successfully")

	return err
}

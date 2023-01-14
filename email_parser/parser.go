package email_parser

import (
	// "fmt"
	"bufio"
	"errors"
	"io"
	"log"
	"net/mail"
	"os"
	"strconv"
	"strings"
)

type Email struct {
	From        string
	Subject     string
	ContentType string
	Message     string
	SpamScore   int
	IsSpam      bool
}

func readPiped() (string, error) {
	// Check pipe
	fileInfo, _ := os.Stdin.Stat()
	if fileInfo.Mode()&os.ModeCharDevice != 0 {
		return "", errors.New("Nothing is being piped to the program")
	}

	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	piped := strings.Builder{}
	for scanner.Scan() {
		t := scanner.Text()
		log.Printf("Scanner text: %v\n", t)
		piped.WriteString(t)
		piped.WriteString("\n")
	}

	return piped.String(), nil
}

func EmailFromPipe() (*Email, error) {
	// var emailInput string
	// fmt.Scanf("%s", emailInput)

	emailInput, err := readPiped()
	if err != nil {
		return nil, err
	}

	inputReader := strings.NewReader(emailInput)
	msg, err := mail.ReadMessage(inputReader)
	if err != nil {
		return nil, err
	}

	log.Println("Received: ", msg)

	var email Email

	// Headers
	header := msg.Header
	email.From = header.Get("From")
	email.Subject = header.Get("Subject")
	contentType := header.Get("Content-Type")
	if strings.Contains(contentType, "plain") {
		email.ContentType = "text/plain"
	} else {
		email.ContentType = "text/html"
	}

	spamScoreS := header.Get("X-Spam-Score")
	spamScoreInt, err := strconv.Atoi(spamScoreS)
	if err == nil {
		email.SpamScore = spamScoreInt
	}
	email.IsSpam = header.Get("X-Is-Flag") == "YES"

	// Body
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return nil, err
	}
	email.Message = parseMessage(string(body))

	return &email, nil
}

// TODO: when sending HTML, remove apple stuff
// TODO: substitutions: %to%, ...
// TODO: attachementsè!!!
func parseMessage(msg string) string {
	return msg
}

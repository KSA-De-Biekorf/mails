package email_parser

import (
	// "fmt"
	"bufio"
	"errors"
	"io"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Email struct {
	From        string
	Subject     string
	ContentType string
	Messages    []SubMessage
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
	contentTypeParts := strings.Split(contentType, ";")
	contentTypeName := strings.Trim(contentTypeParts[0], " ")
	email.ContentType = contentTypeName
	var multipartBoundary *string = nil
	if strings.Contains(contentTypeName, "multipart") {
		r := regexp.MustCompile(`boundary=(.*)`)
		match := r.FindStringSubmatch(contentTypeParts[1])
		multipartBoundary = &match[1]
		email.ContentType = contentTypeName
	}

	spamScoreS := header.Get("X-Spam-Score")
	spamScoreInt, err := strconv.Atoi(spamScoreS)
	if err == nil {
		email.SpamScore = spamScoreInt
	}
	email.IsSpam = header.Get("X-Spam-Flag") == "YES"

	// Body
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		log.Printf("Error while reading body")
		return nil, err
	}
	email.Messages, err = parseMessage(string(body), email.ContentType, multipartBoundary)
	if err != nil {
		log.Printf("Error in `parseMessage`")
		return nil, err
	}

	return &email, nil
}

// TODO: check if attachements work
// TODO: substitutions: %to%, ...
func parseMessage(msg string, mainMessageContentType string, boundary *string) (msgs []SubMessage, err error) {
	if boundary == nil {
		// No additional content types need to be returned, because the main content type has already been pased
		return []SubMessage{{ContentType: mainMessageContentType, Message: msg}}, nil
	}

	// Parse multipart message
	boundRegex := regexp.MustCompile(`--` + *boundary)

	started := false
	var subMessage []string
	for _, str := range strings.Split(msg, "\n") {
		if boundRegex.MatchString(str) {
			if !started {
				started = true
				if len(subMessage) != 0 {
					msgs = append(msgs, SubMessage{ContentType: mainMessageContentType, Message: strings.Join(subMessage, "\n")})
					subMessage = subMessage[:0]
				}
			} else {
				subMsg, err := parseSubMessage(strings.Join(subMessage, "\n"))
				if err != nil {
					return nil, err
				}
				subMessage = subMessage[:0] // reset array, but keep the allocated memory

				msgs = append(msgs, *subMsg)
			}
		} else {
			subMessage = append(subMessage, str)
		}
	}

	return msgs, nil
}

type SubMessage struct {
	ContentType string
	Message     string
}

func parseSubMessage(msg string) (*SubMessage, error) {
	var subMessage SubMessage

	reader := strings.NewReader(msg)
	parsedMsg, err := mail.ReadMessage(reader)
	if err != nil {
		log.Printf("Error in `mail.ReadMessage` (in function `parseSubMessage`). `msg` was %s\n", msg)
		return nil, err
	}

	subMessage.ContentType = parsedMsg.Header.Get("Content-Type")
	msgBytes, err := io.ReadAll(parsedMsg.Body)
	if err != nil {
		return nil, err
	}
	subMessage.Message = string(msgBytes)
	return &subMessage, nil
}

__build:
	GOOS=linux go build -o build/mails -ldflags="-s -w" ./cli

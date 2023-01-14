# Mailinglijsten utility

## Running

```bash
go run ./cli
```

## Bulding

```bash
GOOS=linux go build -o build/mails -ldflags="-s -w" ./cli
```

or

```bash
make
```

## Example

```bash
mails send \
  --reply-to "Marie Jo <marie.jo@hotmail.com>" \
  --ban 1,5 \
  --subject "Waar ben ik?" \
  --message "Ik ben verdwenen.<br/>mvg,<br/>Marie Jo"
```

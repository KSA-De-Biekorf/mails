# Mailinglijsten utility

## Running

```bash
go run ./cli
```

## Bulding

```bash
go build -o build/mails ./cli
```

## Example

```bash
mails send \
  --reply-to "Marie Jo <marie.jo@hotmail.com>" \
  --ban 1,5 \
  --subject "Waar ben ik?" \
  --message "Ik ben verdwenen.<br/>mvg,<br/>Marie Jo"
```

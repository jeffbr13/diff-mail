# diff-mail

Emails you hourly diffs for a website

Watches a webpage, and uses the command-line `diff` tool to send changes in textual content to an email address.


Running
-------

The following environment variables must be set (fairly self-explanatory):

- `SMTP_HOST`
- `SMTP_USERNAME`
- `SMTP_PASSWORD`

```sh
go build main.go
./main -url="http://www.bbc.co.uk" -me="me@example.com"
```

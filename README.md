Simple golang SMTP server that receives incoming email and output to stdout in JSON format.

### Usage

* `make` or `make run` to run

* `make build` to build ARM 7 and host CPU/OS versions

The following flags can be set to customise the server:

```
Usage of ./smtpd:
  -max int
        max body size KB (default 256)
  -name string
        server name (default "mail-server")
  -port int
        listening port (default 25)
  -v    verbose logging
```

### JSON Output

Each message is printed to stdout as a JSON string terminated with a `\n`.
```
{
	"to": "email@domain.com",
	"from": "from@domain.com",
	"subject": "The subject",
	"body": "The body in text format"
	"html": "<html>...</html>"
}
```

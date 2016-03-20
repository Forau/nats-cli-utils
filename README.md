# nats-cli-utils
A small collection of command line tools for working with nats messaging system.
This is a work in progress...

To get all commands, run:
```
go get -v -u github.com/Forau/nats-cli-utils/...                                                                             
```

Tool        | Description                              | Status
----        | -----------                              | ----------
nats-tail   | Do tail on one or several nats subjects  | Beta
nats-srv    | Tool to quickly make adhoc services      | Experimental

nats-tail
---------
```
Usage: nats-tail <flags> subject

Arguments:
     subject: The NATS subject to listen on. Wildcards are supported.

Flags:
  -raw
        Short for template that just prints the data as it comes. Is equal to -t "{{.Data|printf "%s"}}"
  -templ string
        Template for output in golang text/template format with some added functions (default "{{time.Unix}}\t{{.Subject}}\t{{.Reply}}\n\t{{.Data | hex }}")
  -url string
        Uri to connect with NATS (default "nats://localhost:4222")
```

nats-srv
--------
Experimental state. Dont expect it to be useful yet.
Services are chained together, sort of like unix pipes. Currently just a few exists, but more will come.

```
Usage: nats-srv <flags> (subject) service <flags> <args> service <flags> <args>...

Flags:
  -raw
        If set, any service can be first in chain instead of 'nats-srv'
  -url string
        Connection url to pass to service 'nats-srv' (default "nats://localhost:4222")

Services:
nats-srv  <flags>  arg0
  func(chan interface {}) error
         -url   Connection to nats      (nats://localhost:4222)
         Arg0 Subject to subscribe to

echo
  func(interface {}) interface {}

teelog  <flags>
  func(interface {}) (interface {}, error)
         -file  Where to log stream to. Filename, or - which is stdout  (-)
         -format        Printf format for the value     (%!s(MISSING))
         -prefix        The prefix to put on log rows   ()

printf  arg0
  func(interface {}) (interface {}, error)
         Arg0 Format to print the value in

exec  <flags>  arg0
  func(interface {}) (interface {}, error)
         -args  The arguments for the command. Make sure it is 'one flag'. It will be split on spaces for the final arguments   ()
         Arg0 Command to execute

```

Example usages can be:
```
nats-srv subj printf "I got your %s" echo
```
Which will reply on calls to 'subj' with some appended text.

Other example can be:
```
nats-srv subj teelog -file log.log exec -args "www.google.com 80" nc
```
Which would log the request to log.log, then run execute netcat and send what was posted to subj to google, and return the result.




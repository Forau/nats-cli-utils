# nats-cli-utils
A small collection of command line tools for working with nats messaging system.
This is a work in progress...

To get all commands, run:
```
go get -v -u github.com/Forau/nats-cli-utils/...                                                                             
```

Tool        | Description
----        | -----------
nats-tail   | Do tail on one or several nats subjects


nats-tail
---------
```
Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -u, --uri="nats://localhost:4222"  
                       Uri to connect with NATS
  -s, --sub="debug.>"  The subject(s) to listen to
  -t, --templ="{{time.Unix}}\t{{.Subject}}\t{{.Reply}}\n\t{{.Data | hex }}"  
                       Template for output in golang text/template format with some added functions
  -r, --raw            Short for template that just prints the data as it comes. Is equal to -t "{{.Data|printf "%s"}}"
```


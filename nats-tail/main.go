// Copyright (c) 2016 Forau @ github.com. MIT License.

package main

import (
	"github.com/nats-io/nats"

	"flag"
	"fmt"

	"encoding/hex"
	"os"
	"runtime"
	"text/template"
	"time"
)

var (
	natsUri string
	templ   *template.Template
)

func init() {
	flag.StringVar(&natsUri, "url", nats.DefaultURL, "Uri to connect with NATS")
	templText := flag.String("templ", "{{time.Unix}}\t{{.Subject}}\t{{.Reply}}\n\t{{.Data | hex }}", "Template for output in golang text/template format with some added functions")
	raw := flag.Bool("raw", false, "Short for template that just prints the data as it comes. Is equal to -t \"{{.Data|printf \"%s\"}}\"")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: nats-tail <flags> subject\n\n")
		fmt.Fprint(os.Stderr, "Arguments:\n     subject: The NATS subject to listen on. Wildcards are supported.\n\n")
		fmt.Fprint(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	funcMap := template.FuncMap{
		"hex":  hex.Dump,
		"time": time.Now,
	}
	tmplText := *templText
	if *raw {
		*templText = "{{.Data|printf \"%s\"}}"
	}

	templ = template.Must(template.New("output").Funcs(funcMap).Parse(tmplText))
}

func printMsg(msg *nats.Msg) {
	err := templ.Execute(os.Stdout, msg)
	if err != nil {
		panic(err)
	}
}

func main() {
	nc, err := nats.Connect(natsUri)
	if err != nil {
		panic(err)
	}

	for _, s := range flag.Args() {
		_, err = nc.Subscribe(s, printMsg)
		if err != nil {
			panic(err)
		}
	}

	runtime.Goexit()
}

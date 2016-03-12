// Copyright (c) 2016 Forau @ github.com. MIT License.

package main

import (
	"github.com/nats-io/nats"
	"gopkg.in/alecthomas/kingpin.v2"

	"encoding/hex"
	"os"
	"runtime"
	"text/template"
	"time"
)

var (
	natsUri   = kingpin.Flag("uri", "Uri to connect with NATS").Default(nats.DefaultURL).Short('u').String()
	sub       = kingpin.Flag("sub", "The subject(s) to listen to").Default("debug.>").Short('s').String()
	templText = kingpin.Flag("templ", "Template for output in golang text/template format with some added functions").
			Default("{{time.Unix}}\t{{.Subject}}\t{{.Reply}}\n\t{{.Data | hex }}").
			Short('t').String()

	raw = kingpin.Flag("raw", "Short for template that just prints the data as it comes. Is equal to -t \"{{.Data|printf \"%s\"}}\"").
		Short('r').Bool()

	templ *template.Template
)

func init() {
	kingpin.Parse()

	funcMap := template.FuncMap{
		"hex":  hex.Dump,
		"time": time.Now,
	}
	tmplText := *templText
	if *raw {
		tmplText = "{{.Data|printf \"%s\"}}"
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
	nc, err := nats.Connect(*natsUri)
	if err != nil {
		panic(err)
	}

	_, err = nc.Subscribe(*sub, printMsg)
	if err != nil {
		panic(err)
	}

	runtime.Goexit()
}

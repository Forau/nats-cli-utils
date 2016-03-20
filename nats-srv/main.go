// Copyright (c) 2016 Forau @ github.com. MIT License.

package main

import (
	"log"

	//  "sync"
	//  "runtime"
	//  "io"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Forau/nats-cli-utils/services"
	"github.com/Forau/nats-cli-utils/services/natssrv"
	_ "github.com/Forau/nats-cli-utils/services/simple"
)

var (
	raw     = flag.Bool("raw", false, "If set, any service can be first in chain instead of 'nats-srv'")
	url     = flag.String("url", natssrv.DefaultURL, "Connection url to pass to service 'nats-srv'")
	service services.Service
)

func init() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}

	args := []string{}
	if !*raw {
		args = append(args, "nats-srv", "-url", *url)
	}
	args = append(args, flag.Args()...)
	srv, err := services.FindService(args...)
	if err != nil {
		log.Panic(err)
	}
	service = srv
}

func usage() {
	fmt.Fprint(os.Stderr, "Usage: nats-srv <flags> (subject) service <flags> <args> service <flags> <args>...\n\n")
	fmt.Fprint(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\nServices:\n")
	for _, srv := range services.ServiceRegister {
		fmt.Fprintf(os.Stderr, srv.Usage())
	}
}

func main() {
	end := make(chan interface{})
	for {
		res, err := service.Invoke(end)
		log.Printf("Res: %+v,%+v", res, err)
		if err == nil {
			break
		}
		<-time.After(time.Second)
	}
	<-end
}

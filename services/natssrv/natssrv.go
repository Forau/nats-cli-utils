package natssrv

import (
	"fmt"
	"log"

	"reflect"

	"github.com/Forau/nats-cli-utils/services"
	"github.com/nats-io/nats"
)

var DefaultURL = nats.DefaultURL
var ServiceList = []services.ServiceCreator{}

func regService(srv services.ServiceCreatorFn) {
	ServiceList = append(ServiceList, services.RegisterService(srv))
}

func init() {
	regService(func() (*services.ServiceConfig, interface{}) {
		fset := services.NewServiceConfig("nats-srv")
		natsURI := fset.Flag("url", nats.DefaultURL, "Connection to nats")
		subject := fset.Arg("Subject to subscribe to")

		return fset, func(end chan interface{}) error {
			nc, err := nats.Connect(*natsURI)
			if err != nil {
				return err
			}
			defer nc.Close()

			sub, err := nc.Subscribe(*subject, func(msg *nats.Msg) {
				res, _ := fset.Next.Invoke(msg.Data)
				log.Printf("Sent '%+v', and got '%+v'", msg.Data, res)
				switch restyp := res.(type) {
				case []byte:
					nc.Publish(msg.Reply, restyp)
				case string:
					nc.Publish(msg.Reply, []byte(restyp))
				default:
					log.Panic("No converter for ", restyp, " (", reflect.TypeOf(res), ")")
				}
			})
			fmt.Println("NATS service on subject: ", *subject, " :: ", sub)

			<-end
			return nil
		}
	})

}

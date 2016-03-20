package simple

import (
	"fmt"
	"github.com/Forau/nats-cli-utils/services"
	"io"
	"log"
	"strings"

	"os"
	"os/exec"
)

var ServiceList = []services.ServiceCreator{}

func regService(srv services.ServiceCreatorFn) {
	ServiceList = append(ServiceList, services.RegisterService(srv))
}

func init() {
	regService(func() (*services.ServiceConfig, interface{}) {
		fset := services.NewServiceConfig("echo")
		return fset, func(in interface{}) interface{} {
			return in
		}
	})

	regService(func() (*services.ServiceConfig, interface{}) {
		fset := services.NewServiceConfig("teelog")
		logf := fset.Flag("file", "-", "Where to log stream to. Filename, or - which is stdout")
		format := fset.Flag("format", "%s", "Printf format for the value")
		prefix := fset.Flag("prefix", "", "The prefix to put on log rows")

		var logger log.Logger
		fset.ValidateFn(func() error {
			var writer io.Writer
			switch *logf {
			case "-":
				writer = os.Stdout
			default:
				f, err := os.OpenFile(*logf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0666))
				if err != nil {
					log.Printf("Got error when opening file for logging: %+v\n", err)
					return err
				}
				writer = f
			}
			logger = *log.New(writer, *prefix, log.LstdFlags)
			return nil
		})

		return fset, func(in interface{}) (interface{}, error) {
			logger.Printf(*format, in)
			return fset.Next.Invoke(in)
		}
	})

	regService(func() (*services.ServiceConfig, interface{}) {
		fset := services.NewServiceConfig("printf")
		format := fset.Arg("Format to print the value in")
		return fset, func(in interface{}) (interface{}, error) {
			str := fmt.Sprintf(*format, in)
			return fset.Next.Invoke(str)
		}
	})

	regService(func() (*services.ServiceConfig, interface{}) {
		fset := services.NewServiceConfig("exec")
		args := fset.Flag("args", "", "The arguments for the command. Make sure it is 'one flag'. It will be split on spaces for the final arguments")
		cmdStr := fset.Arg("Command to execute")
		return fset, func(in interface{}) (interface{}, error) {
			argArr := strings.Split(*args, " ")
			var cmd *exec.Cmd
			log.Printf("CMD: '%s', argLen %d\n", *cmdStr, len(argArr))
			if len(argArr) == 0 || *args == "" {
				cmd = exec.Command(*cmdStr)
			} else {
				cmd = exec.Command(*cmdStr, argArr...)
			}
			str := fmt.Sprintf("%s", in)
			cmd.Stdin = strings.NewReader(str)

			buf, err := cmd.CombinedOutput()
			log.Printf("Exec command '%s' -> '%s', %+v\n", *cmd, string(buf), err)
			if err != nil {
				return nil, err
			}

			return string(buf), nil
		}
	})

}

package services

import (
	"fmt"

	"errors"
	"reflect"
)

// A service takes something in, and returns something back when invoked
type Service interface {
	Invoke(input interface{}) (interface{}, error)
}

// ServiceCreatorFn: Creates the flags for creating a service, and the service function.
// The use of flag.FlagSet is not final, and might be changed to another input parser
type ServiceCreatorFn func() (*ServiceConfig, interface{})

func (sc ServiceCreatorFn) Create(args ...string) (Service, error) {
	conf, fn := sc()
	rest, err := conf.ParseFlags(args...)
	if err != nil {
		return nil, err
	}

	if len(rest) > 0 {
		next, err := ServiceRegister.FindService(rest...)
		if err != nil {
			return nil, err
		}
		conf.Next = next
	}

	return &serviceHolder{serviceFn: fn, restArgs: rest}, nil
}

func (sc ServiceCreatorFn) Usage() string {
	conf, fn := sc()
	return conf.Usage(fn)
}

type ServiceCreator interface {
	Create(args ...string) (Service, error)
	//	Type() reflect.Type     // Type probably not needed
	Usage() string // Wait with for now
}

// Open for manipulation so the caller can replace services
var ServiceRegister = ServiceRegisterType{}

type ServiceRegisterType []ServiceCreator

func FindService(args ...string) (Service, error) {
	return ServiceRegister.FindService(args...)
}

func RegisterService(createFn ServiceCreatorFn) ServiceCreator {
	return ServiceRegister.Register(createFn)
}
func (srt *ServiceRegisterType) Register(createFn ServiceCreatorFn) ServiceCreator {
	*srt = append(*srt, createFn)
	return createFn
}

func (srt *ServiceRegisterType) FindService(args ...string) (Service, error) {
	for _, sv := range *srt {
		srv, err := sv.Create(args...)
		if err == nil {
			return srv, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Could not find service for '%+v'", args))
}

// A simple struct to implement the Service interface for arbitar function
type serviceHolder struct {
	serviceFn interface{}
	restArgs  []string
}

// Invokes the function. We might want to have checks and bound for this. Implements Service interface
func (bs *serviceHolder) Invoke(input interface{}) (interface{}, error) {
	var err error
	fn := reflect.ValueOf(bs.serviceFn)
	res := fn.Call([]reflect.Value{reflect.ValueOf(input)})
	if len(res) > 1 && !res[1].IsNil() {
		err = res[1].Interface().(error) // TODO: Fix better conversion
	}
	return res[0].Interface(), err
}

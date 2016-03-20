package services

import (
	"flag" // For now we use flag for parsing

	"bytes"
	"fmt"
	"reflect"

	"errors"
)

func NewServiceConfig(name string) *ServiceConfig {
	sc := &ServiceConfig{
		Name: name,

		fset: flag.NewFlagSet(name, flag.ContinueOnError),
	}

	return sc
}

type argHolder struct {
	usage string
	value *string
}

// A holder for service metadata
type ServiceConfig struct {
	Name string

	Next Service
	fset *flag.FlagSet
	args []argHolder

	validateFn func() error
}

// Keep it simple for now, and thus only string flag is supported
func (sc *ServiceConfig) Flag(name, def, usage string) *string {
	return sc.fset.String(name, def, usage)
}

func (sc *ServiceConfig) Arg(usage string) *string {
	val := new(string)
	ah := &argHolder{usage: usage, value: val}
	sc.args = append(sc.args, *ah)
	return val
}

// Parse flags. If anything is left, it is returned in rest, if it doesnt match us, an error is returned
func (sc *ServiceConfig) ParseFlags(args ...string) (rest []string, err error) {
	if len(args) < 1 || args[0] != sc.Name {
		return nil, errors.New("Arguments did not match " + sc.Name)
	}
	err = sc.fset.Parse(args[1:])
	if err == nil {
		rest = sc.fset.Args()
	}

	for _, a := range sc.args {
		if len(rest) < 1 {
			return nil, errors.New("Required argument not set: Usage:\n" + sc.Usage(nil))
		}
		*a.value = rest[0]
		rest = rest[1:]
	}

	if sc.validateFn != nil {
		if err := sc.validateFn(); err != nil {
			return nil, err
		}
	}

	return // rest is set if no err, and err is from Parse
}

func (sc *ServiceConfig) Usage(fn interface{}) string {
	var buf bytes.Buffer

	var flagbuf bytes.Buffer
	sc.fset.VisitAll(func(f *flag.Flag) {
		flagbuf.WriteString(fmt.Sprintf("\t -%s\t%s\t(%s)\n", f.Name, f.Usage, f.DefValue))
	})

	buf.WriteString(sc.Name)
	if flagbuf.Len() > 0 {
		buf.WriteString("  <flags>")
	}
	for i := 0; i < len(sc.args); i++ {
		buf.WriteString(fmt.Sprintf("  arg%d", i))
	}
	buf.WriteString("\n")
	if fn != nil {
		buf.WriteString("  ")
		buf.WriteString(reflect.TypeOf(fn).String())
		buf.WriteString("\n")
	}
	buf.Write(flagbuf.Bytes())

	for idx, a := range sc.args {
		buf.WriteString(fmt.Sprintf("\t Arg%d %s\n", idx, a.usage))
	}
	buf.WriteString("\n")
	return buf.String()
}

func (sc *ServiceConfig) ValidateFn(fn func() error) {
	sc.validateFn = fn
}

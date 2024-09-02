package core

import (
	"io"
)

type Executor interface {
	Execute(variables []string) (string, error)
	ExecuteToWriter(variables []string, writer io.Writer) error
}

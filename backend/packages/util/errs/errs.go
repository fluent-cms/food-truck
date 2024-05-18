package errs

import (
	"fmt"
	"path"
	"runtime"
)

func Err(err error) error {
	if err == nil {
		return nil
	}
	_, file, line, _ := runtime.Caller(1)
	_, fileName := path.Split(file)
	return fmt.Errorf("%s:%d %w", fileName, line, err)
}

func Errf(format string, args ...any) error {
	_, file, line, _ := runtime.Caller(1)
	_, fileName := path.Split(file)
	args = append([]any{fileName, line}, args...)
	return fmt.Errorf("%s:%d "+format, args...)
}

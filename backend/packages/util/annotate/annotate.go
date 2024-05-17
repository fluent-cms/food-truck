package annotate

import (
	"fmt"
	"path"
	"runtime"
)

// Error wraps the given error with the file name and line number.
// Returns nil if err is nil, so the following usage is safe:
//
//	return annotate.Error(funcThatReturnsError())
func Error(err error) error {
	if err == nil {
		return nil
	}
	// can't simply call Errorf because the call stack would be affected
	_, file, line, _ := runtime.Caller(1)
	_, fileName := path.Split(file)
	return fmt.Errorf("%s:%d %w", fileName, line, err)
}

// Errorf creates an error using fmt.Errorf on the given format string, but additionally adds the file name and line number.
// An error may be wrapped by providing the %w verb (see fmt.Errorf for more details). Example usages:
//
//	annotate.Errorf("not found")
//	annotate.Errorf("failed to get item %s, err: %w", id, err)
func Errorf(format string, args ...any) error {
	_, file, line, _ := runtime.Caller(1)
	_, fileName := path.Split(file)
	args = append([]any{fileName, line}, args...)
	return fmt.Errorf("%s:%d "+format, args...)
}

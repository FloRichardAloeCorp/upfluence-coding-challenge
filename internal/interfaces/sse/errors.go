package sse

import "fmt"

type InvalidStatusCodeError struct {
	Target  int
	Current int
}

func (e *InvalidStatusCodeError) Error() string {
	return fmt.Sprintf("invalid status code, want %d got %d", e.Target, e.Current)
}

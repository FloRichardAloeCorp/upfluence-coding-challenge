package logs

import "fmt"

type Field struct {
	Key   string
	Value string
}

func (f *Field) toJSON() string {
	return fmt.Sprintf(`"%s":"%s"`, f.Key, f.Value)
}

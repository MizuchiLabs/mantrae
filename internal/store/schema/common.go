package schema

import (
	"encoding/json"
	"fmt"
)

func scanJSON[T any](value any, out *T) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during Scan: %v", r)
		}
	}()

	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}
	return json.Unmarshal(bytes, out)
}

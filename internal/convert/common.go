// Convert package contains functions to convert between db and proto types
package convert

import (
	"encoding/json"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func SafeInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func SafeInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func SafeFloat(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

func SafeTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func UnmarshalStruct[T any](s *structpb.Struct) (*T, error) {
	// Marshal the proto Struct to JSON bytes
	data, err := s.MarshalJSON()
	if err != nil {
		return nil, err
	}

	// Unmarshal into your target struct
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func MarshalStruct[T any](s *T) (*structpb.Struct, error) {
	// Marshal the target struct to JSON bytes
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	// Unmarshal into your proto Struct
	var out structpb.Struct
	if err := out.UnmarshalJSON(data); err != nil {
		return nil, err
	}
	return &out, nil
}

package service

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func SafeBool(b *bool) *wrapperspb.BoolValue {
	if b == nil {
		return wrapperspb.Bool(false)
	}
	return wrapperspb.Bool(*b)
}

func SafeTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

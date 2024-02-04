package scalar

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MarshalProtobufTimestamp(t *timestamppb.Timestamp) graphql.Marshaler {
	if t == nil {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(t.AsTime().Format(time.RFC3339Nano)))
	})
}

func UnmarshalProtobufTimestamp(v interface{}) (*timestamppb.Timestamp, error) {
	if tmpStr, ok := v.(string); ok {
		asTime, err := time.Parse(time.RFC3339Nano, tmpStr)
		if err != nil {
			return nil, err
		}
		return timestamppb.New(asTime), nil
	}
	return nil, errors.New("time should be RFC3339Nano formatted string")
}

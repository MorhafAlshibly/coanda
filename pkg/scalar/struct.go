package scalar

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/MorhafAlshibly/coanda/pkg/conversion"
	"google.golang.org/protobuf/types/known/structpb"
)

func MarshalProtobufStruct(s *structpb.Struct) graphql.Marshaler {
	if s == nil {
		return graphql.Null
	}
	return graphql.WriterFunc(func(w io.Writer) {
		m, err := conversion.ProtobufStructToMap(s)
		if err != nil {
			panic(err)
		}
		err = json.NewEncoder(w).Encode(m)
		if err != nil {
			panic(err)
		}
	})
}

func UnmarshalProtobufStruct(v interface{}) (*structpb.Struct, error) {
	if m, ok := v.(map[string]interface{}); ok {
		pbStruct, err := conversion.MapToProtobufStruct(m)
		if err != nil {
			return nil, err
		}
		return pbStruct, nil
	}
	return nil, errors.New("Struct must be a map (JSON object)")
}

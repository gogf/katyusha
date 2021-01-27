package grpcutils

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var (
	jsonPbMarshaller = &jsonpb.Marshaler{
		EmitDefaults: true,
	}
)

// MarshalPbMessageToJsonString marshals protobuf message to json string.
func MarshalPbMessageToJsonString(v interface{}) (msg string) {
	var err error
	pb, ok := v.(proto.Message)
	if ok {
		msg, err = jsonPbMarshaller.MarshalToString(pb)
	}
	if err != nil || !ok {
		msg = fmt.Sprintf("%+v", v)
	}
	return msg
}

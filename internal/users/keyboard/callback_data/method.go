package callback_data

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
)

func (m QueryDataType) GetBase64ProtoString() string {
	bytes, _ := proto.Marshal(&m)
	base64str := base64.StdEncoding.EncodeToString(bytes)
	return base64str
}

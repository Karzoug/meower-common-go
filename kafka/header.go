//nolint:gosec // false positive on md5 - we use it only for hashing
package kafka

import (
	"crypto/md5"
	"fmt"

	"google.golang.org/protobuf/proto"
)

const MessageTypeHeaderKey = "fngpnt"

func MessageTypeHeaderValue(msg proto.Message) string {
	res := md5.Sum([]byte(proto.MessageName(msg)))
	return fmt.Sprintf("%x", res)
}

package msgcode

import (
	"encoding/binary"

	"google.golang.org/protobuf/proto"
)

var MsgEndian = binary.LittleEndian // 写入消息头时的大小端

// 编码客户端消息
func EncodeMsg(msg *Msg) ([]byte, error) {
	data := make([]byte, 4) // 预留四个四个字节写长度
	data, err := proto.MarshalOptions{}.MarshalAppend(data, msg)
	if err != nil {
		return nil, err
	}
	msgLen := len(data) - 4
	MsgEndian.PutUint32(data, uint32(msgLen))
	return data, nil
}

// 解码客户端消息
func DecodeMsg(data []byte) (*Msg, int, error) {
	// 消息是否收全了
	if len(data) < 4 {
		return nil, 0, nil
	}
	msgLen := MsgEndian.Uint32(data)
	if int(msgLen+4) > len(data) {
		return nil, 0, nil
	}

	// 解析消息
	msg := &Msg{}
	err := proto.Unmarshal(data[4:4+msgLen], msg)
	if err != nil {
		return nil, int(msgLen + 4), err
	}
	return msg, int(msgLen + 4), nil
}

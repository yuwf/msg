package msgcode

import (
	"encoding/binary"
	"errors"

	"google.golang.org/protobuf/proto"
)

var MsgEndian = binary.BigEndian // 写入消息头时的大小端

// 编码客户端消息
func EncodeMsg(msg *Msg) ([]byte, error) {
	if msg.Head == nil {
		return nil, errors.New("Head is nil")
	}
	if msg.BodyMsg == nil && msg.Body == nil {
		return nil, errors.New("BodyMsg and Body is nil")
	}

	data := make([]byte, 26) // 预留消息头长度
	if msg.BodyMsg != nil {
		var err error
		data, err = proto.MarshalOptions{}.MarshalAppend(data, msg.BodyMsg)
		if err != nil {
			return nil, err
		}
	} else {
		data = append(data, msg.Body...)
	}

	msg.Head.Size = uint32(len(data)) - 26
	MsgEndian.PutUint32(data, msg.Head.MsgID)
	MsgEndian.PutUint16(data[4:], msg.Head.Flag)
	MsgEndian.PutUint64(data[6:], msg.Head.Param)
	MsgEndian.PutUint64(data[14:], msg.Head.Seq)
	MsgEndian.PutUint32(data[22:], msg.Head.Size)
	return data, nil
}

// 解码客户端消息
func DecodeMsg(data []byte) (*Msg, int, error) {
	// 消息是否收全了
	if len(data) < 26 {
		return nil, 0, nil
	}
	// 先读取长度
	msgLen := MsgEndian.Uint32(data[22:])
	if int(26+msgLen) > len(data) {
		return nil, 0, nil
	}

	h := &Head{}
	h.MsgID = MsgEndian.Uint32(data)
	h.Flag = MsgEndian.Uint16(data[4:])
	h.Param = MsgEndian.Uint64(data[6:])
	h.Seq = MsgEndian.Uint64(data[14:])
	h.Size = msgLen

	// 解析消息
	msg := &Msg{
		Head: h,
		Body: data[26 : 26+msgLen],
	}
	return msg, int(26 + msgLen), nil
}

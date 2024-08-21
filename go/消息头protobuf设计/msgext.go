package msgcode

import (
	"fmt"
	reflect "reflect"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

// 获取消息ID
var GetMsgIDByType = func(msgType reflect.Type) string {
	msg, ok := reflect.New(msgType).Interface().(proto.Message) // 创建一个实例
	if ok {
		return string(msg.ProtoReflect().Descriptor().Name())
	}
	return ""
}
var GetMsgIDByObj = func(msg proto.Message) string {
	return string(msg.ProtoReflect().Descriptor().Name())
}

// Msg结构加入的自定义消息
func (m *Msg) MsgID() string {
	return m.Head.MsgID
}
func (m *Msg) MsgMarshal() ([]byte, error) {
	buf, err := EncodeMsg(m)
	return buf, err
}
func (m *Msg) BodyUnMarshal(msgType reflect.Type) (interface{}, error) {
	msg, ok := reflect.New(msgType).Interface().(proto.Message) // 创建一个实例
	if ok {
		err := proto.Unmarshal(m.Body, msg)
		if err == nil {
			m.BodyMsg = msg
			return msg, nil
		}
		return nil, err
	}
	return nil, fmt.Errorf("%s not LMsg", msgType.Name())
}

// 日志输出时调用
func (m *Msg) MarshalZerologObject(e *zerolog.Event) {
	e.Interface("Head", m.Head)
	if m.BodyMsg != nil {
		e.Interface("Body", m.BodyMsg)
	} else if m.Body != nil {
		e.Interface("BodyLen", len(m.Body))
	}
	if len(m.Ext) > 0 {
		e.Str("Body", m.Ext)
	}
}

func (h *Head) MarshalZerologObject(e *zerolog.Event) {
	e.Str("MsgID", h.MsgID)
	e.Int64("Seq", h.Seq)
	e.Int32("Flag", h.Flag)
	if h.Flag == int32(HeadFlag_RPCReq) || h.Flag == int32(HeadFlag_RPCResp) {
		e.Interface("Rpc", h.Rpc)
	} else if h.Flag == int32(HeadFlag_User) {
		e.Interface("User", h.User)
	}
}

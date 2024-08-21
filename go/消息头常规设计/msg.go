package msgcode

import (
	"errors"
	"fmt"
	reflect "reflect"
	"strconv"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

const (
	HeadFlag_Notify  uint16 = 0 // 正常的消息
	HeadFlag_RPCReq  uint16 = 1 // rpc请求 需要填充MsgHead.param，全局唯一
	HeadFlag_RPCResp uint16 = 2 // rpc返回 需要填充rpc请求时的MsgHead.param
	HeadFlag_User    uint16 = 3 // 转发用户的消息 需要填充MsgHead.param为用户ID
)

// 消息头
type Head struct {
	MsgID uint32 // 消息id，使用消息名本身来定义
	Flag  uint16 // 见 HeadFlag
	Param uint64 // 消息参数 根据HeadFlag来解析
	Seq   uint64 // 消息序列号，每发一个消息自增1，客户端登录成功后，服务器会分配一个开始序列号，后续会验证
	Size  uint32 // 消息大小 不包括头 编解码时填充
}

// 消息结构
type Msg struct {
	Head    *Head         // 消息头
	Body    []byte        // 消息体 protobuf编码 Recv时DecodeMsg填充 特殊的情况就是消息透传使用，编码时BodyMsg为空时就用这个
	BodyMsg proto.Message // 消息结构, Msg.BodyUnMarshal和发送时填充

	Ext string // Log时额外输出的内容
}

// Msg结构加入的自定义消息
func (m *Msg) MsgID() string {
	return strconv.FormatUint(uint64(m.Head.MsgID), 10)
}

func (m *Msg) MsgMarshal() ([]byte, error) {
	buf, err := EncodeMsg(m)
	return buf, err
}

func (m *Msg) BodyUnMarshal(msgType reflect.Type) (interface{}, error) {
	if m.Body == nil {
		return nil, errors.New("Body is nil")
	}
	msg, ok := reflect.New(msgType).Interface().(proto.Message) // 创建一个实例
	if ok {
		err := proto.Unmarshal(m.Body, msg)
		if err == nil {
			m.BodyMsg = msg
			return msg, nil
		}
		return nil, err
	}
	return nil, fmt.Errorf("%s is not proto.Message", msgType.Name())
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
	e.Uint32("MsgID", h.MsgID)
	e.Uint16("Flag", h.Flag)
	e.Uint64("Param", h.Param)
	e.Uint64("Seq", h.Seq)
	e.Uint32("Size", h.Size)
}

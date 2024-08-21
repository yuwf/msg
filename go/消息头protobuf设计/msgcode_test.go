package msgcode

import (
	"context"
	"hash/crc32"
	"testing"
	"time"

	"github.com/yuwf/gobase/dispatch"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func BenchmarkCode(b *testing.B) {
	req := TestMsg{
		Desc: "123",
	}
	head := Head{
		MsgID: GetMsgIDByObj(&req),
		Seq:   6553614,
		Flag:  int32(HeadFlag_Notify),
	}
	msg := Msg{
		Head: &head,
	}
	msg.Body, _ = proto.Marshal(&req)
	head.CheckSum = crc32.ChecksumIEEE(msg.Body)

	buf, err := EncodeMsg(&msg)
	if err != nil {
		log.Error().Err(err).Msg("....")
		return
	}
	log.Info().Int("len", len(buf)).Msg("EncodeMsg")
	msg2, _, err := DecodeMsg(buf)
	if err != nil {
		log.Error().Err(err).Msg("....")
		return
	}
	log.Info().Interface("msg2", msg2).Msg("DecodeMsg")
}

type Client[T any] struct {
	name T
}

func (c *Client[T]) SendMsg(msg interface{}) error {
	m, _ := msg.(*Msg)
	log.Info().Interface("Head", m.Head).Interface("Body", m.BodyMsg).Msg("SendMsg")
	return nil
}

type Server struct {
	dispatch.MsgDispatch[*Msg, Client[string]]
}

func (h *Server) onMsgHandle1(ctx context.Context, m *Msg, msg *TestMsg, t *Client[string]) {
}

func (h *Server) onMsgHandle2(ctx context.Context, msg *TestMsg, t *Client[string]) {

}

func (h *Server) onRPCHandle1(ctx context.Context, req *TestMsg, resp *TestMsg, t *Client[string]) {
	resp.Desc = "rpc resp"
	time.Sleep(time.Second * 10)
}

func BenchmarkRegister(b *testing.B) {
	s := &Server{}
	s.RegMsgID = GetMsgIDByType
	// 注销消息
	//s.RegMsg1(s.onMsgHandle1)
	//s.RegMsg2(s.onMsgHandle2)
	s.RegReqResp(s.onRPCHandle1)

	// 组织一个消息
	req := TestMsg{
		Desc: "456",
	}
	head := Head{
		MsgID: GetMsgIDByObj(&req),
	}
	msg := Msg{
		Head: &head,
	}
	msg.Body, _ = proto.Marshal(&req)

	// 消息调用
	t := &Client[string]{}
	s.Dispatch(context.TODO(), &msg, t)

	s.WaitAllMsgDone(time.Second * 30)
}

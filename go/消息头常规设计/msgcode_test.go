package msgcode

import (
	"context"
	"testing"
	"time"

	"github.com/yuwf/gobase/dispatch"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	_ "github.com/yuwf/gobase/log"
)

func BenchmarkCode(b *testing.B) {
	req := TestMsg{
		Desc: "123",
	}
	head := Head{
		MsgID: 123,
		Seq:   6553614,
		Flag:  HeadFlag_Notify,
	}
	msg := Msg{
		Head:    &head,
		BodyMsg: &req,
	}

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
	req2 := &TestMsg{}
	proto.Unmarshal(msg2.Body, req2)
	if req.Desc != req2.Desc {
		log.Info().Interface("msg2", msg2).Msg("Err")
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
	// 注册
	s.RegReqRespI("123", "123", s.onRPCHandle1)
	s.SendResp = func(ctx context.Context, req *Msg, c *Client[string], respid string, resp interface{}) {
		log.Info().Interface("Head", req.Head).Interface("Body", req.BodyMsg).Msg("SendRespMsg")
	}

	// 组织一个消息
	req := TestMsg{
		Desc: "456",
	}
	head := Head{
		MsgID: 123,
	}
	msg := Msg{
		Head:    &head,
		BodyMsg: &req,
	}
	msg.Body, _ = proto.Marshal(&req)

	// 消息调用
	t := &Client[string]{}
	s.Dispatch(context.TODO(), &msg, t)

	s.WaitAllMsgDone(time.Second * 30)
}

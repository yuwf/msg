
#ifndef _MSGHEAD_H_
#define _MSGHEAD_H_

#include <stdint.h>

enum MsgFlag
{
	MsgFlag_Notify  = 0,    // 正常的消息
	MsgFlag_RPCReq  = 1,    // rpc请求 需要填充MsgHead.param，全局唯一
	MsgFlag_RPCResp = 2,    // rpc返回 需要填充rpc请求时的MsgHead.param
	MsgFlag_User    = 3,    // 转发用户的消息 需要填充MsgHead.param为用户ID
	MsgFlag_Max
};


#pragma pack(2)

struct MsgHead
{
	uint32_t  msgId;       // 消息ID
	uint16_t  flag;        // 见 MsgFlag
	uint64_t  param;       // 根据 MsgFlag 使用
	uint64_t  seq;         // 消息序列号，每发一个消息自增1，客户端登录成功后，服务器会分配一个开始序列号，后续会验证
	uint32_t  size;        // 消息大小

	MsgHead()
	{
		memset(this, 0, sizeof(MsgHead));
	}
};

#pragma pack()


// 返回值： 解码成功返回消耗buf的长度 返回-1表示有错误数据解码失败
unsigned int DecodeDemo(const void* buf, unsigned int len)
{
	if (len < sizeof(MsgHead))
	{
		return 0;
	}
	// 消息头
	const MsgHead* pHead = (const MsgHead*)buf;
	buf = (const char*)buf + sizeof(MsgHead);
	len -= sizeof(MsgHead);

	// 判断数据大小是否够了
	if (pHead->size > len)
	{
		return 0;
	}

	// todo ...

	return pHead->size + sizeof(MsgHead);
}

#endif

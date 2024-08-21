cd /d %~dp0

..\tool\protoc.exe --plugin=protoc-gen-go=..\tool\protoc-gen-go.exe --go_out=. msg.proto
if %errorlevel%==0 (
    ..\tool\sed -i "s/\/\/ keep, do not delete/\n\n\t\/\/ 导出脚本自动加入，用于消息输出\n\tExt string        \/\/ Log时额外输出的内容\n\tBodyMsg interface{}  \/\/ 消息结构, Msg.BodyUnMarshal和发送时填充/g" msg.pb.go
    goto 0
)
pause
:0
exit
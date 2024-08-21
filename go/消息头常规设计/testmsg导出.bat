cd /d %~dp0

.\protoc.exe --plugin=protoc-gen-go=.\protoc-gen-go.exe --go_out=. testmsg.proto
if %errorlevel%==0 (
    exit
)
pause
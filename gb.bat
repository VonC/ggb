@echo off
setlocal
set GOPATH=%~dp0/deps
set GOBIN=%~dp0/bin
cd %~dp0
if not exist %~dp0\deps\src\github.com\VonC\ggb (
mkdir %~dp0\deps\src\github.com\VonC
mklink /J %~dp0\deps\src\github.com\VonC\ggb %~dp0
)
if not exist %~dp0\deps\src\github.com\spf13\pflag\.git (
	git submodule update --init
)
rem cd
rem set GO
mklink /J deps\src\github.com\VonC\ggb %~dp0 2>/NUL
%GOROOT%\bin\go.exe build ^
 && %GOROOT%\bin\go.exe test 1>NUL && %GOROOT%\bin\go.exe test -coverprofile=coverage.out %*|grep -v -e "^\.\.*$"|grep -v "^$"|grep -v "thus far" ^
 && %GOROOT%\bin\go.exe install -a .
endlocal
doskey ggb=%~dp0\bin\ggb.exe $*
gc|grep -v main.go

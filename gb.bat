@echo off
setlocal
set GOPATH=%~dp0/deps
set GOBIN=%~dp0/bin
cd %~dp0
rem cd
rem set GO
mklink /J deps\src\github.com\VonC\ggb %~dp0 2>/NUL
%GOROOT%\bin\go.exe install -a .
endlocal
doskey ggb=%~dp0\bin\ggb.exe $*

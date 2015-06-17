@echo off
setlocal
set GOPATH=%~dp0/deps
set GOBIN=%~dp0/bin
cd %~dp0
rem cd
rem set GO
%GOROOT%\bin\go.exe install -a .
endlocal
doskey ggb=%~dp0\bin\ggb.exe $*

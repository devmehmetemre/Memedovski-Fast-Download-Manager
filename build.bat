@echo off
chcp 65001 >nul
echo Building Memedovski Fast Download Manager...
echo.

set GOROOT=%~dp0go
set PATH=%GOROOT%\bin;%PATH%
set GO111MODULE=on
set CGO_ENABLED=1
set CC=gcc

echo [1/3] Creating manifest...
"%USERPROFILE%\go\bin\rsrc.exe" -manifest app.manifest -o rsrc.syso -arch amd64

echo [2/3] Compiling...
go build -ldflags="-H windowsgui" -o FastDownloader.exe .

echo [3/3] Cleaning up...
del /q rsrc.syso 2>nul

echo.
echo Done: FastDownloader.exe
pause

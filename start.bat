@echo off
rem 说明。
chcp 65001 >nul 2>&1
rem 说明。
rem
rem 说明。
rem 说明。
rem 说明。
rem
rem 说明。
rem
rem 说明。
rem 说明。
rem 说明。
rem
rem 说明。
rem 说明。

setlocal enabledelayedexpansion
cd /d "%~dp0"

set "BIN=artex.exe"
if not exist "%BIN%" (
	echo [artex] 找不到可执行文件: %BIN% 1>&2
	exit /b 1
)

set "RESTART_CODE=75"
set "MAX_DELAY=60"
set /a delay=1

:loop
"%BIN%" %*
set "code=!ERRORLEVEL!"

if "!code!"=="0" (
	echo [artex] 正常退出
	exit /b 0
)

if "!code!"=="%RESTART_CODE%" (
	rem 说明。
	echo [artex] 收到重启请求：应用新版本…
	set /a delay=1
	goto loop
)

echo [artex] 异常退出 ^(code=!code!^), !delay!s 秒后重启 1>&2
rem 说明。
set /a pings=!delay!+1
ping -n !pings! 127.0.0.1 >nul 2>&1
set /a delay=!delay!*2
if !delay! gtr %MAX_DELAY% set /a delay=%MAX_DELAY%
goto loop

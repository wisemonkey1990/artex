@echo off
rem 把控制台切换到 UTF-8，否则本文件里的中文在 GBK 终端会显示乱码。
chcp 65001 >nul 2>&1
rem ARTEX 守护启动脚本（Windows）
rem
rem 用法：
rem   start.bat                  前台运行（Ctrl-C 停止）
rem   start.bat -addr :9000      额外参数会原样传给 artex
rem
rem 本脚本只做一件事：运行 artex.exe，进程退出后根据退出码决定是否重新拉起。
rem
rem   0      用户正常停止  -> 结束循环
rem   75     程序请求重启   -> 立即重新运行（界面上点了"一键更新"或"回滚"）
rem   其他   异常退出       -> 退避后重新运行（1->2->4…最多 60 秒）
rem
rem 下载·SHA256 校验·版本替换不在这里处理，全部由 artex 启动时自行完成
rem （selfupdate 包）。脚本保持简单。详见 start.sh 开头的说明。

setlocal enabledelayedexpansion
cd /d "%~dp0"

set "BIN=artex.exe"
if not exist "%BIN%" (
	echo [artex] 找不到可执行文件：%BIN% 1>&2
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
	rem 更新/回滚已就位。重新运行后 artex 会在启动时完成版本替换。
	echo [artex] 请求重启：应用新版本…
	set /a delay=1
	goto loop
)

echo [artex] 异常退出 ^(code=!code!^)，!delay!s 后重启 1>&2
rem timeout 在被重定向的控制台会失败，所以用 ping 代替（延迟 N 秒需要 N+1 次）。
set /a pings=!delay!+1
ping -n !pings! 127.0.0.1 >nul 2>&1
set /a delay=!delay!*2
if !delay! gtr %MAX_DELAY% set /a delay=%MAX_DELAY%
goto loop

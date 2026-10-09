#!/bin/sh
# 说明。
#
# 说明。
# 说明。
# 说明。
# 说明。
#
# 说明。
#
# 说明。
# 说明。
# 说明。
#
# 说明。
# 说明。
# 说明。
# 说明。
set -u

cd "$(dirname "$0")" || exit 1

BIN=./artex
[ -x "$BIN" ] || { echo "[artex] 找不到可执行文件：$BIN" >&2; exit 1; }

RESTART_CODE=75
MAX_DELAY=60

child=0
stopping=0

# 说明。
#
# 说明。
# 说明。
# 说明。
forward() {
	stopping=1
	if [ "$child" -ne 0 ]; then
		kill -TERM "$child" 2>/dev/null || true
	fi
}
trap forward INT TERM

delay=1
while :; do
	"$BIN" "$@" &
	child=$!

	# 说明。
	# 说明。
	wait "$child"
	code=$?
	if [ "$code" -gt 128 ]; then
		wait "$child"
		code=$?
	fi
	child=0

	if [ "$stopping" -eq 1 ]; then
		echo "[artex] 已停止"
		exit 0
	fi

	case "$code" in
		0)
			echo "[artex] 正常退出"
			exit 0
			;;
		"$RESTART_CODE")
			# 说明。
			echo "[artex] 收到重启请求（应用新版本）…"
			delay=1
			;;
		*)
			echo "[artex] 异常退出（code=$code），${delay}s 后重启" >&2
			sleep "$delay"
			delay=$((delay * 2))
			[ "$delay" -gt "$MAX_DELAY" ] && delay=$MAX_DELAY
			;;
	esac
done

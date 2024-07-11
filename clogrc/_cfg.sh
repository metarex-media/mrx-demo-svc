#!/usr/bin/env bash
#          _            _                   _
#   __ _  | |_    ___  | |__   __ _   _ _  | |__
#  / _` | | ' \  |___| | '_ \ / _` | | '_| | / /
#  \__, | |_||_|       |_.__/ \__,_| |_|   |_\_\
#  |___/
# shellcheck source=/dev/null
[ -n "$(clog -v 2>/dev/null)" ] && source <(clog Inc)
export PROJECT=mrx-demo-svc
export bEXE="mrx-demo-svc"
export callingSCRIPT="${0##*/}"
export vCodeType="golang"
export vCodeSrc="releases.yaml"
# export & assign separately: don't mask errors in $() that crash  cicd jobs
vCODE=$(cat $vCodeSrc | grep version | head -1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
bMSG=$(cat $vCodeSrc | grep note | head -1 | sed -nr "s/note: (.*)/\1/p" | xargs)
export vCODE
export bMSG

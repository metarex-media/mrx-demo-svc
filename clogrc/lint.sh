#!/usr/bin/env bash
#  clog> lint
# short> Use megalinter to inspect the code
# long>  use when tweaking a build to pass all the minor tests
#        _                  _   _          _
#   __  | |  ___   __ _    | | (_)  _ _   | |_
#  / _| | | / _ \ / _` |   | | | | | ' \  |  _|
#  \__| |_| \___/ \__, |   |_| |_| |_||_|  \__|
#                 |___/
#-----------------------------------------------------------------------------
[ -f clogrc/_cfg.sh ] && source clogrc/_cfg.sh

if [ -n "$(which mga-linter-runner)" ]; then
  mega-linter-runner
else
  echo "megalinter not found, install node & docker, then try:"
  echo "  npm install mega-linter-runner -g"
fi
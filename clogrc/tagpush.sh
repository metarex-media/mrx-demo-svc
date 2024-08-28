#!/usr/bin/env bash
#  clog> tagpush
# short> Re-Tag local HEAD & remote to the right version
# long>  use when tweaking a build to pass all the minor tests
#        _                  _                                      _
#   __  | |  ___   __ _    | |_   __ _   __ _   _ __   _  _   ___ | |_
#  / _| | | / _ \ / _` |   |  _| / _` | / _` | | '_ \ | || | (_-< | ' \
#  \__| |_| \___/ \__, |    \__| \__,_| \__, | | .__/  \_,_| /__/ |_||_|
#                 |___/                 |___/  |_|
# shellcheck disable=SC2154
#-----------------------------------------------------------------------------
[ -f clogrc/_cfg.sh ] && source clogrc/_cfg.sh

fInfo "Tag & Push to version$cW $vCODE$cT $vCodeType ($cW $vCodeSrc $cT)$cX"

CommitMSG="$(git log -1 --pretty=%B)"
fInfo "Retagging to$cW $vCODE$cT ($cS$CommitMSG$cT)...$cX"

git tag -d "$vCODE"
git tag -a "$vCODE" HEAD -m "$CommitMSG"

if mycmd; then
	fOk "HEAD is now$cW $vCODE$cT ($cS$CommitMSG$cT)$cX"
fi

git push --delete origin "$vCODE" 2>/dev/null
if [ -n "$vCODEmajor" ]; then
	# also tag the major version to this push
	git tag -d "$vCODEmajor"
	git tag -a "$vCODEmajor" HEAD -m "$CommitMSG"
	if mycmd; then
		fOk "HEAD is now also$cW $vCODEmajor$cT ($cS$CommitMSG$cT)$cX"1
	fi

	git push --delete origin "$vCODEmajor" 2>/dev/null
fi
# push with all tags
git push origin --follow-tags

if mycmd; then
	fOk "Remote is in sync with local$cX"
fi

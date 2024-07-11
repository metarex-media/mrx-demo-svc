#!/usr/bin/env bash
#  clog> tagpush
# short> Re-Tag local HEAD & remote to the right version
# long>  use when tweaking a build to pass all the minor tests
#          _            _                   _
#   __ _  | |_    ___  | |__   __ _   _ _  | |__
#  / _` | | ' \  |___| | '_ \ / _` | | '_| | / /
#  \__, | |_||_|       |_.__/ \__,_| |_|   |_\_\
#  |___/
# shellcheck disable=SC2154,
#-----------------------------------------------------------------------------
[ -f clogrc/_cfg.sh ] && source clogrc/_cfg.sh

fInfo "Tag & Push to version$cW $vCODE$cT $vCodeType ($cW $vCodeSrc $cT)$cX"

CommitMSG="$(git log -1 --pretty=%B)"
fInfo "Retagging to$cW $vCODE$cT ($cS$CommitMSG$cT)...$cX"

git tag -d "$vCODE"
git tag -a "$vCODE" HEAD -m "$CommitMSG"
[ $? -eq 0 ] && fOk "HEAD is now$cW $vCODE$cT ($cS$CommitMSG$cT)$cX"
git push --delete origin "$vCODE" 2>/dev/null
git push origin --follow-tags
[ $? -eq 0 ] && fOk "Remote is in sync with local$cX"
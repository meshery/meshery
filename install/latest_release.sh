#!/bin/bash
# GIT_REF=`git symbolic-ref HEAD`
# if [[ $GIT_REF = refs/tags* ]]
# then
# 	RELEASE_CHANNEL="stable"
# 	# export RELEASE_CHANNEL="stable"
# else
# 	RELEASE_CHANNEL="edge"
#  	# export RELEASE_CHANNEL="edge"
#  fi
# # echo "Release channel determined to be $RELEASE_CHANNEL"
# LATEST_VERSION=$(git describe --tags `git rev-list --tags --max-count=1`)
GIT_REF=`git symbolic-ref HEAD`
if [[ $GIT_REF = refs/tags* ]]
then
RELEASE_CHANNEL="stable"
# export RELEASE_CHANNEL="stable"
else
RELEASE_CHANNEL="edge"
# export RELEASE_CHANNEL="edge"
fi
echo "Release channel determined to be $RELEASE_CHANNEL"
echo LATEST_VERSION=$(git describe --tags `git rev-list --tags --max-count=1` --always)
echo GIT_VERSION=$(git describe --tags `git rev-list --tags --max-count=1` --always)
echo GIT_STRIPPED_VERSION=$(git describe --tags `git rev-list --tags --max-count=1` --always | cut -c 2-)
# env
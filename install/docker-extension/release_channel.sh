#!/bin/bash
GIT_REF=`git symbolic-ref HEAD`
if [[ $GIT_REF = refs/tags* ]]
then
	echo "stable"
	# export RELEASE_CHANNEL="stable"
else
	echo "edge"
 	# export RELEASE_CHANNEL="edge"
 fi
#  MESHERY_VERSION=$(cd ../../ && git describe --tags `git rev-list --tags --max-count=1`)
# echo $MESHERY_VERSION
# echo "Release channel determined to be $RELEASE_CHANNEL"

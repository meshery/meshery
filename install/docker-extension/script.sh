#!/bin/bash
	echo "RELEASE_CHANNEL=edge" 
GIT_REF=`git symbolic-ref HEAD`
( "refs" = refs/tags* ) ? "stable" : "edge"
	MESHERY_VERSION=$(cd ../../ && git describe --tags `git rev-list --tags --max-count=1`)
	echo $MESHERY_VERSION
echo "Release channel determined to be $RELEASE_CHANNEL"

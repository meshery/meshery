# Copyright 2019 Layer5 Authors
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

#!/usr/bin/env bash

OSARCHITECTURE="$(uname -m)"
OS="$(uname)"


if [ "x${OS}" = "xDarwin" ] ; then
  OSEXT="Darwin"
else
  OSEXT="Linux"
fi


if [ "x${MESHERY_VERSION}" = "x" ] ; then
  MESHERY_VERSION=$(curl -L -s https://api.github.com/repos/layer5io/meshery/releases | \
                  grep tag_name | sed "s/ *\"tag_name\": *\"\\(.*\\)\",*/\\1/" | \
                  grep -v "rc\.[0-9]$" | sort -t"." -k 1,1 -k 2,2 -k 3,3 -k 4,4 | tail -n 1)
fi

if [ "x${MESHERY_VERSION}" = "x" ] ; then
  printf "Unable to get latest mesheryctl version. Set MESHERY_VERSION env var and re-run. For example: export MESHERY_VERSION=v0.1.6"
  exit;
fi


NAME="mesheryctl-$MESHERY_VERSION"
URL="https://github.com/layer5io/meshery/releases/download/${MESHERY_VERSION}/mesheryctl_${MESHERY_VERSION:1}_${OSEXT}_${OSARCHITECTURE}.zip"
printf "Downloading %s from %s ...\n" "$NAME" "$URL"
#curl -L "$URL" | tar xz

curl -L ${URL} -o ${PWD}/meshery.zip
unzip ${PWD}/meshery.zip 
mv ${PWD}/mesheryctl /usr/local/bin/mesheryctl

printf ""
printf "Mesheryctl %s Download Complete!\n" "$MESHERY_VERSION"
printf "\n"
printf "Mesheryctl has been successfully downloaded into the /usr/local/bin folder on your system.\n"
#printf "\n"
rm -rf meshery.zip LICENSE README.md
printf "cleanup done!"

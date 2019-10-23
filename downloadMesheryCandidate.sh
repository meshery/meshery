#!/usr/bin/sh

OSARCHITECTURE="$(uname -m)"
OS="$(uname)"
if [ "x${OS}" = "xDarwin" ] ; then
  OSEXT="osx"
else
  OSEXT="linux"
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


#mesheryctl_0.1.6_Darwin_i386.zip
#echo ${A:1}


#NAME="mesheryctl-$MESHERY_VERSION"
URL="https://github.com/layer5io/meshery/releases/download/${MESHERY_VERSION}/mesheryctl_${MESHERY_VERSION:1}_${OSEXT}_${OSARCHITECTURE}.zip"
#URL="https://github.com/layer5io/meshery/releases/download/v0.1.6/mesheryctl_0.1.6_Darwin_i386.zip"
#printf "Downloading %s from %s ..." "$NAME" "$URL"
#curl -L "$URL" | tar xz

curl -L ${URL} -o /tmp/meshery.zip
unzip /tmp/meshery.zip 
cp ${PWD}/mesheryctl /usr/local/bin/mesheryctl

echo $URL
echo $MESHERY_VERSION
echo $OSARCHITECTURE
echo $OSEXT
#printf ""
#printf "Mesheryctl %s Download Complete!\n" "$MESHERY_VERSION"
#printf "\n"
#printf "Meshery has been successfully downloaded into the %s folder on your system.\n" "$NAME"
#printf "\n"
#BINDIR="$(cd "$NAME/bin" && pwd)"
#printf "Next Steps:\n"
#printf "See https://istio.io/docs/setup/kubernetes/install/ to add Istio to your Kubernetes cluster.\n"
#printf "\n"
#printf "To configure the istioctl client tool for your workstation,\n"
#printf "add the %s directory to your environment path variable with:\n" "$BINDIR"
#printf "\t export PATH=\"\$PATH:%s\"\n" "$BINDIR"
#printf "\n"

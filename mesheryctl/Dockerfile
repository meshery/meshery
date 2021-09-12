FROM golang:1.16

RUN adduser --disabled-login appuser
WORKDIR /github.com/meshery/meshery
ADD . .

RUN go clean -modcache; cd mesheryctl/cmd/mesheryctl/; GOPROXY=https://proxy.golang.org GOSUMDB=off go build -ldflags="-w -s -X main.version=$GIT_VERSION -X main.commitsha=$GIT_COMMITSHA -X main.releasechannel=$RELEASE_CHANNEL" -tags draft -a -o mesheryctl .
RUN mv mesheryctl/cmd/mesheryctl/mesheryctl /home/appuser/.
RUN yes | /home/appuser/mesheryctl version

ENTRYPOINT ["/home/appuser/mesheryctl"]

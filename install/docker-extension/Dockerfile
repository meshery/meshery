FROM golang:1.19-alpine AS builder
ENV CGO_ENABLED=0
RUN apk add --update make
WORKDIR /backend
COPY go.* .
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  make bin

FROM node:14.17-alpine3.13 AS client-builder
WORKDIR /ui
# cache packages in layer
COPY ui/package.json /ui/package.json
COPY ui/package-lock.json /ui/package-lock.json
ARG TARGETARCH
RUN --mount=type=cache,target=/usr/src/app/.npm-${TARGETARCH} \
  npm set cache /usr/src/app/.npm-${TARGETARCH} && \
  npm ci
# install
COPY ui /ui
RUN npm run build

FROM alpine
ARG GIT_STRIPPED_VERSION
ARG GIT_VERSION
ARG RELEASE_CHANNEL
LABEL org.opencontainers.image.title="Meshery" \
  org.opencontainers.image.description="Design and operate your cloud native deployments with the extensible management plane, Meshery." \
  org.opencontainers.image.vendor="Layer5, Inc." \
  com.docker.desktop.extension.api.version=">= 0.2.0" \
   com.docker.extension.screenshots="[ \
      { \
        \"alt\": \"Meshery Docker Extension\", \
        \"url\": \"https://raw.githubusercontent.com/meshery/meshery/master/install/docker-extension/docs/img/meshmap-docker-extension-for-meshery.png\" \
      },{ \
        \"alt\": \"Meshery Docker Extension\", \
        \"url\": \"https://raw.githubusercontent.com/meshery/meshery/master/install/docker-extension/docs/img/Docker-extension-meshery.png\" \
      } \
    ]" \
  com.docker.extension.detailed-description="\
  <h2>MeshMap is the world's only visual designer for Kubernetes and service mesh deployments.</h2> \
  <p>The Meshery Docker Extension extends Docker Desktop’s position as the cloud native developer’s go-to Kubernetes environment with easy access to the next layer of cloud native infrastructure: service meshes.</p> \
  <ul> \
    <li><b>Detection of Kubernetes environments -</b> scan your kubeconfig and select your current Kubernetes environment. Switch from one environment to another one.</li> \
    <li><b>Service mesh support for your Docker Compose apps -</b> Import your Docker Compose apps. Configure and deploy them to Kubernetes and any service mesh.</li> \
    <li><b>Visual designer for Docker Compose apps -</b> Early access to the Docker Extension for Meshery that offers a visual topology for designing Docker Compose applications, operating Kubernetes, service meshes, and their workloads.</li> \
    <li><b>Single-click deployment of any service mesh -</b> Support of 10 different service meshes to the fingertips of developers in connection with Docker Desktop’s ability to deliver Kubernetes locally.</li> \
  </ul>" \
  com.docker.desktop.extension.icon="https://raw.githubusercontent.com/meshery/meshery/master/install/docker-extension/docs/img/meshery-logo-light.svg" \
  com.docker.extension.publisher-url="https://meshery.io" \
  com.docker.extension.additional-urls="[{\"title\":\"Documentation\",\"url\":\"https://docs.meshery.io\"},{\"title\":\"Project\",\"url\":\"https://layer5.io/meshery\"},{\"title\":\"Slack\",\"url\":\"https://slack.layer5.io\"},{\"title\":\"Discussion Forum\",\"url\":\"https://discuss.layer5.io\"}]"
COPY --from=builder /backend/bin/service /
COPY docker-compose.yaml .
COPY metadata.json .
RUN mkdir -p /mesheryctl/darwin && mkdir -p /mesheryctl/windows && mkdir -p /mesheryctl/linux
ADD https://github.com/meshery/meshery/releases/download/${GIT_VERSION}/mesheryctl_${GIT_STRIPPED_VERSION}_Darwin_x86_64.tar.gz /mesheryctl/darwin/
ADD https://github.com/meshery/meshery/releases/download/${GIT_VERSION}/mesheryctl_${GIT_STRIPPED_VERSION}_Windows_x86_64.zip /mesheryctl/windows/
ADD https://github.com/meshery/meshery/releases/download/${GIT_VERSION}/mesheryctl_${GIT_STRIPPED_VERSION}_Linux_x86_64.tar.gz /mesheryctl/linux/
RUN tar -xzf /mesheryctl/darwin/mesheryctl_${GIT_STRIPPED_VERSION}_Darwin_x86_64.tar.gz -C /mesheryctl/darwin/
RUN unzip /mesheryctl/windows/mesheryctl_${GIT_STRIPPED_VERSION}_Windows_x86_64.zip -d /mesheryctl/windows/
RUN tar -xzf /mesheryctl/linux/mesheryctl_${GIT_STRIPPED_VERSION}_Linux_x86_64.tar.gz -C /mesheryctl/linux/
COPY meshery-logo-light.svg .
COPY --from=client-builder /ui/build ui
EXPOSE 7877/tcp
CMD ./service

ARG GO_VERSION=1.25
ARG HTMLTEST_VERSION=0.17.0

FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /src
RUN apk --update add nodejs npm git gcompat  
RUN npm -v

FROM base AS node
COPY package*.json .
ENV NODE_ENV=production
RUN npm install --prefix /src/  package.json 
RUN ls -li
#RUN npm install --loglevel=verbose; echo "exit=$?"
#RUN find / -maxdepth 4 -type d -name node_modules -print || true


FROM base AS hugo
ARG HUGO_VERSION=0.146.0
ARG TARGETARCH
WORKDIR /tmp/hugo
RUN wget -O "hugo.tar.gz" "https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-${TARGETARCH}.tar.gz"
RUN tar -xf "hugo.tar.gz" hugo

FROM base AS build-base
COPY --from=hugo /tmp/hugo/hugo /bin/hugo
COPY --from=node /src/node_modules /src/node_modules
COPY . .

FROM build-base AS dev

FROM build-base AS build
ARG HUGO_ENV
ARG DOCS_URL
RUN hugo --gc --minify -d /out -e $HUGO_ENV -b $DOCS_URL

FROM scratch AS release
COPY --from=build /out /

FROM scratch AS update-stats
COPY --from=build /src/hugo_stats.json /hugo_stats.json

FROM build AS validate-stats
RUN <<EOF
if [ -n "$(git status --porcelain -- hugo_stats.json)" ]; then
  echo >&2 'ERROR: hugo_stats.json differs. Update with `docker buildx bake update-stats`'
  exit 1
fi
EOF

FROM davidanson/markdownlint-cli2:v0.12.1 AS lint
USER root
RUN --mount=type=bind,target=. \
    /usr/local/bin/markdownlint-cli2 \
    "content/**/*.md" \
    "#content/engine/release-notes/*.md" \
    "#content/desktop/previous-versions/*.md"

FROM wjdp/htmltest:v${HTMLTEST_VERSION} AS test
WORKDIR /test
COPY --from=build /out ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest

FROM build-base AS update-modules
ARG MODULE
RUN <<"EOT"
set -ex
if [ -n "$MODULE" ]; then
    hugo mod get ${MODULE}
    RESOLVED=$(cat go.mod | grep -m 1 "${MODULE/@*/}" | awk '{print $1 "@" $2}')
    go mod edit -replace "${MODULE/@*/}=${RESOLVED}";
else
    echo "no module set";
fi
EOT
RUN hugo mod vendor

FROM scratch AS vendor
COPY --from=update-modules /src/_vendor /_vendor
COPY --from=update-modules /src/go.* /

FROM build-base AS build-upstream
ARG UPSTREAM_MODULE_NAME
ARG UPSTREAM_REPO
ARG UPSTREAM_COMMIT
ENV HUGO_MODULE_REPLACEMENTS="github.com/${UPSTREAM_MODULE_NAME} -> github.com/${UPSTREAM_REPO} ${UPSTREAM_COMMIT}"
RUN hugo --ignoreVendorPaths "github.com/${UPSTREAM_MODULE_NAME}" -d /out

FROM wjdp/htmltest:v${HTMLTEST_VERSION} AS validate-upstream
WORKDIR /test
COPY --from=build-upstream /out ./public
ADD .htmltest.yml .htmltest.yml
RUN htmltest

FROM dev


ARG GO_VERSION=1.25
ARG HUGO_VERSION=0.146.0
ARG HTMLTEST_VERSION=0.17.0

FROM golang:${GO_VERSION} AS base
WORKDIR /src
COPY Makefile ./
RUN make docker-install-base-deps

FROM base AS node
COPY package.json package-lock.json ./
RUN make docker-deps

FROM base AS hugo
ARG HUGO_VERSION
ARG TARGETARCH
RUN make docker-download-hugo HUGO_VERSION=${HUGO_VERSION} TARGETARCH=${TARGETARCH}

FROM base AS build-base
COPY --from=hugo /tmp/hugo/hugo /usr/local/bin/hugo
COPY --from=node /src/node_modules /src/node_modules
COPY . .

FROM build-base AS dev

FROM build-base AS build
ARG HUGO_ENV
ARG DOCS_URL
RUN make docker-build HUGO_ENV="${HUGO_ENV}" DOCS_URL="${DOCS_URL}"

FROM scratch AS release
COPY --from=build /out /

FROM scratch AS update-stats
COPY --from=build /src/hugo_stats.json /hugo_stats.json

FROM build-base AS validate-stats
COPY --from=build /src/hugo_stats.json /src/hugo_stats.json
RUN make docker-validate-stats

FROM build-base AS lint
RUN make docker-lint

FROM base AS htmltest-base
ARG HTMLTEST_VERSION
RUN make docker-install-htmltest HTMLTEST_VERSION=${HTMLTEST_VERSION}

FROM build-base AS test
COPY --from=htmltest-base /usr/local/bin/htmltest /usr/local/bin/htmltest
COPY --from=build /out /src/public
RUN make docker-htmltest HTMLTEST_DIR=/src/public

FROM build-base AS update-modules
ARG MODULE
RUN make docker-update-modules MODULE="${MODULE}"

FROM scratch AS vendor
COPY --from=update-modules /src/_vendor /_vendor
COPY --from=update-modules /src/go.* /

FROM build-base AS build-upstream
ARG UPSTREAM_MODULE_NAME
ARG UPSTREAM_REPO
ARG UPSTREAM_COMMIT
ENV HUGO_MODULE_REPLACEMENTS="github.com/${UPSTREAM_MODULE_NAME} -> github.com/${UPSTREAM_REPO} ${UPSTREAM_COMMIT}"
RUN make docker-build-upstream UPSTREAM_MODULE_NAME="${UPSTREAM_MODULE_NAME}"

FROM build-base AS validate-upstream
COPY --from=htmltest-base /usr/local/bin/htmltest /usr/local/bin/htmltest
COPY --from=build-upstream /out /src/public
RUN make docker-htmltest HTMLTEST_DIR=/src/public

FROM dev


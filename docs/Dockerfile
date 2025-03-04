FROM jekyll/jekyll:4.0

RUN apk update && apk add --no-cache \
    build-base \
    libffi-dev \
    readline-dev \
    zlib-dev \
    openssl-dev \
    bash \
    wget \
    curl

ENV RUBY_BASE_VERSION=3.2
ENV RUBY_VERSION=3.2.6

RUN wget https://cache.ruby-lang.org/pub/ruby/${RUBY_BASE_VERSION}/ruby-${RUBY_VERSION}.tar.gz && \
    tar -xzf ruby-${RUBY_VERSION}.tar.gz && \
    cd ruby-${RUBY_VERSION} && \
    ./configure && \
    make && \
    make install && \
    cd .. && \
    rm -rf ruby-${RUBY_VERSION} ruby-${RUBY_VERSION}.tar.gz

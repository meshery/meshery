FROM ghcr.io/gohugoio/hugo:v0.148.1-ext

WORKDIR /src

COPY . /src

EXPOSE 1313

CMD ["server", "--bind", "0.0.0.0", "--baseURL", "http://localhost", "-D"]

/* eslint-disable no-unused-vars */
const { createServer } = require("http");
const { parse } = require("url");
const next = require("next");

const port = parseInt(process.env.PORT, 10) || 3000;
const dev = process.env.NODE_ENV !== "production";
const app = next({ dev });
const handle = app.getRequestHandler();
var httpProxy = require("http-proxy");

var proxy = httpProxy.createProxyServer({ target : { host : "localhost", port : 9081 } });

proxy.on("error", function (err, req, res) {
  res.writeHead(500, { "Content-Type" : "text/plain" });
  res.end('Proxy issue in Meshery-UI');
});

app.prepare().then(() => {
  let server = createServer((req, res) => {
    // Be sure to pass `true` as the second argument to `url.parse`.
    // This tells it to parse the query portion of the URL.
    const { pathname } = parse(req.url, true);
    if (
      pathname.startsWith("/api") ||
      pathname.startsWith("/user/logout") ||
      pathname.startsWith("/user/login") ||
      pathname.startsWith("/provider")
    ) {
      proxy.web(req, res);
    } else {
      handle(req, res);
    }
  });
  server.on("upgrade", (req, socket, head) => {
    proxy.ws(req, socket, head, (err) => {
      socket.write("HTTP/" + req.httpVersion + " 500 Connection error\r\n\r\n");
      socket.end();
    });
  });

  server.listen(port, (err) => {
    if (err) throw err;
    console.log(`> Ready on http://localhost:${port}`);
  });
});

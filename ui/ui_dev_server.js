/* eslint-disable no-unused-vars */
const { createServer } = require('http')
const detect = require('detect-port');
const { parse } = require('url')
const next = require('next')

var port = parseInt(process.env.PORT, 10) || 4000
detect(port, (err, _port) => {
  while (port !== _port) {
    port = port + 1;
  }
})

var dev = process.env.NODE_ENV !== 'production'
var app = next({ dev })
var handle = app.getRequestHandler()
var httpProxy = require('http-proxy');


var proxy = httpProxy.createProxyServer({ target : { host : "localhost",
  port : 9081 } });




proxy.on('error', function (err, req, res) {
  res.writeHead(500, { 'Content-Type' : 'text/plain' });
  res.end('Unexpected issue.');
});



app.prepare().then(() => {
  let server = createServer((req, res) => {
    // Be sure to pass `true` as the second argument to `url.parse`.
    // This tells it to parse the query portion of the URL.
    const { pathname } = parse(req.url, true);
    if (pathname.startsWith("/api") || pathname.startsWith("/user/logout") || pathname.startsWith("/user/login")){
      proxy.web(req, res);
    } else {
      handle(req, res)
    }

  })
  server.on("upgrade", (req,socket,head) => {
    proxy.ws(req,socket, head, err => {
      socket.write("HTTP/" + req.httpVersion + " 500 Connection error\r\n\r\n");
      socket.end();
    })
  })

  server.listen(port, err => {
    if (err) throw err
    console.log(`> Ready on http://localhost:${port}`)
  })
})
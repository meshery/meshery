const { createServer } = require('http')
const { parse } = require('url')
const next = require('next')

const dev = process.env.NODE_ENV !== 'production'
const hostname = 'localhost'
const port = 3001
// when using middleware `hostname` and `port` must be provided below
const app = next({ dev, hostname, port })
const handle = app.getRequestHandler()
const httpProxy = require('http-proxy');

const proxy = httpProxy.createProxyServer({});

proxy.on("error", function (err, req, res) {
  res.writeHead(500, { "Content-Type" : "text/plain" });
  res.end('Proxy issue in Meshery Provider-UI');
});

app.prepare().then(() => {
  createServer(async (req, res) => {
    try {
      // Be sure to pass `true` as the second argument to `url.parse`.
      // This tells it to parse the query portion of the URL.
      const parsedUrl = parse(req.url, true)
      const { pathname } = parsedUrl

      if (pathname.startsWith('/api') ||
                pathname.startsWith('/user/logout') ||
                pathname.startsWith('/user/login')) {
        proxy.web(req, res, { target : `http://${hostname}:9081`, })
      } else {
        if (req.url.startsWith('/provider')) {
          req.url = req.url.replace('/provider', '')
        }
        await handle(req, res, parsedUrl)
      }
    } catch (err) {
      console.error('Error occurred handling', req.url, err)
      res.statusCode = 500
      res.end('internal server error')
    }
  }).listen(port, (err) => {
    if (err) throw err
    console.log(`> Ready on http://${hostname}:${port}`)
  })
})

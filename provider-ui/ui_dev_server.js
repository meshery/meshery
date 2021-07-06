const { createServer } = require('http');
const { parse } = require('url');
const next = require('next');

const port = parseInt(process.env.PORT, 10) || 3001;
const dev = process.env.NODE_ENV !== 'production';
const app = next({ dev });
const handle = app.getRequestHandler();
const httpProxy = require('http-proxy');

const proxy = httpProxy.createProxyServer({});

proxy.on('error', (err, req, res) => {
  res.writeHead(500, {
    'Content-Type': 'text/plain',
  });
  res.end('Unexpected issue.');
});


app.prepare().then(() => {
  createServer((req, res) => {
    // Be sure to pass `true` as the second argument to `url.parse`.
    // This tells it to parse the query portion of the URL.
    const { pathname } = parse(req.url, true);
    if (pathname.startsWith('/api') || pathname.startsWith('/user/logout') || pathname.startsWith('/user/login')) {
      proxy.web(req, res, { target: 'http://localhost:9081' });
    } else {
      if (req.url.startsWith('/provider')) {
        req.url = req.url.replace('/provider', '');
      }
      handle(req, res);
    }
    // app.render(req, res, '/blog', Object.assign(params, query))
  }).listen(port, (err) => {
    if (err) throw err;
    console.log(`> Ready on http://localhost:${port}`);
  });
});

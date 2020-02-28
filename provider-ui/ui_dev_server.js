const { createServer } = require('http');
const { parse } = require('url');
const next = require('next');
// const pathMatch = require('path-match')

const port = parseInt(process.env.PORT, 10) || 3001;
const dev = process.env.NODE_ENV !== 'production';
const app = next({ dev });
const handle = app.getRequestHandler();
// const route = pathMatch()
// const match = route('/blog/:id')
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
    const { pathname } = parse(req.url, true);
    if (pathname.startsWith('/api') || pathname.startsWith('/logout') || pathname.startsWith('/login')) {
      proxy.web(req, res, { target: 'http://localhost:9081' });
    } else {
      if (req.url.startsWith('/provider')) {
        req.url = req.url.replace('/provider', '');
      }
      handle(req, res);
    }
    // const params = match(pathname)
    // if (params === false) {
    //   handle(req, res)
    //   return
    // }
    // // assigning `query` into the params means that we still
    // // get the query string passed to our application
    // // i.e. /blog/foo?show-comments=true
    // app.render(req, res, '/blog', Object.assign(params, query))
  }).listen(port, (err) => {
    if (err) throw err;
    console.log(`> Ready on http://localhost:${port}`);
  });
});

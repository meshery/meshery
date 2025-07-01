const PROXY_ORIGIN = "http://127.0.0.1:7877"; // ✅ Your proxy origin

self.addEventListener("install", (event) => {
  console.log("Service worker installing...");
  self.skipWaiting(); // activate immediately
});

self.addEventListener("activate", (event) => {
  console.log("Service worker activating...");
  self.clients.claim();
});

self.addEventListener("fetch", (event) => {
  const { request } = event;
  const url = new URL(request.url);

  console.log("Intercepted request:", request.url);

  const isFileProtocol = url.protocol === "file:";
  const isSameOrigin = url.origin === self.location.origin;
  const isNavigation = request.mode === "navigate";

  // Only intercept file:// or same-origin, non-navigation requests
  if ((isSameOrigin && !isNavigation) || isFileProtocol) {
    // Create proxy URL
    const proxyUrl = new URL(PROXY_ORIGIN);
    proxyUrl.pathname = url.pathname; // preserve full path

    // Preserve query string
    proxyUrl.search = url.search;

    const newRequest = new Request(proxyUrl.href, {
      method: request.method,
      headers: request.headers,
      body:
        request.method !== "GET" && request.method !== "HEAD"
          ? request.body
          : undefined,
      mode: "cors", // avoid 'navigate' mode
      credentials: request.credentials,
      redirect: request.redirect,
      referrer: request.referrer,
      referrerPolicy: request.referrerPolicy,
    });

    console.log("Proxying request to:", newRequest, newRequest.url);
    event.respondWith(fetch(newRequest));
  }
});

# Old Compatibility Matrix

> An installation compatibility matrix and project test status dashboard.

Source: /pr-preview/pr-21670/project/compatibility-matrix/compatibility-matrix/

Meshery Server and Meshery Adapters are tested daily for their compatibility with the infrastructure they manage and the platforms Meshery deploys on (Kubernetes and Docker). End-to-end test results are automatically posted to the following compatibility matrix.

<style>
  td:hover,tr:hover {
      background-color: var(--color-primary-dark);
      cursor:pointer;
    }
    td.details {
      background-color: #fafafa;
      cursor:text;
    }
    .yellowCheckbox{
      width:2.5rem
    }
    .tooltipss{
      position:relative;
      width:fit-content;
      cursor:pointer;
      margin-right: auto;
      margin-left: auto;
    }
    .tooltipss .tooltiptext {
    visibility: hidden;
    width: 120px;
    background-color: #555;
    color: #fff;
    text-align: center;
    border-radius: 6px;
    padding: 5px 0;
    position: absolute;
    z-index: 1;
    bottom: 125%;
    left: 50%;
    margin-left: -60px;
    opacity: 0;
    transition: opacity 0.3s;
    }

  .tooltipss .tooltiptext::after {
    content: "";
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -5px;
    border-width: 5px;
    border-style: solid;
    border-color: #555 transparent transparent transparent;
  }
  .tooltipss:hover .tooltiptext {
    visibility: visible;
    opacity: 1;
}
</style>

<table class="table table-striped" >
  <thead>
    <tr>
      <th style="text-align: left; padding: 10px; border-bottom: 2px solid #ccc;">Kubernetes Version</th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">
        <a href="https://github.com/meshery/meshery-istio" style="text-decoration: none; color: #1dd1a1; cursor: pointer;">
          <img style="height: 1.5rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/istio.svg" />
          meshery-istio
        </a>
      </th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">
        <a href="https://github.com/meshery/meshery-linkerd" style="text-decoration: none; color: #1dd1a1; cursor: pointer;">
          <img style="height: 1.5rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/linkerd.svg" />
          meshery-linkerd
        </a>
      </th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">
        <a href="https://github.com/meshery/meshery-kuma" style="text-decoration: none; color: #1dd1a1; cursor: pointer;">
          <img style="height: 1.5rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kuma.svg" />
          meshery-kuma
        </a>
      </th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">
        <a href="https://github.com/meshery/meshery-nginx-sm" style="text-decoration: none; color: #1dd1a1; cursor: pointer;">
          <img style="height: 1.5rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/nginx-sm.svg" />
          meshery-nginx-sm
        </a>
      </th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">
        <a href="https://github.com/meshery/meshery-traefik-mesh" style="text-decoration: none; color: #1dd1a1; cursor: pointer;">
          <img style="height: 1.5rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/traefik-mesh.svg" />
          meshery-traefik-mesh
        </a>
      </th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">
        <a href="https://github.com/meshery/meshery-cilium" style="text-decoration: none; color: #1dd1a1; cursor: pointer;">
          <img style="height: 1.5rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/cilium.svg" />
          meshery-cilium
        </a>
      </th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">
        <a href="https://github.com/meshery/meshery-consul" style="text-decoration: none; color: #1dd1a1; cursor: pointer;">
          <img style="height: 1.5rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/consul.svg" />
          meshery-consul
        </a>
      </th>
    </tr>
  </thead>
  <tbody><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.32.0</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">N/A</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.31.4</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">N/A</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.29.0</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">N/A</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.28.4</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">N/A</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.27.8</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">N/A</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.26.11</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">N/A</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.25.3</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">✗</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.25.2</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">✗</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.24.7</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">✗</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.23.9</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">✗</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.23.13</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">✓</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.22.2</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">✗</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.21.5</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">✗</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.21.0</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">N/A</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">N/A</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.20.11</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">✗</td>
      </tr><tr style="border-bottom: 1px solid #eee;">
        <td style="padding: 10px;">v1.20.1</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-istio" onclick="redirectToPastResults('meshery-istio')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-linkerd" onclick="redirectToPastResults('meshery-linkerd')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-kuma" onclick="redirectToPastResults('meshery-kuma')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-nginx-sm" onclick="redirectToPastResults('meshery-nginx-sm')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-traefik-mesh" onclick="redirectToPastResults('meshery-traefik-mesh')">✓</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-cilium" onclick="redirectToPastResults('meshery-cilium')">✗</td><td style="text-align: center; padding: 10px; cursor: pointer;" class="compatibility-status" data-adapter="meshery-consul" onclick="redirectToPastResults('meshery-consul')">N/A</td>
      </tr></tbody>
</table>

<script>
  function redirectToPastResults(adapter) {
    
    const adapterMap = {
      'meshery-istio': 'meshery-istio-past-results',
      'meshery-linkerd': 'meshery-linkerd-past-results',
      'meshery-kuma': 'meshery-kuma-past-results',
      'meshery-nginx-sm': 'meshery-nginx-past-results',
      'meshery-traefik-mesh': 'meshery-traefik-mesh-past-results',
      'meshery-cilium': 'meshery-cilium-past-results',
      'meshery-consul': 'meshery-consul-past-results'
    };
    
    const page = adapterMap[adapter];
    if (page) {
      window.location.href = `\/pr-preview\/pr-21670\/project\/compatibility-matrix\/${page}`;
    }
  }
  
  function showCompatability() {
    let statusCells = document.querySelectorAll(".compatibility-status");
    for(let i = 0; i < statusCells.length; i++) {
      let status = statusCells[i].innerHTML.trim();
      if (status === "✓") {
        statusCells[i].innerHTML = `
          <div class="tooltipss" style="text-align:center; cursor: pointer;">
            <img src="\/pr-preview\/pr-21670\/project\/compatibility-matrix\/images\/service-meshes\/passing.svg" class="yellowCheckbox" alt="Passing" />
            <span class="tooltiptext">Passing</span>
          </div>`;
      } else if (status === "✗") {
        statusCells[i].innerHTML = `
          <div class="tooltipss" style="text-align:center; cursor: pointer;">
            <img src="\/pr-preview\/pr-21670\/project\/compatibility-matrix\/images\/service-meshes\/failing.svg" class="yellowCheckbox" alt="Failing" />
            <span class="tooltiptext">Failing</span>
          </div>`;
      } else if (status === "—") {
        statusCells[i].innerHTML = `
          <div class="tooltipss" style="text-align:center; cursor: pointer;">
            <img src="\/pr-preview\/pr-21670\/project\/compatibility-matrix\/images\/service-meshes\/passing.svg" class="yellowCheckbox" alt="Partial" />
            <span class="tooltiptext">Partial</span>
          </div>`;
      } else if (status === "N/A") {
        statusCells[i].innerHTML = `
          <div class="tooltipss" style="text-align:center; cursor: pointer;">
            <img src="\/pr-preview\/pr-21670\/project\/compatibility-matrix\/images\/service-meshes\/na-icon.svg" class="yellowCheckbox" alt="Not Applicable" />
            <span class="tooltiptext">Not Applicable</span>
          </div>`;
      }
    }
  }
  
  
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', showCompatability);
  } else {
    showCompatability();
  }
</script>


## Integration Tests

As a key aspect of Meshery, its integrations with other systems are routinely tested. Unit and integration tests before and after every pull request (before code is to be merged into the project and after code is merged into the project). End-to-end tests are run nightly and automatically posted to the following test matrix.

<style>
  .edge_visible{
    display: table-row !important;
    visibility: visible !important;
  }
  .stable_visible{
    display: table-row !important;
    visibility: visible !important;
  }
  .checkbox{
    display: flex;
    justify-content: flex-end;
    align-items: center;
    text-align: left;
  }
  td.details {
    background-color: #fafafa;
    cursor:text;
  }
  .edge_test_text{
    margin-right: 20px;
  }
  .status-passing {
    background-color: #56B257 !important;
    color: white !important;
  }
  .status-partial {
    background-color: #EBC017 !important;
    color: white !important;
  }
  .status-failing {
    background-color: #B32700 !important;
    color: white !important;
  }
  .status-default {
    background-color: #6c757d !important;
    color: white !important;
  }
  #test-table tbody tr {
    background-color: #3a3f47 !important;
    color: #ccc;
    cursor: pointer;
  }
  #test-table tbody tr:hover {
    background-color: #3a3f47 !important;
    color: #ccc;
  }
</style>

<div class="checkbox">
  <div>
    <input onchange="handleEdgeCheckboxChange();" type="checkbox" id="checkbox_edge" value="Edge Tests" checked />
    <label for="checkbox_edge" class="edge_test_text">Edge Channel</label>
  </div>
  <div>
    <input onchange="handleStableCheckboxChange();" type="checkbox" id="checkbox_stable" value="Stable Tests" checked />
    <label for="checkbox_stable">Stable Channel</label>
  </div>
</div>

<table id="test-table" style="text-align: center;">
  <thead>
    <tr>
      <th style="text-align: center;">Status</th>
      <th style="text-align: center;">Meshery Component</th>
      <th style="text-align: center;">Meshery Component Version</th>
      <th style="text-align: center;">Meshery Server Version</th>
      <th style="text-align: center;">Infrastructure</th>
      <th style="text-align: center;">Infrastructure Version</th>
    </tr>
  </thead>
  <tbody><tr style="visibility: hidden; display: none;" class="test-details edge edge_visible" onclick="toggle_visibility('meshery-consul-edge');">
            <td class="status-failing">2023-11-02 23:47:58 UTC Thu</td>
            <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-consul">meshery-consul</a></td>
            <td><a href="https://github.com/meshery/meshery-consul/releases">edge</a></td>
            <td><a href="https://github.com/meshery/meshery/releases/tag/"></a></td><td style="white-space: nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/consul.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/consul/">Consul</a></td>
            <td>v1.2.3</td>
          </tr>
          <tr class="hidden-details" id="meshery-consul-edge" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <ul style="margin: 5px 0;">
                <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.21.5</li>
              </ul>
            </td>
            <td colspan="4" class="details">
              <i>Test results:</i>
              <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>consul-client</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>consul-server</td>
                    </tr></table>
            </td>
          </tr><tr style="visibility: hidden; display: none;" class="test-details edge edge_visible" onclick="toggle_visibility('meshery-istio-edge');">
            <td class="status-failing">2025-01-16 01:06:39 UTC Thu</td>
            <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-istio">meshery-istio</a></td>
            <td><a href="https://github.com/meshery/meshery-istio/releases">edge</a></td>
            <td><a href="https://github.com/meshery/meshery/releases/tag/v0.8.11">v0.8.11</a></td><td style="white-space: nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/istio.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/istio/">istio</a></td>
            <td>1.23.4</td>
          </tr>
          <tr class="hidden-details" id="meshery-istio-edge" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <ul style="margin: 5px 0;">
                <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.31.4</li>
              </ul>
            </td>
            <td colspan="4" class="details">
              <i>Test results:</i>
              <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/grafana-addon</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/istio-egressgateway</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/istio-ingressgateway</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/istiod</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/prometheus-addon</td>
                    </tr></table>
            </td>
          </tr><tr style="visibility: hidden; display: none;" class="test-details edge edge_visible" onclick="toggle_visibility('meshery-kuma-edge');">
            <td class="status-failing">2025-10-31 00:26:32 UTC Fri</td>
            <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-kuma">meshery-kuma</a></td>
            <td><a href="https://github.com/meshery/meshery-kuma/releases">edge</a></td>
            <td><a href="https://github.com/meshery/meshery/releases/tag/v0.8.147">v0.8.147</a></td><td style="white-space: nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kuma.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/kuma/">Kuma</a></td>
            <td>2.12.3</td>
          </tr>
          <tr class="hidden-details" id="meshery-kuma-edge" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <ul style="margin: 5px 0;">
                <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.32.0</li>
              </ul>
            </td>
            <td colspan="4" class="details">
              <i>Test results:</i>
              <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/kuma-control-plane</td>
                    </tr></table>
            </td>
          </tr><tr style="visibility: hidden; display: none;" class="test-details edge edge_visible" onclick="toggle_visibility('meshery-linkerd-edge');">
            <td class="status-failing">2025-11-14 00:24:21 UTC Fri</td>
            <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-linkerd">meshery-linkerd</a></td>
            <td><a href="https://github.com/meshery/meshery-linkerd/releases">edge</a></td>
            <td><a href="https://github.com/meshery/meshery/releases/tag/v0.8.160">v0.8.160</a></td><td style="white-space: nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/linkerd.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/linkerd/">Linkerd</a></td>
            <td></td>
          </tr>
          <tr class="hidden-details" id="meshery-linkerd-edge" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <ul style="margin: 5px 0;">
                <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.31.4</li>
              </ul>
            </td>
            <td colspan="4" class="details">
              <i>Test results:</i>
              <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/linkerd-destination</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/linkerd-identity</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/linkerd-proxy-injector</td>
                    </tr></table>
            </td>
          </tr><tr style="visibility: hidden; display: none;" class="test-details edge edge_visible" onclick="toggle_visibility('meshery-nginx-sm-edge');">
            <td class="status-failing">2024-01-25 07:23:35 UTC Thu</td>
            <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-nginx-sm">meshery-nginx-sm</a></td>
            <td><a href="https://github.com/meshery/meshery-nginx-sm/releases">edge</a></td>
            <td><a href="https://github.com/meshery/meshery/releases/tag/v0.7.12">v0.7.12</a></td><td style="white-space: nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/nginx-sm.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/nginx-sm/">nginx-sm</a></td>
            <td>v2.0.0</td>
          </tr>
          <tr class="hidden-details" id="meshery-nginx-sm-edge" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <ul style="margin: 5px 0;">
                <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.22.2</li>
              </ul>
            </td>
            <td colspan="4" class="details">
              <i>Test results:</i>
              <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/nginx-mesh-api</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/nginx-mesh-metrics</td>
                    </tr></table>
            </td>
          </tr><tr style="visibility: hidden; display: none;" class="test-details edge edge_visible" onclick="toggle_visibility('meshery-osm-edge');">
            <td class="status-failing">2023-07-19 23:01:53 UTC Wed</td>
            <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-osm">meshery-osm</a></td>
            <td><a href="https://github.com/meshery/meshery-osm/releases">edge</a></td>
            <td><a href="https://github.com/meshery/meshery/releases/tag/"></a></td><td style="white-space: nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/osm.svg" />OSM</td>
            <td>v1.2.4</td>
          </tr>
          <tr class="hidden-details" id="meshery-osm-edge" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <ul style="margin: 5px 0;">
                <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.22.2</li>
              </ul>
            </td>
            <td colspan="4" class="details">
              <i>Test results:</i>
              <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/osm-bootstrap</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/osm-controller</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                      <td>pod/osm-injector</td>
                    </tr></table>
            </td>
          </tr><tr style="visibility: hidden; display: none;" class="test-details edge edge_visible" onclick="toggle_visibility('meshery-traefik-mesh-edge');">
            <td class="status-passing">2023-10-26 23:44:53 UTC Thu</td>
            <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-traefik-mesh">meshery-traefik-mesh</a></td>
            <td><a href="https://github.com/meshery/meshery-traefik-mesh/releases">edge</a></td>
            <td><a href="https://github.com/meshery/meshery/releases/tag/"></a></td><td style="white-space: nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/traefik-mesh.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/">traefik-mesh</a></td>
            <td></td>
          </tr>
          <tr class="hidden-details" id="meshery-traefik-mesh-edge" style="visibility:hidden; display:none;">
            <td colspan="2" class="details">
              <i>Platform:</i>
              <ul style="margin: 5px 0;">
                <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.21.5</li>
              </ul>
            </td>
            <td colspan="4" class="details">
              <i>Test results:</i>
              <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                      <td>pod/grafana-core</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                      <td>pod/jaeger</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                      <td>pod/prometheus-core</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                      <td>pod/traefik-mesh-controller</td>
                    </tr><tr style="background-color: transparent;">
                      <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                      <td>pod/traefik-mesh-proxy</td>
                    </tr></table>
            </td>
          </tr><tr style="visibility: hidden; display: none;" class="test-details stable stable_visible" onclick="toggle_visibility('meshery-cilium-stable');">
          <td class="status-failing">2023-11-10 06:53:32 UTC Fri</td>
          <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-cilium">meshery-cilium</a></td>
          <td><a href="https://github.com/meshery/meshery-cilium/releases/tag/v0.6.8">v0.6.8</a></td>
          <td><a href="https://github.com/meshery/meshery/releases/tag/"></a></td><td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/cilium.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/cilium/">Cilium</a></td>
          <td>v1.14.3</td>
        </tr>
        <tr class="hidden-details" id="meshery-cilium-stable" style="visibility:hidden; display:none;">
          <td colspan="2" class="details">
            <i>Platform:</i>
            <ul style="margin: 5px 0;">
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.22.2</li>
            </ul>
          </td>
          <td colspan="4" class="details">
            <i>Test results:</i>
            <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/cilium</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/cilium-operator</td>
                  </tr></table>
          </td>
        </tr><tr style="visibility: hidden; display: none;" class="test-details stable stable_visible" onclick="toggle_visibility('meshery-consul-stable');">
          <td class="status-failing">2023-10-31 04:00:04 UTC Tue</td>
          <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-consul">meshery-consul</a></td>
          <td><a href="https://github.com/meshery/meshery-consul/releases/tag/v1.0.0">v1.0.0</a></td>
          <td><a href="https://github.com/meshery/meshery/releases/tag/"></a></td><td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/consul.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/consul/">Consul</a></td>
          <td>v1.2.2</td>
        </tr>
        <tr class="hidden-details" id="meshery-consul-stable" style="visibility:hidden; display:none;">
          <td colspan="2" class="details">
            <i>Platform:</i>
            <ul style="margin: 5px 0;">
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.21.5</li>
            </ul>
          </td>
          <td colspan="4" class="details">
            <i>Test results:</i>
            <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>consul-client</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>consul-server</td>
                  </tr></table>
          </td>
        </tr><tr style="visibility: hidden; display: none;" class="test-details stable stable_visible" onclick="toggle_visibility('meshery-istio-stable');">
          <td class="status-failing">2024-02-26 21:27:39 UTC Mon</td>
          <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-istio">meshery-istio</a></td>
          <td><a href="https://github.com/meshery/meshery-istio/releases/tag/v0.7.0">v0.7.0</a></td>
          <td><a href="https://github.com/meshery/meshery/releases/tag/v0.7.23">v0.7.23</a></td><td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/istio.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/istio/">istio</a></td>
          <td>1.20.3</td>
        </tr>
        <tr class="hidden-details" id="meshery-istio-stable" style="visibility:hidden; display:none;">
          <td colspan="2" class="details">
            <i>Platform:</i>
            <ul style="margin: 5px 0;">
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.29.0</li>
            </ul>
          </td>
          <td colspan="4" class="details">
            <i>Test results:</i>
            <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/istio-egressgateway</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/istio-ingressgateway</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/istiod</td>
                  </tr></table>
          </td>
        </tr><tr style="visibility: hidden; display: none;" class="test-details stable stable_visible" onclick="toggle_visibility('meshery-kuma-stable');">
          <td class="status-failing">2024-12-27 20:42:19 UTC Fri</td>
          <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-kuma">meshery-kuma</a></td>
          <td><a href="https://github.com/meshery/meshery-kuma/releases/tag/v0.8.0">v0.8.0</a></td>
          <td><a href="https://github.com/meshery/meshery/releases/tag/v0.8.1">v0.8.1</a></td><td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kuma.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/kuma/">Kuma</a></td>
          <td>2.9.2</td>
        </tr>
        <tr class="hidden-details" id="meshery-kuma-stable" style="visibility:hidden; display:none;">
          <td colspan="2" class="details">
            <i>Platform:</i>
            <ul style="margin: 5px 0;">
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.22.2</li>
            </ul>
          </td>
          <td colspan="4" class="details">
            <i>Test results:</i>
            <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/kuma-control-plane</td>
                  </tr></table>
          </td>
        </tr><tr style="visibility: hidden; display: none;" class="test-details stable stable_visible" onclick="toggle_visibility('meshery-linkerd-stable');">
          <td class="status-failing">2023-02-07 19:39:19 UTC Tue</td>
          <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-linkerd">meshery-linkerd</a></td>
          <td><a href="https://github.com/meshery/meshery-linkerd/releases/tag/v0.6.7">v0.6.7</a></td>
          <td><a href="https://github.com/meshery/meshery/releases/tag/v0.6.50">v0.6.50</a></td><td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/linkerd.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/linkerd/">Linkerd</a></td>
          <td></td>
        </tr>
        <tr class="hidden-details" id="meshery-linkerd-stable" style="visibility:hidden; display:none;">
          <td colspan="2" class="details">
            <i>Platform:</i>
            <ul style="margin: 5px 0;">
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.24.7</li>
            </ul>
          </td>
          <td colspan="4" class="details">
            <i>Test results:</i>
            <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/linkerd-destination</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/linkerd-identity</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/linkerd-proxy-injector</td>
                  </tr></table>
          </td>
        </tr><tr style="visibility: hidden; display: none;" class="test-details stable stable_visible" onclick="toggle_visibility('meshery-nginx-sm-stable');">
          <td class="status-failing">2023-01-16 13:40:30 UTC Mon</td>
          <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-nginx-sm">meshery-nginx-sm</a></td>
          <td><a href="https://github.com/meshery/meshery-nginx-sm/releases/tag/v0.6.3">v0.6.3</a></td>
          <td><a href="https://github.com/meshery/meshery/releases/tag/v0.6.43">v0.6.43</a></td><td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/nginx-sm.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/nginx-sm/">nginx-sm</a></td>
          <td>v1.6.0</td>
        </tr>
        <tr class="hidden-details" id="meshery-nginx-sm-stable" style="visibility:hidden; display:none;">
          <td colspan="2" class="details">
            <i>Platform:</i>
            <ul style="margin: 5px 0;">
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.25.3</li>
            </ul>
          </td>
          <td colspan="4" class="details">
            <i>Test results:</i>
            <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/nginx-mesh-api</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/failing.svg" alt="Test Status" /></td>
                    <td>pod/nginx-mesh-metrics</td>
                  </tr></table>
          </td>
        </tr><tr style="visibility: hidden; display: none;" class="test-details stable stable_visible" onclick="toggle_visibility('meshery-osm-stable');">
          <td class="status-passing">2022-12-30 21:13:32 UTC Fri</td>
          <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-osm">meshery-osm</a></td>
          <td><a href="https://github.com/meshery/meshery-osm/releases/tag/v0.6.4">v0.6.4</a></td>
          <td><a href="https://github.com/meshery/meshery/releases/tag/v0.6.39">v0.6.39</a></td><td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/osm.svg" />OSM</td>
          <td>v1.2.3</td>
        </tr>
        <tr class="hidden-details" id="meshery-osm-stable" style="visibility:hidden; display:none;">
          <td colspan="2" class="details">
            <i>Platform:</i>
            <ul style="margin: 5px 0;">
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.24.7</li>
            </ul>
          </td>
          <td colspan="4" class="details">
            <i>Test results:</i>
            <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                    <td>pod/osm-bootstrap</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                    <td>pod/osm-controller</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                    <td>pod/osm-injector</td>
                  </tr></table>
          </td>
        </tr><tr style="visibility: hidden; display: none;" class="test-details stable stable_visible" onclick="toggle_visibility('meshery-traefik-mesh-stable');">
          <td class="status-passing">2023-07-22 05:51:44 UTC Sat</td>
          <td style="white-space:nowrap;"><a href="https://github.com/meshery/meshery-traefik-mesh">meshery-traefik-mesh</a></td>
          <td><a href="https://github.com/meshery/meshery-traefik-mesh/releases/tag/v0.6.8">v0.6.8</a></td>
          <td><a href="https://github.com/meshery/meshery/releases/tag/"></a></td><td style="white-space:nowrap;"><img style="height: 2rem; vertical-align: text-bottom; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/traefik-mesh.svg" /><a href="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/">traefik-mesh</a></td>
          <td></td>
        </tr>
        <tr class="hidden-details" id="meshery-traefik-mesh-stable" style="visibility:hidden; display:none;">
          <td colspan="2" class="details">
            <i>Platform:</i>
            <ul style="margin: 5px 0;">
              <li><img style="height: 1rem; vertical-align: text-bottom;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/kubernetes-icon-color.svg" /> minikube v1.22.2</li>
            </ul>
          </td>
          <td colspan="4" class="details">
            <i>Test results:</i>
            <table style="border:0; margin-top: 5px;"><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                    <td>pod/grafana-core</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                    <td>pod/jaeger</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                    <td>pod/prometheus-core</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                    <td>pod/traefik-mesh-controller</td>
                  </tr><tr style="background-color: transparent;">
                    <td><img style="height: 24px; width: 24px; margin-right: 8px;" src="/pr-preview/pr-21670/project/compatibility-matrix/images/service-meshes/passing.svg" alt="Test Status" /></td>
                    <td>pod/traefik-mesh-proxy</td>
                  </tr></table>
          </td>
        </tr></tbody>
</table>

<script>
  function toggle_visibility(id) {
    var e = document.getElementById(id);
    if (e.style.visibility == "visible") {
      e.style.display = "none";
      e.style.visibility = "hidden";
    } else {
      e.style.display = "table-row";
      e.style.visibility = "visible";
    }
  }

  function handleEdgeCheckboxChange() {
    let e = document.getElementsByClassName("edge");
    let stable = document.getElementsByClassName("stable");
    let stable_box = document.getElementById("checkbox_stable");
    for (let i = 0; i < e.length; i++) {
      if (e[i].classList.contains("edge_visible")) {
        e[i].classList.remove("edge_visible");
        if (!stable_box.checked) {
          stable_box.checked = true;
          handleStableCheckboxChange();
        }
      } else {
        e[i].classList.add("edge_visible");
      }
    }
  }

  function handleStableCheckboxChange() {
    let e = document.getElementsByClassName("stable");
    let edge_box = document.getElementById("checkbox_edge");
    for (let i = 0; i < e.length; i++) {
      if (e[i].classList.contains("stable_visible")) {
        e[i].classList.remove("stable_visible");
        if (!edge_box.checked) {
          edge_box.checked = true;
          handleEdgeCheckboxChange();
        }
      } else {
        e[i].classList.add("stable_visible");
      }
    }
  }
</script>

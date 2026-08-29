# Adapters

> Adapters extend Meshery's management capabilities in any number of ways, including lifecycle, configuration, performance, governance, identity...

Source: /pr-preview/pr-21670/concepts/architecture/adapters/

## What are Meshery Adapters?

Part of Meshery's extensibility as a platform, Meshery Adapters are purpose-built to address an area in need of management that is either considered optional to the platform and/or is considered an area in which additional depth of control is needed. Adapters extend Meshery's management capabilities in any number of ways, including lifecycle, configuration, performance, governance, identity and so on. Meshery Adapters come in different form factors, and depending on their purpose, deliver different sets or capabilities. Each Adapter registers its capabilities with Meshery Server. Meshery Server, in-turn, exposes those capabilities for you to control.

## Meshery Adapters for Lifecycle Management

Adapters that extend Meshery's lifecycle management capabilities for infrastructure do so, by offering an infrastructure-specific interface to increase the depth of control that Meshery has over a particular technology. Meshery uses adapters to offer choice of load generator (for performance management) and for managing different layers of your infrastructure. Lifecycle adapters allow Meshery to interface with the different cloud native infrastructure, exposing their differentiated value to users.

Meshery has lifecycle adapters for managing the following cloud native infrastructure.

<table class="table table-striped">
  <thead>
    <tr>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">Adapter Status</th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">Adapter</th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">Port</th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">Earliest Version Supported</th>
    </tr>
  </thead>
  <tbody><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#4caf50;
              color: white; border-radius: 4px; font-size: 0.85em;">
              stable
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/images/traefik-mesh.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/images/traefik-mesh.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/images/traefik-mesh.svg" 
                   loading="lazy" alt="Meshery Adapter for Traefik Mesh" />
              <span>Meshery Adapter for Traefik Mesh</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10006/gRPC</td>
          <td style="text-align: center; padding: 10px;">v1.0</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#4caf50;
              color: white; border-radius: 4px; font-size: 0.85em;">
              stable
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/nginx-sm/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/nginx-sm/images/nginx-sm.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/nginx-sm/images/nginx-sm-white.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/nginx-sm/images/nginx-sm.svg" 
                   loading="lazy" alt="Meshery Adapter for NGINX Service Mesh" />
              <span>Meshery Adapter for NGINX Service Mesh</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10010/gRPC</td>
          <td style="text-align: center; padding: 10px;">v1.2.0</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#4caf50;
              color: white; border-radius: 4px; font-size: 0.85em;">
              stable
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/nsm/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/nsm/images/nsm.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/nsm/images/nsm.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/nsm/images/nsm.svg" 
                   loading="lazy" alt="Meshery Adapter for Network Service Mesh" />
              <span>Meshery Adapter for Network Service Mesh</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10004/gRPC</td>
          <td style="text-align: center; padding: 10px;">v0.2.1</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#4caf50;
              color: white; border-radius: 4px; font-size: 0.85em;">
              stable
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/linkerd/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/linkerd/images/linkerd.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/linkerd/images/linkerd-white.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/linkerd/images/linkerd.svg" 
                   loading="lazy" alt="Meshery Adapter for Linkerd" />
              <span>Meshery Adapter for Linkerd</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10001/gRPC</td>
          <td style="text-align: center; padding: 10px;">v2.10.2</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#4caf50;
              color: white; border-radius: 4px; font-size: 0.85em;">
              stable
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/kuma/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/kuma/images/kuma.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/kuma/images/kuma-white.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/kuma/images/kuma.svg" 
                   loading="lazy" alt="Meshery Adapter for Kuma" />
              <span>Meshery Adapter for Kuma</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10007/gRPC</td>
          <td style="text-align: center; padding: 10px;">v1.2.2</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#4caf50;
              color: white; border-radius: 4px; font-size: 0.85em;">
              stable
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/istio/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/istio/images/istio.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/istio/images/istio-white.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/istio/images/istio.svg" 
                   loading="lazy" alt="Meshery Adapter for Istio" />
              <span>Meshery Adapter for Istio</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10000/gRPC</td>
          <td style="text-align: center; padding: 10px;">v1.6.0</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#4caf50;
              color: white; border-radius: 4px; font-size: 0.85em;">
              stable
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/consul/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/consul/images/consul.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/consul/images/consul-white.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/consul/images/consul.svg" 
                   loading="lazy" alt="Meshery Adapter for Consul" />
              <span>Meshery Adapter for Consul</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10002/gRPC</td>
          <td style="text-align: center; padding: 10px;">v1.8.4</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#4caf50;
              color: white; border-radius: 4px; font-size: 0.85em;">
              stable
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/cilium/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/cilium/images/cilium.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/cilium/images/cilium-white.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/cilium/images/cilium.svg" 
                   loading="lazy" alt="Meshery Adapter for Cilium Service Mesh" />
              <span>Meshery Adapter for Cilium Service Mesh</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10012/gRPC</td>
          <td style="text-align: center; padding: 10px;">v1.10.6</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#ff9800;
              color: white; border-radius: 4px; font-size: 0.85em;">
              beta
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/app-mesh/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/app-mesh/images/app-mesh.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/app-mesh/images/app-mesh-white.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/app-mesh/images/app-mesh.svg" 
                   loading="lazy" alt="Meshery Adapter for App Mesh" />
              <span>Meshery Adapter for App Mesh</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10005/gRPC</td>
          <td style="text-align: center; padding: 10px;">v1.4.1</td>
        </tr><tr style="border-bottom: 1px solid #eee;">
          <td style="text-align: center; padding: 10px;">
            <span style="display: inline-block; padding: 4px 8px; background-color:#9e9e9e;
              color: white; border-radius: 4px; font-size: 0.85em;">
              alpha
            </span>
          </td>
          <td style="text-align: center; padding: 10px;">
            <a href="/pr-preview/pr-21670/extensions/adapters/tanzu-sm/" style="text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 8px;">
              <img src="/pr-preview/pr-21670/extensions/adapters/tanzu-sm/images/tanzu.svg" style="height: 24px; width: auto; display: inline-block;" 
                   data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/tanzu-sm/images/tanzu.svg" 
                   data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/tanzu-sm/images/tanzu.svg" 
                   loading="lazy" alt="Meshery Adapter for Tanzu Service Mesh" />
              <span>Meshery Adapter for Tanzu Service Mesh</span>
            </a>
          </td>
          <td style="text-align: center; padding: 10px;">10011/gRPC</td>
          <td style="text-align: center; padding: 10px;">pre-GA</td>
        </tr></tbody>
</table>


## Meshery Adapters for Performance Management

Meshery Server allows users to generate traffic load tests using fortio.

## Adapter Deployment and Registration

Like every Meshery component, Meshery Adapters use MeshKit.

When an Adapter registers its models with Meshery Server, it also tells the server where to reach it. That is what makes an Adapter more than a source of definitions: at deployment time, any component belonging to a model an Adapter registered is provisioned *by that Adapter*, over gRPC, rather than by Meshery Server. See [Deployment Engine](/pr-preview/pr-21670/concepts/architecture/deployment-engine/) for how Meshery chooses between the two.

### Adapter FAQs

#### Is each Meshery adapter made equal?

No, different Meshery adapters are written to expose the unique value of each cloud native infrastructure. Consequently, they are not equally capable just as each cloud native infrastructure is not equally capable as the other. Each Adapter has a set of operations which are grouped based on predefined operation types. See the [extensibility](/pr-preview/pr-21670/reference/extensibility/) page for more details on adapter operations.

#### How can I create a new adapter?

Yes, see the [extensibility](/pr-preview/pr-21670/reference/extensibility/) documentation for details how to create a new Meshery Adapter. See the Meshery Adapter Template repository as boilerplate for your new adapter.

#### Do adapters have to be written in Golang?

No. Adapters much interface with Meshery Server via gRPC. What language is used in that adapter is the prerogative of a given adapter's maintainers.

#### Can I run more than one instance of the same Meshery adapter?

Yes. The default configuration of a Meshery deployment includes one instance of each of the Meshery adapters (that have reached a stable version status). You may choose to run multiple instances of the same type of Meshery adapter; e.g. two instances of the `meshery-istio` adapter. To do so, modify ~/.meshery/meshery.yaml to include multiple copies of the given adapter.

See the "[Multiple Adapters](/pr-preview/pr-21670/installation/advanced/multiple-adapters/)" guide for more information.

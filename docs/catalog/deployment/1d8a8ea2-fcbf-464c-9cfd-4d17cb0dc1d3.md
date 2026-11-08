---
layout: item
name: On Prem Kubernetes with Secure Cloud Logging
publishedVersion: 0.0.7
userId: 43643949-e1f8-4d55-b65a-c61a2b99a0a1
userName: Chandru M
userAvatarURL: https://avatars.githubusercontent.com/u/121534233?v=4
type: deployment
compatibility:
    - kubernetes
patternId: 1d8a8ea2-fcbf-464c-9cfd-4d17cb0dc1d3
image: /assets/images/logos/service-mesh-pattern.svg
patternInfo: |
  This%20design%20represents%20a%20secure%20on-premises%20Kubernetes%20environment%20for%20running%20and%20automatically%20scaling%20a%20client%20application.%20Fluent%20Bit%20runs%20as%20a%20DaemonSet%20and%20collects%20application%20and%20cluster%20logs.%20Network%20policies%20use%20a%20default-deny%20approach%20and%20permit%20only%20the%20required%20internal%20application%20traffic%20and%20outbound%20log%20delivery%20to%20an%20authorized%20cloud%20logging%20endpoint.%20Cloud%20credentials%20are%20referenced%20through%20a%20Kubernetes%20Secret%2C%20while%20Fluent%20Bit%20configuration%20is%20stored%20in%20a%20ConfigMap
patternCaveats: |
  This%20is%20a%20reference%20architecture%20and%20requires%20environment-specific%20configuration%20before%20production%20use.%20Replace%20placeholder%20cloud%20endpoints%20and%20credentials%2C%20enable%20TLS%20certificate%20verification%2C%20apply%20least-privilege%20RBAC%2C%20encrypt%20Kubernetes%20Secrets%20at%20rest%2C%20and%20rotate%20credentials%20regularly.%20NetworkPolicy%20enforcement%20requires%20a%20compatible%20CNI%20plugin.%20Kubernetes%20NetworkPolicy%20alone%20may%20not%20restrict%20traffic%20by%20domain%20name%2C%20so%20an%20egress%20gateway%2C%20firewall%2C%20or%20proxy%20should%20be%20used%20to%20guarantee%20that%20only%20approved%20logging%20destinations%20are%20reachable.%20Fluent%20Bit%20buffering%2C%20retry%20behavior%2C%20log%20filtering%2C%20and%20sensitive-data%20redaction%20should%20also%20be%20configured
createdAt: 2026-09-07T07:25:44Z
permalink: catalog/deployment/on-prem-kubernetes-with-secure-cloud-logging-1d8a8ea2-fcbf-464c-9cfd-4d17cb0dc1d3.html
URL: 'https://raw.githubusercontent.com/meshery/meshery.io/master/catalog/1d8a8ea2-fcbf-464c-9cfd-4d17cb0dc1d3/0.0.7/design.yml'
downloadLink: 1d8a8ea2-fcbf-464c-9cfd-4d17cb0dc1d3/design.yml
---
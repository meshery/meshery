---
layout: item
name: Network policy
publishedVersion: 0.0.1
userId: 9222bde1-64c6-4fb2-971a-3402d5ae2fd9
userName: Deepak Reddy
userAvatarURL: https://lh3.googleusercontent.com/a/ACg8ocIGbiDtE0q65qVvAUdzHw8Qky81rM0kSAknIqbgysfDCw=s96-c
type: traffic-management
compatibility: 
    - kubernetes
patternId: 7b2e40b0-3cc8-4da3-bccd-b66bc6cd206a
image: https://raw.githubusercontent.com/layer5labs/meshery-extensions-packages/master/action-assets/design-assets/7b2e40b0-3cc8-4da3-bccd-b66bc6cd206a-light.png,https://raw.githubusercontent.com/layer5labs/meshery-extensions-packages/master/action-assets/design-assets/7b2e40b0-3cc8-4da3-bccd-b66bc6cd206a-dark.png
patternInfo: |
  If%20you%20want%20to%20control%20traffic%20flow%20at%20the%20IP%20address%20or%20port%20level%20for%20TCP%2C%20UDP%2C%20and%20SCTP%20protocols%2C%20then%20you%20might%20consider%20using%20Kubernetes%20NetworkPolicies%20for%20particular%20applications%20in%20your%20cluster.%20NetworkPolicies%20are%20an%20application-centric%20construct%20which%20allow%20you%20to%20specify%20how%20a%20pod%20is%20allowed%20to%20communicate%20with%20various%20network%20%22entities%22%20(we%20use%20the%20word%20%22entity%22%20here%20to%20avoid%20overloading%20the%20more%20common%20terms%20such%20as%20%22endpoints%22%20and%20%22services%22%2C%20which%20have%20specific%20Kubernetes%20connotations)%20over%20the%20network.%20NetworkPolicies%20apply%20to%20a%20connection%20with%20a%20pod%20on%20one%20or%20both%20ends%2C%20and%20are%20not%20relevant%20to%20other%20connections.%20
patternCaveats: |
  This%20is%20an%20sample%20network%20%20policy%20with%20ingress%2Cegress%20defined%20%2C%20change%20according%20to%20your%20requirements
permalink: catalog/traffic-management/network-policy-7b2e40b0-3cc8-4da3-bccd-b66bc6cd206a.html
URL: 'https://raw.githubusercontent.com/meshery/meshery.io/master/catalog/7b2e40b0-3cc8-4da3-bccd-b66bc6cd206a/0.0.1/design.yml'
downloadLink: 7b2e40b0-3cc8-4da3-bccd-b66bc6cd206a/design.yml
---
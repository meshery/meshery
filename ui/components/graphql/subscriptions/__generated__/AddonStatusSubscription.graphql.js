/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type AddonStatusSubscriptionVariables = {|
  selector?: ?MeshType
|};
export type AddonStatusSubscriptionResponse = {|
  +addonsState: $ReadOnlyArray<{|
    +type: string,
    +config: {|
      +serviceName: string,
      +endpoint: string,
    |},
  |}>
|};
export type AddonStatusSubscription = {|
  variables: AddonStatusSubscriptionVariables,
  response: AddonStatusSubscriptionResponse,
|};
*/


/*
subscription AddonStatusSubscription(
  $selector: MeshType
) {
  addonsState: listenToAddonState(selector: $selector) {
    type
    config {
      serviceName
      endpoint
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "selector"
  }
],
v1 = [
  {
    "alias": "addonsState",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "AddonList",
    "kind": "LinkedField",
    "name": "listenToAddonState",
    "plural": true,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "type",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "AddonConfig",
        "kind": "LinkedField",
        "name": "config",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "serviceName",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "endpoint",
            "storageKey": null
          }
        ],
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "AddonStatusSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AddonStatusSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "fcb8f0fbff656cbf142d60097b3b6825",
    "id": null,
    "metadata": {},
    "name": "AddonStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription AddonStatusSubscription(\n  $selector: MeshType\n) {\n  addonsState: listenToAddonState(selector: $selector) {\n    type\n    config {\n      serviceName\n      endpoint\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '3b3f1991ed11961f1519a64d8db3392a';

module.exports = node;

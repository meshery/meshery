/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL_MESH" | "APP_MESH" | "CITRIX_SERVICE_MESH" | "CONSUL" | "INVALID_MESH" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORK_SERVICE_MESH" | "NGINX_SERVICE_MESH" | "OCTARINE" | "OPEN_SERVICE_MESH" | "TANZU" | "TRAEFIK_MESH" | "%future added value";
export type AddonStatusSubscriptionVariables = {|
  selector?: ?MeshType
|};
export type AddonStatusSubscriptionResponse = {|
  +addonsState: $ReadOnlyArray<{|
    +name: string,
    +owner: string,
    +endpoint: string,
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
    name
    owner
    endpoint
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
        "name": "name",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "owner",
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
    "cacheID": "653204021c7317b7444668f2cc8771d0",
    "id": null,
    "metadata": {},
    "name": "AddonStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription AddonStatusSubscription(\n  $selector: MeshType\n) {\n  addonsState: listenToAddonState(selector: $selector) {\n    name\n    owner\n    endpoint\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'b2bd87e95067d7db68c966ac0ec0717f';

module.exports = node;

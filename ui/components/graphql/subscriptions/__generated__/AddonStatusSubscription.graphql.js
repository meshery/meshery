/**
 * @generated SignedSource<<5745a37785e352c724b509449718f8ac>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type MeshType = "ALL_MESH" | "INVALID_MESH" | "APP_MESH" | "CITRIX_SERVICE_MESH" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "TRAEFIK_MESH" | "OCTARINE" | "NETWORK_SERVICE_MESH" | "TANZU" | "OPEN_SERVICE_MESH" | "NGINX_SERVICE_MESH" | "CILIUM_SERVICE_MESH" | "%future added value";
export type ServiceMeshFilter = {|
  type?: ?MeshType,
  k8sClusterIDs?: ?$ReadOnlyArray<string>,
|};
export type AddonStatusSubscription$variables = {|
  filter?: ?ServiceMeshFilter,
|};
export type AddonStatusSubscription$data = {|
  +addonsState: $ReadOnlyArray<{|
    +name: string,
    +owner: string,
  |}>,
|};
export type AddonStatusSubscription = {|
  variables: AddonStatusSubscription$variables,
  response: AddonStatusSubscription$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "filter"
  }
],
v1 = [
  {
    "alias": "addonsState",
    "args": [
      {
        "kind": "Variable",
        "name": "filter",
        "variableName": "filter"
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
    "cacheID": "b7f42d3b9c1ab13f1d78c3048eb9a5db",
    "id": null,
    "metadata": {},
    "name": "AddonStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription AddonStatusSubscription(\n  $filter: ServiceMeshFilter\n) {\n  addonsState: listenToAddonState(filter: $filter) {\n    name\n    owner\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "7cef73e9bfdcc63d12dfe54d2a0f3fbf";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  AddonStatusSubscription$variables,
  AddonStatusSubscription$data,
>*/);

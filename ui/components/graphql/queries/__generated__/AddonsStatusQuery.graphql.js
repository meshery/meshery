/**
 * @generated SignedSource<<81a68e742a5d0e408ab7c56bbbdbf849>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type MeshType = "ALL_MESH" | "INVALID_MESH" | "APP_MESH" | "CITRIX_SERVICE_MESH" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "TRAEFIK_MESH" | "OCTARINE" | "NETWORK_SERVICE_MESH" | "TANZU" | "OPEN_SERVICE_MESH" | "NGINX_SERVICE_MESH" | "CILIUM_SERVICE_MESH" | "%future added value";
export type ServiceMeshFilter = {|
  type?: ?MeshType,
  k8sClusterIDs?: ?$ReadOnlyArray<string>,
|};
export type AddonsStatusQuery$variables = {|
  filter?: ?ServiceMeshFilter,
|};
export type AddonsStatusQuery$data = {|
  +addonsState: $ReadOnlyArray<{|
    +name: string,
    +owner: string,
  |}>,
|};
export type AddonsStatusQuery = {|
  variables: AddonsStatusQuery$variables,
  response: AddonsStatusQuery$data,
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
    "name": "getAvailableAddons",
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
    "name": "AddonsStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AddonsStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "49f96950f100465ccdab36e9903b9281",
    "id": null,
    "metadata": {},
    "name": "AddonsStatusQuery",
    "operationKind": "query",
    "text": "query AddonsStatusQuery(\n  $filter: ServiceMeshFilter\n) {\n  addonsState: getAvailableAddons(filter: $filter) {\n    name\n    owner\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "9cbf0a827a321dead7e3e6d0c2e9cbe7";

module.exports = ((node/*: any*/)/*: Query<
  AddonsStatusQuery$variables,
  AddonsStatusQuery$data,
>*/);

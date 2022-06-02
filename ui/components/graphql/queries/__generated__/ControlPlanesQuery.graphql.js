/**
 * @generated SignedSource<<1b9a3ed5711eaa6029c44d22c1caef33>>
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
export type ControlPlanesQuery$variables = {|
  filter?: ?ServiceMeshFilter,
|};
export type ControlPlanesQuery$data = {|
  +controlPlanesState: $ReadOnlyArray<{|
    +name: string,
    +members: $ReadOnlyArray<{|
      +name: string,
      +version: string,
      +component: string,
      +namespace: string,
    |}>,
  |}>,
|};
export type ControlPlanesQuery = {|
  variables: ControlPlanesQuery$variables,
  response: ControlPlanesQuery$data,
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
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "name",
  "storageKey": null
},
v2 = [
  {
    "alias": "controlPlanesState",
    "args": [
      {
        "kind": "Variable",
        "name": "filter",
        "variableName": "filter"
      }
    ],
    "concreteType": "ControlPlane",
    "kind": "LinkedField",
    "name": "getControlPlanes",
    "plural": true,
    "selections": [
      (v1/*: any*/),
      {
        "alias": null,
        "args": null,
        "concreteType": "ControlPlaneMember",
        "kind": "LinkedField",
        "name": "members",
        "plural": true,
        "selections": [
          (v1/*: any*/),
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "version",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "component",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "namespace",
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
    "name": "ControlPlanesQuery",
    "selections": (v2/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ControlPlanesQuery",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "2a087a7971334b86b29d478c175cd336",
    "id": null,
    "metadata": {},
    "name": "ControlPlanesQuery",
    "operationKind": "query",
    "text": "query ControlPlanesQuery(\n  $filter: ServiceMeshFilter\n) {\n  controlPlanesState: getControlPlanes(filter: $filter) {\n    name\n    members {\n      name\n      version\n      component\n      namespace\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "82f1d5dbf1eec9d253f34bc6b6a7e6f4";

module.exports = ((node/*: any*/)/*: Query<
  ControlPlanesQuery$variables,
  ControlPlanesQuery$data,
>*/);

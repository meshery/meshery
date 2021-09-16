/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL_MESH" | "APP_MESH" | "CITRIX_SERVICE_MESH" | "CONSUL" | "INVALID_MESH" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORK_SERVICE_MESH" | "NGINX_SERVICE_MESH" | "OCTARINE" | "OPEN_SERVICE_MESH" | "TANZU" | "TRAEFIK_MESH" | "%future added value";
export type ServiceMeshFilter = {|
  type?: ?MeshType
|};
export type ControlPlanesQueryVariables = {|
  filter?: ?ServiceMeshFilter
|};
export type ControlPlanesQueryResponse = {|
  +controlPlanesState: $ReadOnlyArray<{|
    +name: string,
    +members: $ReadOnlyArray<{|
      +name: string,
      +version: string,
      +component: string,
      +namespace: string,
    |}>,
  |}>
|};
export type ControlPlanesQuery = {|
  variables: ControlPlanesQueryVariables,
  response: ControlPlanesQueryResponse,
|};
*/


/*
query ControlPlanesQuery(
  $filter: ServiceMeshFilter
) {
  controlPlanesState: getControlPlanes(filter: $filter) {
    name
    members {
      name
      version
      component
      namespace
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
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
// prettier-ignore
(node/*: any*/).hash = '82f1d5dbf1eec9d253f34bc6b6a7e6f4';

module.exports = node;

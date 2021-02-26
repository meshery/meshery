/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type Status = "DISABLED" | "ENABLED" | "UNKNOWN" | "%future added value";
export type ControlPlaneFilter = {|
  type?: ?MeshType
|};
export type meshScanQueryVariables = {|
  filter?: ?ControlPlaneFilter
|};
export type meshScanQueryResponse = {|
  +getControlPlanes: $ReadOnlyArray<{|
    +name: ?MeshType,
    +version: string,
    +members: $ReadOnlyArray<{|
      +component: string,
      +status: ?Status,
      +namespace: string,
    |}>,
  |}>
|};
export type meshScanQuery = {|
  variables: meshScanQueryVariables,
  response: meshScanQueryResponse,
|};
*/


/*
query meshScanQuery(
  $filter: ControlPlaneFilter
) {
  getControlPlanes(filter: $filter) {
    name
    version
    members {
      component
      status
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
v1 = [
  {
    "alias": null,
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
        "name": "version",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "ControlPlaneMember",
        "kind": "LinkedField",
        "name": "members",
        "plural": true,
        "selections": [
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
            "name": "status",
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
    "name": "meshScanQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "meshScanQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "e6e6230af67c35e191c1e10f5451acfd",
    "id": null,
    "metadata": {},
    "name": "meshScanQuery",
    "operationKind": "query",
    "text": "query meshScanQuery(\n  $filter: ControlPlaneFilter\n) {\n  getControlPlanes(filter: $filter) {\n    name\n    version\n    members {\n      component\n      status\n      namespace\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ea95368731b9b984d27bef1748d71bd3';

module.exports = node;

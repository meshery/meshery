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
    "cacheID": "ac72d42262aa447dc8f700419ae50788",
    "id": null,
    "metadata": {},
    "name": "meshScanQuery",
    "operationKind": "query",
    "text": "query meshScanQuery(\n  $filter: ControlPlaneFilter\n) {\n  getControlPlanes(filter: $filter) {\n    name\n    version\n    members {\n      component\n      status\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '848526434c5c8789743bb7a4fe448f2e';

module.exports = node;

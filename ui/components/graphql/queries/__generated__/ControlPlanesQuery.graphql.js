/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type ControlPlaneFilter = {|
  type?: ?MeshType
|};
export type ControlPlanesQueryVariables = {|
  filter?: ?ControlPlaneFilter
|};
export type ControlPlanesQueryResponse = {|
  +controlPlanesState: $ReadOnlyArray<{|
    +name: ?MeshType,
    +members: $ReadOnlyArray<{|
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
  $filter: ControlPlaneFilter
) {
  controlPlanesState: getControlPlanes(filter: $filter) {
    name
    members {
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
v1 = [
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
        "concreteType": "ControlPlaneMember",
        "kind": "LinkedField",
        "name": "members",
        "plural": true,
        "selections": [
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
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ControlPlanesQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "201a5360248793dda91a8787d69c1376",
    "id": null,
    "metadata": {},
    "name": "ControlPlanesQuery",
    "operationKind": "query",
    "text": "query ControlPlanesQuery(\n  $filter: ControlPlaneFilter\n) {\n  controlPlanesState: getControlPlanes(filter: $filter) {\n    name\n    members {\n      version\n      component\n      namespace\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '90d386dc85c45565e02cc6dbceb98da8';

module.exports = node;

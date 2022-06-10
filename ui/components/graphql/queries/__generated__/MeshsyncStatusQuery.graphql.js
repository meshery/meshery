/**
 * @generated SignedSource<<d35797043a537ffb6e902f975fae8edd>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type Status = "ENABLED" | "CONNECTED" | "DISABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type MeshsyncStatusQuery$variables = {|
  k8scontextID: string,
|};
export type MeshsyncStatusQuery$data = {|
  +controller: {|
    +name: string,
    +version: string,
    +status: Status,
  |},
|};
export type MeshsyncStatusQuery = {|
  variables: MeshsyncStatusQuery$variables,
  response: MeshsyncStatusQuery$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "k8scontextID"
  }
],
v1 = [
  {
    "alias": "controller",
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextID",
        "variableName": "k8scontextID"
      }
    ],
    "concreteType": "OperatorControllerStatus",
    "kind": "LinkedField",
    "name": "getMeshsyncStatus",
    "plural": false,
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
        "kind": "ScalarField",
        "name": "status",
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
    "name": "MeshsyncStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshsyncStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "88b5a4baae947adb519f9e528b995fca",
    "id": null,
    "metadata": {},
    "name": "MeshsyncStatusQuery",
    "operationKind": "query",
    "text": "query MeshsyncStatusQuery(\n  $k8scontextID: String!\n) {\n  controller: getMeshsyncStatus(k8scontextID: $k8scontextID) {\n    name\n    version\n    status\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "57b05e27fe69c3f807a81378bcbc5471";

module.exports = ((node/*: any*/)/*: Query<
  MeshsyncStatusQuery$variables,
  MeshsyncStatusQuery$data,
>*/);

/**
 * @generated SignedSource<<f34592f1e9202c68f218d97eb62a2091>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type Status = "ENABLED" | "CONNECTED" | "DISABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type DeployMeshSyncQuery$variables = {|
  k8scontextID: string,
|};
export type DeployMeshSyncQuery$data = {|
  +deployMeshsync: Status,
|};
export type DeployMeshSyncQuery = {|
  variables: DeployMeshSyncQuery$variables,
  response: DeployMeshSyncQuery$data,
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
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextID",
        "variableName": "k8scontextID"
      }
    ],
    "kind": "ScalarField",
    "name": "deployMeshsync",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "DeployMeshSyncQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "DeployMeshSyncQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "0fabd55f7b3f104df0c380e8296c421b",
    "id": null,
    "metadata": {},
    "name": "DeployMeshSyncQuery",
    "operationKind": "query",
    "text": "query DeployMeshSyncQuery(\n  $k8scontextID: String!\n) {\n  deployMeshsync(k8scontextID: $k8scontextID)\n}\n"
  }
};
})();

(node/*: any*/).hash = "ade80c586a9696106576404cbf968aaf";

module.exports = ((node/*: any*/)/*: Query<
  DeployMeshSyncQuery$variables,
  DeployMeshSyncQuery$data,
>*/);

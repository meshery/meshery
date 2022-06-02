/**
 * @generated SignedSource<<068116f3e3fc4f5a4d9daf7dfdf24001>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type Status = "ENABLED" | "CONNECTED" | "DISABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type MeshSyncStatusSubscription$variables = {|
  k8scontextID: string,
|};
export type MeshSyncStatusSubscription$data = {|
  +listenToMeshSyncEvents: {|
    +name: string,
    +status: Status,
    +version: string,
    +error: ?{|
      +code: string,
      +description: string,
    |},
  |},
|};
export type MeshSyncStatusSubscription = {|
  variables: MeshSyncStatusSubscription$variables,
  response: MeshSyncStatusSubscription$data,
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
    "concreteType": "OperatorControllerStatus",
    "kind": "LinkedField",
    "name": "listenToMeshSyncEvents",
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
        "name": "status",
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
        "concreteType": "Error",
        "kind": "LinkedField",
        "name": "error",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "code",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "description",
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
    "name": "MeshSyncStatusSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshSyncStatusSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "f454c7301e9b17565da71991433c215c",
    "id": null,
    "metadata": {},
    "name": "MeshSyncStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription MeshSyncStatusSubscription(\n  $k8scontextID: String!\n) {\n  listenToMeshSyncEvents(k8scontextID: $k8scontextID) {\n    name\n    status\n    version\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "b2826bba5296fa543072c0320de22d65";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  MeshSyncStatusSubscription$variables,
  MeshSyncStatusSubscription$data,
>*/);

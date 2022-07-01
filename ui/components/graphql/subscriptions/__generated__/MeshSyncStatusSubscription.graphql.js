/**
 * @generated SignedSource<<0d298ea50d3a9a725b1cd2e0d0ef5854>>
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
  k8scontextIDs?: ?$ReadOnlyArray<string>,
|};
export type MeshSyncStatusSubscription$data = {|
  +listenToMeshSyncEvents: ?{|
    +contextID: string,
    +OperatorControllerStatus: {|
      +name: string,
      +status: Status,
      +version: string,
      +error: ?{|
        +code: string,
        +description: string,
      |},
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
    "name": "k8scontextIDs"
  }
],
v1 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextIDs",
        "variableName": "k8scontextIDs"
      }
    ],
    "concreteType": "OperatorControllerStatusPerK8sContext",
    "kind": "LinkedField",
    "name": "listenToMeshSyncEvents",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "contextID",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "OperatorControllerStatus",
        "kind": "LinkedField",
        "name": "OperatorControllerStatus",
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
    "cacheID": "cbe09d18ad6045c56617e15888916d4a",
    "id": null,
    "metadata": {},
    "name": "MeshSyncStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription MeshSyncStatusSubscription(\n  $k8scontextIDs: [String!]\n) {\n  listenToMeshSyncEvents(k8scontextIDs: $k8scontextIDs) {\n    contextID\n    OperatorControllerStatus {\n      name\n      status\n      version\n      error {\n        code\n        description\n      }\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "ad82236368e06dbae0e9f6008f6dc032";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  MeshSyncStatusSubscription$variables,
  MeshSyncStatusSubscription$data,
>*/);

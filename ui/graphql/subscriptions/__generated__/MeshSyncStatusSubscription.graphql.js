/**
 * @generated SignedSource<<16aa4fb154990d0451a4d469d4e27328>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

var node = (function(){
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

node.hash = "ad82236368e06dbae0e9f6008f6dc032";

module.exports = node;

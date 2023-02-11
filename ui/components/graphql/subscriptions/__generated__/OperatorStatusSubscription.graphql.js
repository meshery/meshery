/**
 * @generated SignedSource<<25878a5479eea4a4697a737a7eb27f59>>
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
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "status",
  "storageKey": null
},
v2 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "version",
  "storageKey": null
},
v3 = [
  {
    "alias": "operator",
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextIDs",
        "variableName": "k8scontextIDs"
      }
    ],
    "concreteType": "OperatorStatusPerK8sContext",
    "kind": "LinkedField",
    "name": "listenToOperatorState",
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
        "concreteType": "OperatorStatus",
        "kind": "LinkedField",
        "name": "operatorStatus",
        "plural": false,
        "selections": [
          (v1/*: any*/),
          (v2/*: any*/),
          {
            "alias": null,
            "args": null,
            "concreteType": "OperatorControllerStatus",
            "kind": "LinkedField",
            "name": "controllers",
            "plural": true,
            "selections": [
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "name",
                "storageKey": null
              },
              (v2/*: any*/),
              (v1/*: any*/)
            ],
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
    "name": "OperatorStatusSubscription",
    "selections": (v3/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "OperatorStatusSubscription",
    "selections": (v3/*: any*/)
  },
  "params": {
    "cacheID": "8a1bb7537702e43181f395fdcaeceaa0",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription OperatorStatusSubscription(\n  $k8scontextIDs: [String!]\n) {\n  operator: listenToOperatorState(k8scontextIDs: $k8scontextIDs) {\n    contextID\n    operatorStatus {\n      status\n      version\n      controllers {\n        name\n        version\n        status\n      }\n      error {\n        code\n        description\n      }\n    }\n  }\n}\n"
  }
};
})();

node.hash = "5728620b5666bd13a414080f9d90778e";

module.exports = node;

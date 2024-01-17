/**
 * @generated SignedSource<<0b864e564b4ac213932808374e940c88>>
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
    "concreteType": "MesheryControllersStatusListItem",
    "kind": "LinkedField",
    "name": "subscribeMesheryControllersStatus",
    "plural": true,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "contextId",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "controller",
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
    "name": "MesheryControllersStatusSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MesheryControllersStatusSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "d13cc9027c16da162d9414ab27148cf5",
    "id": null,
    "metadata": {},
    "name": "MesheryControllersStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription MesheryControllersStatusSubscription(\n  $k8scontextIDs: [String!]\n) {\n  subscribeMesheryControllersStatus(k8scontextIDs: $k8scontextIDs) {\n    contextId\n    controller\n    status\n    version\n  }\n}\n"
  }
};
})();

node.hash = "b5d4c39af94cb255a592948102c3abb5";

module.exports = node;

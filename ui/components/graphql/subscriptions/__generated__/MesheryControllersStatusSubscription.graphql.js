/**
 * @generated SignedSource<<23f539688b80e52960ef747e8e37a964>>
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
    "name": "connectionIDs"
  }
],
v1 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "connectionIDs",
        "variableName": "connectionIDs"
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
        "name": "connectionID",
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
    "cacheID": "534d4b9391826e634ad2ab05b1a30e75",
    "id": null,
    "metadata": {},
    "name": "MesheryControllersStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription MesheryControllersStatusSubscription(\n  $connectionIDs: [String!]\n) {\n  subscribeMesheryControllersStatus(connectionIDs: $connectionIDs) {\n    connectionID\n    controller\n    status\n    version\n  }\n}\n"
  }
};
})();

node.hash = "57fdebbc7ad327ad4b1607c9fa7df9e2";

module.exports = node;

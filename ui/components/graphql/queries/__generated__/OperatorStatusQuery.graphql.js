/**
 * @generated SignedSource<<72371de5e888cd09251acd22c1310b50>>
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
    "name": "connectionID"
  }
],
v1 = [
  {
    "alias": "operator",
    "args": [
      {
        "kind": "Variable",
        "name": "connectionID",
        "variableName": "connectionID"
      }
    ],
    "concreteType": "MesheryControllersStatusListItem",
    "kind": "LinkedField",
    "name": "getOperatorStatus",
    "plural": false,
    "selections": [
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
        "name": "controller",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "connectionID",
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
    "name": "OperatorStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "OperatorStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "1edf459162cb7e0a916856a8bb8649e2",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusQuery",
    "operationKind": "query",
    "text": "query OperatorStatusQuery(\n  $connectionID: String!\n) {\n  operator: getOperatorStatus(connectionID: $connectionID) {\n    status\n    controller\n    connectionID\n  }\n}\n"
  }
};
})();

node.hash = "519b77a47c22930c9d148d1e4f4011cb";

module.exports = node;

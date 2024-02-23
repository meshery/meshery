/**
 * @generated SignedSource<<df91191dd49a559f7e6e4d8e3bc44adc>>
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
        "name": "contextId",
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
    "cacheID": "29fa4a0ce4b1961decf629f90ebf0604",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusQuery",
    "operationKind": "query",
    "text": "query OperatorStatusQuery(\n  $connectionID: String!\n) {\n  operator: getOperatorStatus(connectionID: $connectionID) {\n    status\n    controller\n    contextId\n  }\n}\n"
  }
};
})();

node.hash = "68141e509e00cc3bb45906a0f0aae24b";

module.exports = node;

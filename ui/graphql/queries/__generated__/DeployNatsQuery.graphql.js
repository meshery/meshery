/**
 * @generated SignedSource<<5a4a6bd490e55d9afa0faa4b6ebc99a9>>
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
    "name": "connectToNats",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "DeployNatsQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "DeployNatsQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "76180bca9f0c5921831b5c004d7e2679",
    "id": null,
    "metadata": {},
    "name": "DeployNatsQuery",
    "operationKind": "query",
    "text": "query DeployNatsQuery(\n  $k8scontextID: String!\n) {\n  connectToNats(k8scontextID: $k8scontextID)\n}\n"
  }
};
})();

node.hash = "36be030e0bed096e2cfaef1101bfaa27";

module.exports = node;

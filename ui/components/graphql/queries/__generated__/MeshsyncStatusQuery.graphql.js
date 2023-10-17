/**
 * @generated SignedSource<<89d91dbbc41da00eed2c37155a5b5bed>>
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
    "alias": "controller",
    "args": [
      {
        "kind": "Variable",
        "name": "connectionID",
        "variableName": "connectionID"
      }
    ],
    "concreteType": "OperatorControllerStatus",
    "kind": "LinkedField",
    "name": "getMeshsyncStatus",
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
        "name": "version",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "status",
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
    "name": "MeshsyncStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshsyncStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "1b2cb7cffaec7ee51423e6f26c1e6ae4",
    "id": null,
    "metadata": {},
    "name": "MeshsyncStatusQuery",
    "operationKind": "query",
    "text": "query MeshsyncStatusQuery(\n  $connectionID: String!\n) {\n  controller: getMeshsyncStatus(connectionID: $connectionID) {\n    name\n    version\n    status\n  }\n}\n"
  }
};
})();

node.hash = "a9ff2d13ffaf332f9be4ca12e09bd7f9";

module.exports = node;

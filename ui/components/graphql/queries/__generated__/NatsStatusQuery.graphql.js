/**
 * @generated SignedSource<<49e3869624d07be085d19ab99c102219>>
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
    "name": "getNatsStatus",
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
    "name": "NatsStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "NatsStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "f232d0001a24e998f5d993786cfed5f6",
    "id": null,
    "metadata": {},
    "name": "NatsStatusQuery",
    "operationKind": "query",
    "text": "query NatsStatusQuery(\n  $connectionID: String!\n) {\n  controller: getNatsStatus(connectionID: $connectionID) {\n    name\n    version\n    status\n  }\n}\n"
  }
};
})();

node.hash = "10628d273554398f8d127aaa764a94fd";

module.exports = node;

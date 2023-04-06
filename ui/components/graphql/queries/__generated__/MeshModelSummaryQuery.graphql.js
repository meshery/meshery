/**
 * @generated SignedSource<<fc2ef860d09573a73d5b766e731a9d5a>>
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
    "name": "selector"
  }
],
v1 = [
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
    "name": "count",
    "storageKey": null
  }
],
v2 = [
  {
    "alias": "meshmodelSummary",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "MeshModelSummary",
    "kind": "LinkedField",
    "name": "getMeshModelSummary",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "MeshModelComponent",
        "kind": "LinkedField",
        "name": "components",
        "plural": true,
        "selections": (v1/*: any*/),
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "MeshModelRelationship",
        "kind": "LinkedField",
        "name": "relationships",
        "plural": true,
        "selections": (v1/*: any*/),
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
    "name": "MeshModelSummaryQuery",
    "selections": (v2/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshModelSummaryQuery",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "fb8718ecc59096bb7a68d0479193aaab",
    "id": null,
    "metadata": {},
    "name": "MeshModelSummaryQuery",
    "operationKind": "query",
    "text": "query MeshModelSummaryQuery(\n  $selector: MeshModelSummarySelector!\n) {\n  meshmodelSummary: getMeshModelSummary(selector: $selector) {\n    components {\n      name\n      count\n    }\n    relationships {\n      name\n      count\n    }\n  }\n}\n"
  }
};
})();

node.hash = "830ea0f72a52a2f3419b0796d9b3a562";

module.exports = node;

/**
 * @generated SignedSource<<250f6fec42f64cdf2df9bd2be3190563>>
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
    "alias": "catalogPatterns",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "CatalogPattern",
    "kind": "LinkedField",
    "name": "fetchPatternCatalogContent",
    "plural": true,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "id",
        "storageKey": null
      },
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
        "name": "user_id",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "pattern_file",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "visibility",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "catalog_data",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "created_at",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "updated_at",
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
    "name": "CatalogPatternQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "CatalogPatternQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "5819e3c1959f0e6ce5329f2a222a499d",
    "id": null,
    "metadata": {},
    "name": "CatalogPatternQuery",
    "operationKind": "query",
    "text": "query CatalogPatternQuery(\n  $selector: CatalogSelector!\n) {\n  catalogPatterns: fetchPatternCatalogContent(selector: $selector) {\n    id\n    name\n    user_id\n    pattern_file\n    visibility\n    catalog_data\n    created_at\n    updated_at\n  }\n}\n"
  }
};
})();

node.hash = "3662bdbf1b55f72dec9757315e54e8ab";

module.exports = node;

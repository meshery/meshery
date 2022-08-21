/**
 * @generated SignedSource<<e64af0eb31dd5861d3750c81516afd70>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type CatalogSelector = {|
  search: string,
  order: string,
|};
export type CatalogPatternQuery$variables = {|
  selector: CatalogSelector,
|};
export type CatalogPatternQuery$data = {|
  +catlogPatterns: $ReadOnlyArray<{|
    +id: string,
    +name: string,
    +user_id: string,
    +pattern_file: string,
    +visibility: any,
    +catalog_data: any,
    +created_at: ?string,
    +updated_at: ?string,
  |}>,
|};
export type CatalogPatternQuery = {|
  variables: CatalogPatternQuery$variables,
  response: CatalogPatternQuery$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "selector"
  }
],
v1 = [
  {
    "alias": "catlogPatterns",
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
    "cacheID": "3c26483019645508156e0cf1db15058e",
    "id": null,
    "metadata": {},
    "name": "CatalogPatternQuery",
    "operationKind": "query",
    "text": "query CatalogPatternQuery(\n  $selector: CatalogSelector!\n) {\n  catlogPatterns: fetchPatternCatalogContent(selector: $selector) {\n    id\n    name\n    user_id\n    pattern_file\n    visibility\n    catalog_data\n    created_at\n    updated_at\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "2e33f33217c049cd49610319b67e3687";

module.exports = ((node/*: any*/)/*: Query<
  CatalogPatternQuery$variables,
  CatalogPatternQuery$data,
>*/);

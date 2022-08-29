/**
 * @generated SignedSource<<b4d3de0d370e97e088d213d9ad49f6b6>>
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
export type CatalogFilterQuery$variables = {|
  selector: CatalogSelector,
|};
export type CatalogFilterQuery$data = {|
  +catlogFilters: $ReadOnlyArray<{|
    +id: string,
    +name: string,
    +user_id: string,
    +filter_file: string,
    +visibility: any,
    +catalog_data: any,
    +created_at: ?string,
    +updated_at: ?string,
  |}>,
|};
export type CatalogFilterQuery = {|
  variables: CatalogFilterQuery$variables,
  response: CatalogFilterQuery$data,
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
    "alias": "catlogFilters",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "CatalogFilter",
    "kind": "LinkedField",
    "name": "fetchFilterCatalogContent",
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
        "name": "filter_file",
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
    "name": "CatalogFilterQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "CatalogFilterQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "c8ef8b102fc74283be63caed6f4b9426",
    "id": null,
    "metadata": {},
    "name": "CatalogFilterQuery",
    "operationKind": "query",
    "text": "query CatalogFilterQuery(\n  $selector: CatalogSelector!\n) {\n  catlogFilters: fetchFilterCatalogContent(selector: $selector) {\n    id\n    name\n    user_id\n    filter_file\n    visibility\n    catalog_data\n    created_at\n    updated_at\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "7432aee8a3d6aa20f0277e24e574df09";

module.exports = ((node/*: any*/)/*: Query<
  CatalogFilterQuery$variables,
  CatalogFilterQuery$data,
>*/);

/**
 * @generated SignedSource<<84d9394ef265babad032c478ee6ac59a>>
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
  +catalogFilters: $ReadOnlyArray<{|
    +id: string,
    +name: string,
    +user_id: string,
    +filter_file: string,
    +visibility: string,
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
    "alias": "catalogFilters",
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
    "cacheID": "ac1121b539ce2cffd452177869b631dc",
    "id": null,
    "metadata": {},
    "name": "CatalogFilterQuery",
    "operationKind": "query",
    "text": "query CatalogFilterQuery(\n  $selector: CatalogSelector!\n) {\n  catalogFilters: fetchFilterCatalogContent(selector: $selector) {\n    id\n    name\n    user_id\n    filter_file\n    visibility\n    catalog_data\n    created_at\n    updated_at\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "6a8f76476c5d44a58d465d9d7106e6b2";

module.exports = ((node/*: any*/)/*: Query<
  CatalogFilterQuery$variables,
  CatalogFilterQuery$data,
>*/);

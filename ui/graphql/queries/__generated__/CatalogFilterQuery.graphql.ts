/**
 * @generated SignedSource<<4cd48e1b5e8227c4329832a11db0c6c0>>
 * @lightSyntaxTransform
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type CatalogSelector = {
  class?: ReadonlyArray<string | null | undefined> | null | undefined;
  metrics?: string | null | undefined;
  order: string;
  orgID?: ReadonlyArray<string | null | undefined> | null | undefined;
  page: string;
  pagesize: string;
  patternType?: ReadonlyArray<string | null | undefined> | null | undefined;
  populate?: ReadonlyArray<string | null | undefined> | null | undefined;
  search: string;
  technology?: ReadonlyArray<string | null | undefined> | null | undefined;
  userid?: ReadonlyArray<string | null | undefined> | null | undefined;
  workspaceID?: ReadonlyArray<string | null | undefined> | null | undefined;
};
export type CatalogFilterQuery$variables = {
  selector: CatalogSelector;
};
export type CatalogFilterQuery$data = {
  readonly catalogFilters: ReadonlyArray<{
    readonly catalog_data: any | null | undefined;
    readonly created_at: string | null | undefined;
    readonly filter_file: string;
    readonly filter_resource: string;
    readonly id: string;
    readonly name: string;
    readonly owner: string;
    readonly updated_at: string | null | undefined;
    readonly visibility: string;
  }>;
};
export type CatalogFilterQuery = {
  response: CatalogFilterQuery$data;
  variables: CatalogFilterQuery$variables;
};

const node: ConcreteRequest = (function(){
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
        "name": "owner",
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
        "name": "filter_resource",
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
    "argumentDefinitions": (v0/*:: as any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "CatalogFilterQuery",
    "selections": (v1/*:: as any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*:: as any*/),
    "kind": "Operation",
    "name": "CatalogFilterQuery",
    "selections": (v1/*:: as any*/)
  },
  "params": {
    "cacheID": "54a402060d50fc439524839d8ad8c47b",
    "id": null,
    "metadata": {},
    "name": "CatalogFilterQuery",
    "operationKind": "query",
    "text": "query CatalogFilterQuery(\n  $selector: CatalogSelector!\n) {\n  catalogFilters: fetchFilterCatalogContent(selector: $selector) {\n    id\n    name\n    owner\n    filter_file\n    filter_resource\n    visibility\n    catalog_data\n    created_at\n    updated_at\n  }\n}\n"
  }
};
})();

(node as any).hash = "bfa9963f8c464994cd5a0015cfb95f06";

export default node;

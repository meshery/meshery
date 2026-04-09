/**
 * @generated SignedSource<<126d7055ae416a8abf815f9ea968c5cf>>
 * @lightSyntaxTransform
 * @nogrep
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
    readonly updated_at: string | null | undefined;
    readonly user_id: string;
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
    "cacheID": "2178c2622b2e672bcf288e1004c939a0",
    "id": null,
    "metadata": {},
    "name": "CatalogFilterQuery",
    "operationKind": "query",
    "text": "query CatalogFilterQuery(\n  $selector: CatalogSelector!\n) {\n  catalogFilters: fetchFilterCatalogContent(selector: $selector) {\n    id\n    name\n    user_id\n    filter_file\n    filter_resource\n    visibility\n    catalog_data\n    created_at\n    updated_at\n  }\n}\n"
  }
};
})();

(node as any).hash = "391a34d0da3dfd429e7a8a335e07930b";

export default node;

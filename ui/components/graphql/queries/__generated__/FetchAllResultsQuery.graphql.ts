/**
 * @generated SignedSource<<102229c5145a1d2b2d6cecc811bb618e>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type PageFilter = {
  from?: string | null | undefined;
  metrics?: string | null | undefined;
  order?: string | null | undefined;
  page: string;
  pageSize: string;
  populate?: ReadonlyArray<string | null | undefined> | null | undefined;
  search?: string | null | undefined;
  to?: string | null | undefined;
  updated_after?: string | null | undefined;
  visibility?: ReadonlyArray<string> | null | undefined;
};
export type FetchAllResultsQuery$variables = {
  selector: PageFilter;
};
export type FetchAllResultsQuery$data = {
  readonly fetchAllResults: {
    readonly page: number;
    readonly page_size: number;
    readonly results: ReadonlyArray<{
      readonly created_at: string | null | undefined;
      readonly mesh: string | null | undefined;
      readonly meshery_id: string | null | undefined;
      readonly name: string | null | undefined;
      readonly performance_profile: string | null | undefined;
      readonly runner_results: any | null | undefined;
      readonly server_metrics: string | null | undefined;
      readonly test_id: string | null | undefined;
      readonly test_start_time: string | null | undefined;
      readonly updated_at: string | null | undefined;
      readonly user_id: string | null | undefined;
    } | null | undefined> | null | undefined;
    readonly total_count: number;
  };
};
export type FetchAllResultsQuery = {
  response: FetchAllResultsQuery$data;
  variables: FetchAllResultsQuery$variables;
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
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "PerfPageResult",
    "kind": "LinkedField",
    "name": "fetchAllResults",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "page",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "page_size",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "total_count",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "MesheryResult",
        "kind": "LinkedField",
        "name": "results",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "meshery_id",
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
            "name": "mesh",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "performance_profile",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "test_id",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "server_metrics",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "test_start_time",
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
            "name": "user_id",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "updated_at",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "runner_results",
            "storageKey": null
          }
        ],
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
    "name": "FetchAllResultsQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "FetchAllResultsQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "0770a3af73084b2362bfb7bad68fbeb4",
    "id": null,
    "metadata": {},
    "name": "FetchAllResultsQuery",
    "operationKind": "query",
    "text": "query FetchAllResultsQuery(\n  $selector: PageFilter!\n) {\n  fetchAllResults(selector: $selector) {\n    page\n    page_size\n    total_count\n    results {\n      meshery_id\n      name\n      mesh\n      performance_profile\n      test_id\n      server_metrics\n      test_start_time\n      created_at\n      user_id\n      updated_at\n      runner_results\n    }\n  }\n}\n"
  }
};
})();

(node as any).hash = "0aede14fffa5004109535cbff9b07687";

export default node;

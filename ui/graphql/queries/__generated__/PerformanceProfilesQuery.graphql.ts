/**
 * @generated SignedSource<<3a91d48a9bbe2cc55f168d5b23085cfe>>
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
export type PerformanceProfilesQuery$variables = {
  selector: PageFilter;
};
export type PerformanceProfilesQuery$data = {
  readonly getPerformanceProfiles: {
    readonly page: number;
    readonly page_size: number;
    readonly profiles: ReadonlyArray<{
      readonly concurrent_request: number;
      readonly content_type: string | null | undefined;
      readonly created_at: string | null | undefined;
      readonly duration: string;
      readonly endpoints: ReadonlyArray<string | null | undefined> | null | undefined;
      readonly id: string;
      readonly last_run: string | null | undefined;
      readonly load_generators: ReadonlyArray<string | null | undefined> | null | undefined;
      readonly metadata: any | null | undefined;
      readonly name: string | null | undefined;
      readonly qps: number | null | undefined;
      readonly request_body: string | null | undefined;
      readonly request_cookies: string | null | undefined;
      readonly request_headers: string | null | undefined;
      readonly service_mesh: string | null | undefined;
      readonly total_results: number | null | undefined;
      readonly updated_at: string | null | undefined;
      readonly user_id: string;
    } | null | undefined> | null | undefined;
    readonly total_count: number;
  };
};
export type PerformanceProfilesQuery = {
  response: PerformanceProfilesQuery$data;
  variables: PerformanceProfilesQuery$variables;
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
    "concreteType": "PerfPageProfiles",
    "kind": "LinkedField",
    "name": "getPerformanceProfiles",
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
        "concreteType": "PerfProfile",
        "kind": "LinkedField",
        "name": "profiles",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "concurrent_request",
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
            "name": "duration",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "endpoints",
            "storageKey": null
          },
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
            "name": "last_run",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "load_generators",
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
            "name": "qps",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "total_results",
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
            "name": "user_id",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "request_body",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "request_cookies",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "request_headers",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "content_type",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "service_mesh",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "metadata",
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
    "name": "PerformanceProfilesQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "PerformanceProfilesQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "05fcaffb37eda9d52f8bd45bd8da5550",
    "id": null,
    "metadata": {},
    "name": "PerformanceProfilesQuery",
    "operationKind": "query",
    "text": "query PerformanceProfilesQuery(\n  $selector: PageFilter!\n) {\n  getPerformanceProfiles(selector: $selector) {\n    page\n    page_size\n    total_count\n    profiles {\n      concurrent_request\n      created_at\n      duration\n      endpoints\n      id\n      last_run\n      load_generators\n      name\n      qps\n      total_results\n      updated_at\n      user_id\n      request_body\n      request_cookies\n      request_headers\n      content_type\n      service_mesh\n      metadata\n    }\n  }\n}\n"
  }
};
})();

(node as any).hash = "a901902df4422ec63fa999ca605a644b";

export default node;

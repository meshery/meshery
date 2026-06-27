/**
 * @generated SignedSource<<92d3e163ce83d999377e62b54af81186>>
 * @lightSyntaxTransform
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
export type PerformanceProfilesSubscription$variables = {
  selector: PageFilter;
};
export type PerformanceProfilesSubscription$data = {
  readonly subscribePerfProfiles: {
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
      readonly owner: string;
      readonly qps: number | null | undefined;
      readonly request_body: string | null | undefined;
      readonly request_cookies: string | null | undefined;
      readonly request_headers: string | null | undefined;
      readonly service_mesh: string | null | undefined;
      readonly total_results: number | null | undefined;
      readonly updated_at: string | null | undefined;
    } | null | undefined> | null | undefined;
    readonly total_count: number;
  };
};
export type PerformanceProfilesSubscription = {
  response: PerformanceProfilesSubscription$data;
  variables: PerformanceProfilesSubscription$variables;
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
    "name": "subscribePerfProfiles",
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
            "name": "owner",
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
    "argumentDefinitions": (v0/*:: as any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "PerformanceProfilesSubscription",
    "selections": (v1/*:: as any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*:: as any*/),
    "kind": "Operation",
    "name": "PerformanceProfilesSubscription",
    "selections": (v1/*:: as any*/)
  },
  "params": {
    "cacheID": "0f4908d8d1fd2fd8e9d32903a8ea0c86",
    "id": null,
    "metadata": {},
    "name": "PerformanceProfilesSubscription",
    "operationKind": "subscription",
    "text": "subscription PerformanceProfilesSubscription(\n  $selector: PageFilter!\n) {\n  subscribePerfProfiles(selector: $selector) {\n    page\n    page_size\n    total_count\n    profiles {\n      concurrent_request\n      created_at\n      duration\n      endpoints\n      id\n      last_run\n      load_generators\n      name\n      qps\n      total_results\n      updated_at\n      owner\n      request_body\n      request_cookies\n      request_headers\n      content_type\n      service_mesh\n      metadata\n    }\n  }\n}\n"
  }
};
})();

(node as any).hash = "96a468adaa1d6a43f86595c3c3902428";

export default node;

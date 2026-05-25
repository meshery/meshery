/**
 * @generated SignedSource<<990d9b897fc0ca2e397f65be62f59b3d>>
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
export type PerformanceResultQuery$variables = {
  profileID: string;
  selector: PageFilter;
};
export type PerformanceResultQuery$data = {
  readonly fetchResults: {
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
export type PerformanceResultQuery = {
  response: PerformanceResultQuery$data;
  variables: PerformanceResultQuery$variables;
};

const node: ConcreteRequest = (function(){
var v0 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "profileID"
},
v1 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "selector"
},
v2 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "profileID",
        "variableName": "profileID"
      },
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "PerfPageResult",
    "kind": "LinkedField",
    "name": "fetchResults",
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
    "argumentDefinitions": [
      (v0/*: any*/),
      (v1/*: any*/)
    ],
    "kind": "Fragment",
    "metadata": null,
    "name": "PerformanceResultQuery",
    "selections": (v2/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v1/*: any*/),
      (v0/*: any*/)
    ],
    "kind": "Operation",
    "name": "PerformanceResultQuery",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "cb5e2072afbc31af649ca22bc7e5d4d6",
    "id": null,
    "metadata": {},
    "name": "PerformanceResultQuery",
    "operationKind": "query",
    "text": "query PerformanceResultQuery(\n  $selector: PageFilter!\n  $profileID: String!\n) {\n  fetchResults(selector: $selector, profileID: $profileID) {\n    page\n    page_size\n    total_count\n    results {\n      meshery_id\n      name\n      mesh\n      performance_profile\n      test_id\n      server_metrics\n      test_start_time\n      created_at\n      user_id\n      updated_at\n      runner_results\n    }\n  }\n}\n"
  }
};
})();

(node as any).hash = "e42d1c8529e951f7ad055eab6db359b8";

export default node;

/**
 * @generated SignedSource<<5a4a0e50be67f1046811ed99dba67ef1>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type PageFilter = {|
  page: string,
  pageSize: string,
  order?: ?string,
  search?: ?string,
  from?: ?string,
  to?: ?string,
|};
export type PerformanceProfilesQuery$variables = {|
  selector: PageFilter,
|};
export type PerformanceProfilesQuery$data = {|
  +getPerformanceProfiles: {|
    +page: number,
    +page_size: number,
    +total_count: number,
    +profiles: ?$ReadOnlyArray<?{|
      +concurrent_request: number,
      +created_at: ?string,
      +duration: string,
      +endpoints: ?$ReadOnlyArray<?string>,
      +id: string,
      +last_run: ?string,
      +load_generators: ?$ReadOnlyArray<?string>,
      +name: ?string,
      +qps: ?number,
      +total_results: ?number,
      +updated_at: ?string,
      +user_id: string,
      +request_body: ?string,
      +request_cookies: ?string,
      +request_headers: ?string,
      +content_type: ?string,
      +service_mesh: ?string,
    |}>,
  |},
|};
export type PerformanceProfilesQuery = {|
  variables: PerformanceProfilesQuery$variables,
  response: PerformanceProfilesQuery$data,
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
    "cacheID": "9d7eded079e0a1f23d58b8f400fbf1f7",
    "id": null,
    "metadata": {},
    "name": "PerformanceProfilesQuery",
    "operationKind": "query",
    "text": "query PerformanceProfilesQuery(\n  $selector: PageFilter!\n) {\n  getPerformanceProfiles(selector: $selector) {\n    page\n    page_size\n    total_count\n    profiles {\n      concurrent_request\n      created_at\n      duration\n      endpoints\n      id\n      last_run\n      load_generators\n      name\n      qps\n      total_results\n      updated_at\n      user_id\n      request_body\n      request_cookies\n      request_headers\n      content_type\n      service_mesh\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "8c86c93431616aef124858b289e06a4f";

module.exports = ((node/*: any*/)/*: Query<
  PerformanceProfilesQuery$variables,
  PerformanceProfilesQuery$data,
>*/);

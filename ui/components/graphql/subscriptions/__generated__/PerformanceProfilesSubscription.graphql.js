/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PageFilter = {|
  page: string,
  pageSize: string,
  order?: ?string,
  search?: ?string,
  from?: ?string,
  to?: ?string,
|};
export type PerformanceProfilesSubscriptionVariables = {|
  selector: PageFilter
|};
export type PerformanceProfilesSubscriptionResponse = {|
  +subscribePerfProfiles: {|
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
  |}
|};
export type PerformanceProfilesSubscription = {|
  variables: PerformanceProfilesSubscriptionVariables,
  response: PerformanceProfilesSubscriptionResponse,
|};
*/


/*
subscription PerformanceProfilesSubscription(
  $selector: PageFilter!
) {
  subscribePerfProfiles(selector: $selector) {
    page
    page_size
    total_count
    profiles {
      concurrent_request
      created_at
      duration
      endpoints
      id
      last_run
      load_generators
      name
      qps
      total_results
      updated_at
      user_id
      request_body
      request_cookies
      request_headers
      content_type
      service_mesh
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
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
    "name": "PerformanceProfilesSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "PerformanceProfilesSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "44c12a99a803800ffae91a06add3ef0c",
    "id": null,
    "metadata": {},
    "name": "PerformanceProfilesSubscription",
    "operationKind": "subscription",
    "text": "subscription PerformanceProfilesSubscription(\n  $selector: PageFilter!\n) {\n  subscribePerfProfiles(selector: $selector) {\n    page\n    page_size\n    total_count\n    profiles {\n      concurrent_request\n      created_at\n      duration\n      endpoints\n      id\n      last_run\n      load_generators\n      name\n      qps\n      total_results\n      updated_at\n      user_id\n      request_body\n      request_cookies\n      request_headers\n      content_type\n      service_mesh\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '2efb94fd4676d37f9cc54c1c45a8bcdb';

module.exports = node;

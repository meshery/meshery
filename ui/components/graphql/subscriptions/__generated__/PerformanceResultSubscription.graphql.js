/**
 * @generated SignedSource<<961ee8b78a38b33ad9b63126264132db>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type PageFilter = {|
  page: string,
  pageSize: string,
  order?: ?string,
  search?: ?string,
  from?: ?string,
  to?: ?string,
|};
export type PerformanceResultSubscription$variables = {|
  selector: PageFilter,
  profileID: string,
|};
export type PerformanceResultSubscription$data = {|
  +subscribePerfResults: {|
    +page: number,
    +page_size: number,
    +total_count: number,
    +results: ?$ReadOnlyArray<?{|
      +meshery_id: ?string,
      +name: ?string,
      +mesh: ?string,
      +performance_profile: ?string,
      +test_id: ?string,
      +server_metrics: ?string,
      +test_start_time: ?string,
      +created_at: ?string,
      +user_id: ?string,
      +updated_at: ?string,
      +runner_results: ?any,
    |}>,
  |},
|};
export type PerformanceResultSubscription = {|
  variables: PerformanceResultSubscription$variables,
  response: PerformanceResultSubscription$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
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
    "name": "subscribePerfResults",
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
    "name": "PerformanceResultSubscription",
    "selections": (v2/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v1/*: any*/),
      (v0/*: any*/)
    ],
    "kind": "Operation",
    "name": "PerformanceResultSubscription",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "fa344fe9e839347da9809ab6a307034a",
    "id": null,
    "metadata": {},
    "name": "PerformanceResultSubscription",
    "operationKind": "subscription",
    "text": "subscription PerformanceResultSubscription(\n  $selector: PageFilter!\n  $profileID: String!\n) {\n  subscribePerfResults(selector: $selector, profileID: $profileID) {\n    page\n    page_size\n    total_count\n    results {\n      meshery_id\n      name\n      mesh\n      performance_profile\n      test_id\n      server_metrics\n      test_start_time\n      created_at\n      user_id\n      updated_at\n      runner_results\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "4cd62dfb1b40f60dbc14b2ad7ffea32e";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  PerformanceResultSubscription$variables,
  PerformanceResultSubscription$data,
>*/);

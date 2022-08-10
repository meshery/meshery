/**
 * @generated SignedSource<<929bcea13d1700ea379c3f0c28452b8c>>
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
export type PatternsQuery$variables = {|
  selector: PageFilter,
|};
export type PatternsQuery$data = {|
  +fetchPatterns: {|
    +page: number,
    +page_size: number,
    +patterns: ?$ReadOnlyArray<?{|
      +canSupport: boolean,
      +created_at: ?string,
      +errmsg: ?string,
      +id: string,
      +location: {|
        +branch: ?string,
        +host: ?string,
        +path: ?string,
        +type: ?string,
      |},
      +name: string,
      +pattern_file: string,
      +updated_at: ?string,
      +user_id: string,
    |}>,
    +total_count: number,
  |},
|};
export type PatternsQuery = {|
  variables: PatternsQuery$variables,
  response: PatternsQuery$data,
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
    "concreteType": "PatternPageResult",
    "kind": "LinkedField",
    "name": "fetchPatterns",
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
        "concreteType": "PatternResult",
        "kind": "LinkedField",
        "name": "patterns",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "canSupport",
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
            "name": "errmsg",
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
            "concreteType": "Location",
            "kind": "LinkedField",
            "name": "location",
            "plural": false,
            "selections": [
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "branch",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "host",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "path",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "type",
                "storageKey": null
              }
            ],
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
            "name": "pattern_file",
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
          }
        ],
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "total_count",
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
    "name": "PatternsQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "PatternsQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "d3f796bc59feb18d2455c2218ae6eec8",
    "id": null,
    "metadata": {},
    "name": "PatternsQuery",
    "operationKind": "query",
    "text": "query PatternsQuery(\n  $selector: PageFilter!\n) {\n  fetchPatterns(selector: $selector) {\n    page\n    page_size\n    patterns {\n      canSupport\n      created_at\n      errmsg\n      id\n      location {\n        branch\n        host\n        path\n        type\n      }\n      name\n      pattern_file\n      updated_at\n      user_id\n    }\n    total_count\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "b7e08391b5f6608d3158decad80e8614";

module.exports = ((node/*: any*/)/*: Query<
  PatternsQuery$variables,
  PatternsQuery$data,
>*/);

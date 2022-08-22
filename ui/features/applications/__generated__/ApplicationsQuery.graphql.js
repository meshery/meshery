/**
 * @generated SignedSource<<ddc154214672915b9430e50b6d3a8765>>
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
export type ApplicationsQuery$variables = {|
  selector: PageFilter,
|};
export type ApplicationsQuery$data = {|
  +fetchApplications: {|
    +page: number,
    +page_size: number,
    +total_count: number,
    +applications: ?$ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +user_id: string,
      +application_file: string,
      +location: {|
        +branch: ?string,
        +host: ?string,
        +path: ?string,
        +type: ?string,
      |},
      +created_at: ?string,
      +updated_at: ?string,
    |}>,
  |},
|};
export type ApplicationsQuery = {|
  variables: ApplicationsQuery$variables,
  response: ApplicationsQuery$data,
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
    "concreteType": "ApplicationPage",
    "kind": "LinkedField",
    "name": "fetchApplications",
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
        "concreteType": "ApplicationResult",
        "kind": "LinkedField",
        "name": "applications",
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
            "name": "application_file",
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
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "ApplicationsQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ApplicationsQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "5882a09e7cec3dd9b0163293b8f518c6",
    "id": null,
    "metadata": {},
    "name": "ApplicationsQuery",
    "operationKind": "query",
    "text": "query ApplicationsQuery(\n  $selector: PageFilter!\n) {\n  fetchApplications(selector: $selector) {\n    page\n    page_size\n    total_count\n    applications {\n      id\n      name\n      user_id\n      application_file\n      location {\n        branch\n        host\n        path\n        type\n      }\n      created_at\n      updated_at\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "ba8a1ed5408a8812530690ee8b27b740";

module.exports = ((node/*: any*/)/*: Query<
  ApplicationsQuery$variables,
  ApplicationsQuery$data,
>*/);

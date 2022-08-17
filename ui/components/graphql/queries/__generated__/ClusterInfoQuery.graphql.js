/**
 * @generated SignedSource<<b91b1256fcb317e7a7c063384a1ef283>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type ClusterInfoQuery$variables = {|
  k8scontextIDs?: ?$ReadOnlyArray<string>,
|};
export type ClusterInfoQuery$data = {|
  +clusterInfo: {|
    +resources: $ReadOnlyArray<{|
      +kind: string,
      +number: number,
    |}>,
  |},
|};
export type ClusterInfoQuery = {|
  variables: ClusterInfoQuery$variables,
  response: ClusterInfoQuery$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "k8scontextIDs"
  }
],
v1 = [
  {
    "alias": "clusterInfo",
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextIDs",
        "variableName": "k8scontextIDs"
      }
    ],
    "concreteType": "ClusterInfo",
    "kind": "LinkedField",
    "name": "getClusterInfo",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Resource",
        "kind": "LinkedField",
        "name": "resources",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "kind",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "number",
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
    "name": "ClusterInfoQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ClusterInfoQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "5131abb7963f46b954cb00af93196c0d",
    "id": null,
    "metadata": {},
    "name": "ClusterInfoQuery",
    "operationKind": "query",
    "text": "query ClusterInfoQuery(\n  $k8scontextIDs: [String!]\n) {\n  clusterInfo: getClusterInfo(k8scontextIDs: $k8scontextIDs) {\n    resources {\n      kind\n      number\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "87fe4e9d9c0561175fb5e3a8fbc075ce";

module.exports = ((node/*: any*/)/*: Query<
  ClusterInfoQuery$variables,
  ClusterInfoQuery$data,
>*/);

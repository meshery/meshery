/**
 * @generated SignedSource<<78510c1343b7d26369488990fd1716e6>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type ClusterResourcesQuery$variables = {|
  k8scontextIDs?: ?$ReadOnlyArray<string>,
  namespaces?: ?$ReadOnlyArray<string>,
|};
export type ClusterResourcesQuery$data = {|
  +clusterResources: {|
    +resources: $ReadOnlyArray<{|
      +kind: string,
      +count: number,
    |}>,
  |},
|};
export type ClusterResourcesQuery = {|
  variables: ClusterResourcesQuery$variables,
  response: ClusterResourcesQuery$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "k8scontextIDs"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "namespaces"
  }
],
v1 = [
  {
    "alias": "clusterResources",
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextIDs",
        "variableName": "k8scontextIDs"
      },
      {
        "kind": "Variable",
        "name": "namespaces",
        "variableName": "namespaces"
      }
    ],
    "concreteType": "ClusterResources",
    "kind": "LinkedField",
    "name": "getClusterResources",
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
            "name": "count",
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
    "name": "ClusterResourcesQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ClusterResourcesQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "b4049f93cc26eb621c50945342a69b8b",
    "id": null,
    "metadata": {},
    "name": "ClusterResourcesQuery",
    "operationKind": "query",
    "text": "query ClusterResourcesQuery(\n  $k8scontextIDs: [String!]\n  $namespaces: [String!]\n) {\n  clusterResources: getClusterResources(k8scontextIDs: $k8scontextIDs, namespaces: $namespaces) {\n    resources {\n      kind\n      count\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "996dd23bbe7c3ebc88fcf007049c3f06";

module.exports = ((node/*: any*/)/*: Query<
  ClusterResourcesQuery$variables,
  ClusterResourcesQuery$data,
>*/);

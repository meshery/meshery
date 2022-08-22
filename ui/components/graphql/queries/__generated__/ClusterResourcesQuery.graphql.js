/**
 * @generated SignedSource<<1390691b362e5451ff54588040f05cc3>>
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
    "cacheID": "1abdefbc9ce07e8a87c32dbf30edcbdd",
    "id": null,
    "metadata": {},
    "name": "ClusterResourcesQuery",
    "operationKind": "query",
    "text": "query ClusterResourcesQuery(\n  $k8scontextIDs: [String!]\n) {\n  clusterResources: getClusterResources(k8scontextIDs: $k8scontextIDs) {\n    resources {\n      kind\n      count\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "20be3f850535f5fc738c4b759c99215f";

module.exports = ((node/*: any*/)/*: Query<
  ClusterResourcesQuery$variables,
  ClusterResourcesQuery$data,
>*/);

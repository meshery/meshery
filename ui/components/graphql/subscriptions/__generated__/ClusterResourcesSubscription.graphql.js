/**
 * @generated SignedSource<<f8ff39fc0d57fdee68da032e096a60de>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type ClusterResourcesSubscription$variables = {|
  k8scontextIDs?: ?$ReadOnlyArray<string>,
|};
export type ClusterResourcesSubscription$data = {|
  +clusterResources: {|
    +resources: $ReadOnlyArray<{|
      +kind: string,
      +count: number,
    |}>,
  |},
|};
export type ClusterResourcesSubscription = {|
  variables: ClusterResourcesSubscription$variables,
  response: ClusterResourcesSubscription$data,
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
    "name": "subscribeClusterResources",
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
    "name": "ClusterResourcesSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ClusterResourcesSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "08cf7a17c7256b24bf7217396260f4cc",
    "id": null,
    "metadata": {},
    "name": "ClusterResourcesSubscription",
    "operationKind": "subscription",
    "text": "subscription ClusterResourcesSubscription(\n  $k8scontextIDs: [String!]\n) {\n  clusterResources: subscribeClusterResources(k8scontextIDs: $k8scontextIDs) {\n    resources {\n      kind\n      count\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "ab3fce770bc208fe26f2806d070b8995";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  ClusterResourcesSubscription$variables,
  ClusterResourcesSubscription$data,
>*/);

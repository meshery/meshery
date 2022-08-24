/**
 * @generated SignedSource<<9686ae525411d2efdf04211e27fb8023>>
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
  namespaces?: ?$ReadOnlyArray<string>,
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
    "cacheID": "30fcc5c9285042a529056224d20d999b",
    "id": null,
    "metadata": {},
    "name": "ClusterResourcesSubscription",
    "operationKind": "subscription",
    "text": "subscription ClusterResourcesSubscription(\n  $k8scontextIDs: [String!]\n  $namespaces: [String!]\n) {\n  clusterResources: subscribeClusterResources(k8scontextIDs: $k8scontextIDs, namespaces: $namespaces) {\n    resources {\n      kind\n      count\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "d221ceb413e9ea6de94bde64fc229090";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  ClusterResourcesSubscription$variables,
  ClusterResourcesSubscription$data,
>*/);

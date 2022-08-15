/**
 * @generated SignedSource<<3fab72b1d660ff726d690a9b0057c232>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type ClusterInfoSubscription$variables = {|
  k8scontextIDs?: ?$ReadOnlyArray<string>,
|};
export type ClusterInfoSubscription$data = {|
  +clusterInfo: {|
    +resources: $ReadOnlyArray<{|
      +kind: string,
      +number: number,
    |}>,
  |},
|};
export type ClusterInfoSubscription = {|
  variables: ClusterInfoSubscription$variables,
  response: ClusterInfoSubscription$data,
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
    "name": "subscribeClusterInfo",
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
    "name": "ClusterInfoSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ClusterInfoSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "cd90ed04679ef1175d0c20bdddc9cc47",
    "id": null,
    "metadata": {},
    "name": "ClusterInfoSubscription",
    "operationKind": "subscription",
    "text": "subscription ClusterInfoSubscription(\n  $k8scontextIDs: [String!]\n) {\n  clusterInfo: subscribeClusterInfo(k8scontextIDs: $k8scontextIDs) {\n    resources {\n      kind\n      number\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "867b3a2886251da9448c8dfde882c99a";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  ClusterInfoSubscription$variables,
  ClusterInfoSubscription$data,
>*/);

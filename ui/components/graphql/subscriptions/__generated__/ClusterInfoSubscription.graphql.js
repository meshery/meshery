/**
 * @generated SignedSource<<fde52f581f1f0ef747f278cf2e9aa1ff>>
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
  +subscribeClusterInfo: {|
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
    "alias": null,
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
    "cacheID": "1d8aba943dec00a83aab08ae03d818f8",
    "id": null,
    "metadata": {},
    "name": "ClusterInfoSubscription",
    "operationKind": "subscription",
    "text": "subscription ClusterInfoSubscription(\n  $k8scontextIDs: [String!]\n) {\n  subscribeClusterInfo(k8scontextIDs: $k8scontextIDs) {\n    resources {\n      kind\n      number\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "6b46e3d367148a9cd7f2322c6d8092a5";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  ClusterInfoSubscription$variables,
  ClusterInfoSubscription$data,
>*/);

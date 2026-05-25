/**
 * @generated SignedSource<<22fe2142625128f51930ffae62197796>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type ClusterResourcesSubscription$variables = {
  k8scontextIDs?: ReadonlyArray<string> | null | undefined;
  namespace: string;
};
export type ClusterResourcesSubscription$data = {
  readonly clusterResources: {
    readonly resources: ReadonlyArray<{
      readonly count: number;
      readonly kind: string;
    }>;
  };
};
export type ClusterResourcesSubscription = {
  response: ClusterResourcesSubscription$data;
  variables: ClusterResourcesSubscription$variables;
};

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "k8scontextIDs"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "namespace"
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
        "name": "namespace",
        "variableName": "namespace"
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
    "cacheID": "d3c5e48d338830a81ec808ac6acf5d8a",
    "id": null,
    "metadata": {},
    "name": "ClusterResourcesSubscription",
    "operationKind": "subscription",
    "text": "subscription ClusterResourcesSubscription(\n  $k8scontextIDs: [String!]\n  $namespace: String!\n) {\n  clusterResources: subscribeClusterResources(k8scontextIDs: $k8scontextIDs, namespace: $namespace) {\n    resources {\n      kind\n      count\n    }\n  }\n}\n"
  }
};
})();

(node as any).hash = "2722cb19adcd7e786a91856efcff3e1c";

export default node;

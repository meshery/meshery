/**
 * @generated SignedSource<<e11c1691eba99b1c55d832d4aaf0c25c>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

var node = (function(){
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
    "cacheID": "be4502eb8686b7ba1ab73d914f75e299",
    "id": null,
    "metadata": {},
    "name": "ClusterResourcesQuery",
    "operationKind": "query",
    "text": "query ClusterResourcesQuery(\n  $k8scontextIDs: [String!]\n  $namespace: String!\n) {\n  clusterResources: getClusterResources(k8scontextIDs: $k8scontextIDs, namespace: $namespace) {\n    resources {\n      kind\n      count\n    }\n  }\n}\n"
  }
};
})();

node.hash = "df3b0056c596c2f025a87ac56962fc6e";

module.exports = node;

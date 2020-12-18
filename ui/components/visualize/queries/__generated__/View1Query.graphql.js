/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type View1QueryVariables = {|
  nodeID?: ?string,
  showClusters: boolean,
  showNodes: boolean,
  showPods: boolean,
|};
export type View1QueryResponse = {|
  +cluster: ?$ReadOnlyArray<?{|
    +id?: string,
    +clusterNodes: ?$ReadOnlyArray<?{|
      +id?: string,
      +pods: $ReadOnlyArray<?{|
        +id?: string,
        +name?: string,
      |}>,
    |}>,
  |}>
|};
export type View1Query = {|
  variables: View1QueryVariables,
  response: View1QueryResponse,
|};
*/


/*
query View1Query(
  $nodeID: ID
  $showClusters: Boolean!
  $showNodes: Boolean!
  $showPods: Boolean!
) {
  cluster {
    id @include(if: $showClusters)
    clusterNodes(nodeid: $nodeID) {
      id @include(if: $showNodes)
      pods {
        id @include(if: $showPods)
        name @include(if: $showPods)
        id
      }
      id
    }
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "nodeID"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "showClusters"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "showNodes"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "showPods"
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "nodeid",
    "variableName": "nodeID"
  }
],
v2 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v3 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "name",
  "storageKey": null
},
v4 = [
  (v2/*: any*/)
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "View1Query",
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Cluster",
        "kind": "LinkedField",
        "name": "cluster",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": (v1/*: any*/),
            "concreteType": "ClusterNode",
            "kind": "LinkedField",
            "name": "clusterNodes",
            "plural": true,
            "selections": [
              {
                "alias": null,
                "args": null,
                "concreteType": "Pod",
                "kind": "LinkedField",
                "name": "pods",
                "plural": true,
                "selections": [
                  {
                    "condition": "showPods",
                    "kind": "Condition",
                    "passingValue": true,
                    "selections": [
                      (v2/*: any*/),
                      (v3/*: any*/)
                    ]
                  }
                ],
                "storageKey": null
              },
              {
                "condition": "showNodes",
                "kind": "Condition",
                "passingValue": true,
                "selections": (v4/*: any*/)
              }
            ],
            "storageKey": null
          },
          {
            "condition": "showClusters",
            "kind": "Condition",
            "passingValue": true,
            "selections": (v4/*: any*/)
          }
        ],
        "storageKey": null
      }
    ],
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "View1Query",
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Cluster",
        "kind": "LinkedField",
        "name": "cluster",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": (v1/*: any*/),
            "concreteType": "ClusterNode",
            "kind": "LinkedField",
            "name": "clusterNodes",
            "plural": true,
            "selections": [
              {
                "alias": null,
                "args": null,
                "concreteType": "Pod",
                "kind": "LinkedField",
                "name": "pods",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  {
                    "condition": "showPods",
                    "kind": "Condition",
                    "passingValue": true,
                    "selections": [
                      (v3/*: any*/)
                    ]
                  }
                ],
                "storageKey": null
              },
              (v2/*: any*/)
            ],
            "storageKey": null
          },
          (v2/*: any*/)
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "3a1759e510fd531aef6b5f9858a0e098",
    "id": null,
    "metadata": {},
    "name": "View1Query",
    "operationKind": "query",
    "text": "query View1Query(\n  $nodeID: ID\n  $showClusters: Boolean!\n  $showNodes: Boolean!\n  $showPods: Boolean!\n) {\n  cluster {\n    id @include(if: $showClusters)\n    clusterNodes(nodeid: $nodeID) {\n      id @include(if: $showNodes)\n      pods {\n        id @include(if: $showPods)\n        name @include(if: $showPods)\n        id\n      }\n      id\n    }\n    id\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'c576161bae31d2013820a9356b21defe';

module.exports = node;

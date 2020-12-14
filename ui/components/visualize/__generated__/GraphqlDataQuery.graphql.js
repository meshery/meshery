/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type GraphqlDataQueryVariables = {|
  nodeID?: ?string
|};
export type GraphqlDataQueryResponse = {|
  +cluster: ?$ReadOnlyArray<?{|
    +clusterNodes: ?$ReadOnlyArray<?{|
      +id: string,
      +parentid: string,
      +pods: $ReadOnlyArray<?{|
        +id: string,
        +parentid: string,
        +name: string,
      |}>,
    |}>
  |}>
|};
export type GraphqlDataQuery = {|
  variables: GraphqlDataQueryVariables,
  response: GraphqlDataQueryResponse,
|};
*/


/*
query GraphqlDataQuery(
  $nodeID: ID
) {
  cluster {
    clusterNodes(nodeid: $nodeID) {
      id
      parentid
      pods {
        id
        parentid
        name
      }
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
  }
],
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v2 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "parentid",
  "storageKey": null
},
v3 = {
  "alias": null,
  "args": [
    {
      "kind": "Variable",
      "name": "nodeid",
      "variableName": "nodeID"
    }
  ],
  "concreteType": "ClusterNode",
  "kind": "LinkedField",
  "name": "clusterNodes",
  "plural": true,
  "selections": [
    (v1/*: any*/),
    (v2/*: any*/),
    {
      "alias": null,
      "args": null,
      "concreteType": "Pod",
      "kind": "LinkedField",
      "name": "pods",
      "plural": true,
      "selections": [
        (v1/*: any*/),
        (v2/*: any*/),
        {
          "alias": null,
          "args": null,
          "kind": "ScalarField",
          "name": "name",
          "storageKey": null
        }
      ],
      "storageKey": null
    }
  ],
  "storageKey": null
};
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "GraphqlDataQuery",
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Cluster",
        "kind": "LinkedField",
        "name": "cluster",
        "plural": true,
        "selections": [
          (v3/*: any*/)
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
    "name": "GraphqlDataQuery",
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Cluster",
        "kind": "LinkedField",
        "name": "cluster",
        "plural": true,
        "selections": [
          (v3/*: any*/),
          (v1/*: any*/)
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "3be50f84716323401e3831d42b13fe80",
    "id": null,
    "metadata": {},
    "name": "GraphqlDataQuery",
    "operationKind": "query",
    "text": "query GraphqlDataQuery(\n  $nodeID: ID\n) {\n  cluster {\n    clusterNodes(nodeid: $nodeID) {\n      id\n      parentid\n      pods {\n        id\n        parentid\n        name\n      }\n    }\n    id\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'cb2e0b12aaa20f7d04b49f32832d4a1a';

module.exports = node;

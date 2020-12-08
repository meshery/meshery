/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type GraphqlDataQueryVariables = {||};
export type GraphqlDataQueryResponse = {|
  +cluster: ?$ReadOnlyArray<?{|
    +clusterNodes: ?$ReadOnlyArray<?{|
      +id: string,
      +parentid: string,
    |}>
  |}>
|};
export type GraphqlDataQuery = {|
  variables: GraphqlDataQueryVariables,
  response: GraphqlDataQueryResponse,
|};
*/


/*
query GraphqlDataQuery {
  cluster {
    clusterNodes {
      id
      parentid
    }
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v1 = {
  "alias": null,
  "args": null,
  "concreteType": "ClusterNode",
  "kind": "LinkedField",
  "name": "clusterNodes",
  "plural": true,
  "selections": [
    (v0/*: any*/),
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "parentid",
      "storageKey": null
    }
  ],
  "storageKey": null
};
return {
  "fragment": {
    "argumentDefinitions": [],
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
          (v1/*: any*/)
        ],
        "storageKey": null
      }
    ],
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
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
          (v1/*: any*/),
          (v0/*: any*/)
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "bc178dded9175a88d33145f3a7137b69",
    "id": null,
    "metadata": {},
    "name": "GraphqlDataQuery",
    "operationKind": "query",
    "text": "query GraphqlDataQuery {\n  cluster {\n    clusterNodes {\n      id\n      parentid\n    }\n    id\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd5f0f8ccb3f9b8592ae2ffe56b746098';

module.exports = node;

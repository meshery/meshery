/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type MeshsyncStatusQueryVariables = {||};
export type MeshsyncStatusQueryResponse = {|
  +controller: {|
    +name: string,
    +version: string,
    +status: Status,
  |}
|};
export type MeshsyncStatusQuery = {|
  variables: MeshsyncStatusQueryVariables,
  response: MeshsyncStatusQueryResponse,
|};
*/


/*
query MeshsyncStatusQuery {
  controller: getMeshsyncStatus {
    name
    version
    status
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": "controller",
    "args": null,
    "concreteType": "OperatorControllerStatus",
    "kind": "LinkedField",
    "name": "getMeshsyncStatus",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "name",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "version",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "status",
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "MeshsyncStatusQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "MeshsyncStatusQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "d3b1a1a67c5b3f740e169a86665e1751",
    "id": null,
    "metadata": {},
    "name": "MeshsyncStatusQuery",
    "operationKind": "query",
    "text": "query MeshsyncStatusQuery {\n  controller: getMeshsyncStatus {\n    name\n    version\n    status\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'c009f646493c799c1dd4d6fd9e996ac8';

module.exports = node;

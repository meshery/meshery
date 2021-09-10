/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type DeployMeshSyncQueryVariables = {||};
export type DeployMeshSyncQueryResponse = {|
  +deployMeshsync: Status
|};
export type DeployMeshSyncQuery = {|
  variables: DeployMeshSyncQueryVariables,
  response: DeployMeshSyncQueryResponse,
|};
*/


/*
query DeployMeshSyncQuery {
  deployMeshsync
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": null,
    "args": null,
    "kind": "ScalarField",
    "name": "deployMeshsync",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "DeployMeshSyncQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "DeployMeshSyncQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "4ccc7a359f435f939905d3db4fef6e62",
    "id": null,
    "metadata": {},
    "name": "DeployMeshSyncQuery",
    "operationKind": "query",
    "text": "query DeployMeshSyncQuery {\n  deployMeshsync\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '0f9c347cda862c0d3aa97933d85e5814';

module.exports = node;

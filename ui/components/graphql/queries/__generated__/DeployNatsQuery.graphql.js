/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type DeployNatsQueryVariables = {||};
export type DeployNatsQueryResponse = {|
  +connectToNats: Status
|};
export type DeployNatsQuery = {|
  variables: DeployNatsQueryVariables,
  response: DeployNatsQueryResponse,
|};
*/


/*
query DeployNatsQuery {
  connectToNats
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": null,
    "args": null,
    "kind": "ScalarField",
    "name": "connectToNats",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "DeployNatsQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "DeployNatsQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "52b1fa366093c490c90977eddb30b357",
    "id": null,
    "metadata": {},
    "name": "DeployNatsQuery",
    "operationKind": "query",
    "text": "query DeployNatsQuery {\n  connectToNats\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '518fa6beed8e9831105a05d05e904948';

module.exports = node;

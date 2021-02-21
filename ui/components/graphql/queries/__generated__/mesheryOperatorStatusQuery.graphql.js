/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "UNKNOWN" | "%future added value";
export type mesheryOperatorStatusQueryVariables = {||};
export type mesheryOperatorStatusQueryResponse = {|
  +getOperatorStatus: ?Status
|};
export type mesheryOperatorStatusQuery = {|
  variables: mesheryOperatorStatusQueryVariables,
  response: mesheryOperatorStatusQueryResponse,
|};
*/


/*
query mesheryOperatorStatusQuery {
  getOperatorStatus
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": null,
    "args": null,
    "kind": "ScalarField",
    "name": "getOperatorStatus",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "mesheryOperatorStatusQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "mesheryOperatorStatusQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "9cfc892c596379577a122481f2c28f28",
    "id": null,
    "metadata": {},
    "name": "mesheryOperatorStatusQuery",
    "operationKind": "query",
    "text": "query mesheryOperatorStatusQuery {\n  getOperatorStatus\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'b3a18d8178377756fad98f2f4c28c64f';

module.exports = node;

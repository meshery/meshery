/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorStatusMutationVariables = {|
  targetStatus: Status
|};
export type OperatorStatusMutationResponse = {|
  +operatorStatus: Status
|};
export type OperatorStatusMutation = {|
  variables: OperatorStatusMutationVariables,
  response: OperatorStatusMutationResponse,
|};
*/


/*
mutation OperatorStatusMutation(
  $targetStatus: Status!
) {
  operatorStatus: changeOperatorStatus(targetStatus: $targetStatus)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "targetStatus"
  }
],
v1 = [
  {
    "alias": "operatorStatus",
    "args": [
      {
        "kind": "Variable",
        "name": "targetStatus",
        "variableName": "targetStatus"
      }
    ],
    "kind": "ScalarField",
    "name": "changeOperatorStatus",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "OperatorStatusMutation",
    "selections": (v1/*: any*/),
    "type": "Mutation",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "OperatorStatusMutation",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "f0f81d6685889dd15748788bd1e2d206",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusMutation",
    "operationKind": "mutation",
    "text": "mutation OperatorStatusMutation(\n  $targetStatus: Status!\n) {\n  operatorStatus: changeOperatorStatus(targetStatus: $targetStatus)\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f18b68bf4a817f8c1752cabe5235bf1f';

module.exports = node;

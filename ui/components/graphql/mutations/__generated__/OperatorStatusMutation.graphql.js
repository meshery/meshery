/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorStatusInput = {|
  targetStatus: Status
|};
export type OperatorStatusMutationVariables = {|
  input?: ?OperatorStatusInput
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
  $input: OperatorStatusInput
) {
  operatorStatus: changeOperatorStatus(input: $input)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "input"
  }
],
v1 = [
  {
    "alias": "operatorStatus",
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
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
    "cacheID": "86d68e3b96cd8684338daf88a6f49ab5",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusMutation",
    "operationKind": "mutation",
    "text": "mutation OperatorStatusMutation(\n  $input: OperatorStatusInput\n) {\n  operatorStatus: changeOperatorStatus(input: $input)\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'fc306fe156aa1a0f9984281666bd7693';

module.exports = node;

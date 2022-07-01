/**
 * @generated SignedSource<<3f8e706210e59168563a58ae4724c63c>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Mutation } from 'relay-runtime';
export type Status = "ENABLED" | "CONNECTED" | "DISABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorStatusInput = {|
  targetStatus: Status,
  contextID: string,
|};
export type OperatorStatusMutation$variables = {|
  input?: ?OperatorStatusInput,
|};
export type OperatorStatusMutation$data = {|
  +operatorStatus: Status,
|};
export type OperatorStatusMutation = {|
  variables: OperatorStatusMutation$variables,
  response: OperatorStatusMutation$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
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

(node/*: any*/).hash = "fc306fe156aa1a0f9984281666bd7693";

module.exports = ((node/*: any*/)/*: Mutation<
  OperatorStatusMutation$variables,
  OperatorStatusMutation$data,
>*/);

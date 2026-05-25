/**
 * @generated SignedSource<<0dd9419c0a7558919b38722c00dfcffc>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorStatusInput = {
  contextID: string;
  targetStatus: Status;
};
export type OperatorStatusMutation$variables = {
  input?: OperatorStatusInput | null | undefined;
};
export type OperatorStatusMutation$data = {
  readonly operatorStatus: Status;
};
export type OperatorStatusMutation = {
  response: OperatorStatusMutation$data;
  variables: OperatorStatusMutation$variables;
};

const node: ConcreteRequest = (function(){
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

(node as any).hash = "fc306fe156aa1a0f9984281666bd7693";

export default node;

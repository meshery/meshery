/**
 * @generated SignedSource<<ed1ec0402fa8b4861339501488a3e28e>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type AdapterStatusInput = {
  adapter: string;
  targetPort: string;
  targetStatus: Status;
};
export type AdapterStatusMutation$variables = {
  input?: AdapterStatusInput | null | undefined;
};
export type AdapterStatusMutation$data = {
  readonly adapterStatus: Status;
};
export type AdapterStatusMutation = {
  response: AdapterStatusMutation$data;
  variables: AdapterStatusMutation$variables;
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
    "alias": "adapterStatus",
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "kind": "ScalarField",
    "name": "changeAdapterStatus",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "AdapterStatusMutation",
    "selections": (v1/*: any*/),
    "type": "Mutation",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AdapterStatusMutation",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "910bac5b935be9f44edf2e0b45186c89",
    "id": null,
    "metadata": {},
    "name": "AdapterStatusMutation",
    "operationKind": "mutation",
    "text": "mutation AdapterStatusMutation(\n  $input: AdapterStatusInput\n) {\n  adapterStatus: changeAdapterStatus(input: $input)\n}\n"
  }
};
})();

(node as any).hash = "3d97311fb835925f575991945093af8e";

export default node;

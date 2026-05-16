/**
 * @generated SignedSource<<af9b90fa300571b33bcd5c7f70347fcd>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type MesheryController = "BROKER" | "MESHSYNC" | "OPERATOR" | "%future added value";
export type MesheryControllerStatus = "CONNECTED" | "DEPLOYED" | "DEPLOYING" | "ENABLED" | "NOTDEPLOYED" | "RUNNING" | "UNDEPLOYED" | "UNKOWN" | "%future added value";
export type OperatorStatusQuery$variables = {
  connectionID: string;
};
export type OperatorStatusQuery$data = {
  readonly operator: {
    readonly connectionID: string;
    readonly controller: MesheryController;
    readonly status: MesheryControllerStatus;
  } | null | undefined;
};
export type OperatorStatusQuery = {
  response: OperatorStatusQuery$data;
  variables: OperatorStatusQuery$variables;
};

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "connectionID"
  }
],
v1 = [
  {
    "alias": "operator",
    "args": [
      {
        "kind": "Variable",
        "name": "connectionID",
        "variableName": "connectionID"
      }
    ],
    "concreteType": "MesheryControllersStatusListItem",
    "kind": "LinkedField",
    "name": "getOperatorStatus",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "status",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "controller",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "connectionID",
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "OperatorStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "OperatorStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "1edf459162cb7e0a916856a8bb8649e2",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusQuery",
    "operationKind": "query",
    "text": "query OperatorStatusQuery(\n  $connectionID: String!\n) {\n  operator: getOperatorStatus(connectionID: $connectionID) {\n    status\n    controller\n    connectionID\n  }\n}\n"
  }
};
})();

(node as any).hash = "519b77a47c22930c9d148d1e4f4011cb";

export default node;

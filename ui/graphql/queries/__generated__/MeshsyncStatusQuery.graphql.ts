/**
 * @generated SignedSource<<cd252b762f4c91a7b1ee54846a2eb9fd>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type MeshsyncStatusQuery$variables = {
  connectionID: string;
};
export type MeshsyncStatusQuery$data = {
  readonly controller: {
    readonly name: string;
    readonly status: Status;
    readonly version: string;
  };
};
export type MeshsyncStatusQuery = {
  response: MeshsyncStatusQuery$data;
  variables: MeshsyncStatusQuery$variables;
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
    "alias": "controller",
    "args": [
      {
        "kind": "Variable",
        "name": "connectionID",
        "variableName": "connectionID"
      }
    ],
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
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "MeshsyncStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshsyncStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "1b2cb7cffaec7ee51423e6f26c1e6ae4",
    "id": null,
    "metadata": {},
    "name": "MeshsyncStatusQuery",
    "operationKind": "query",
    "text": "query MeshsyncStatusQuery(\n  $connectionID: String!\n) {\n  controller: getMeshsyncStatus(connectionID: $connectionID) {\n    name\n    version\n    status\n  }\n}\n"
  }
};
})();

(node as any).hash = "a9ff2d13ffaf332f9be4ca12e09bd7f9";

export default node;

/**
 * @generated SignedSource<<ab8619d4497e08066a968ab897a59fa1>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type NatsStatusQuery$variables = {
  connectionID: string;
};
export type NatsStatusQuery$data = {
  readonly controller: {
    readonly name: string;
    readonly status: Status;
    readonly version: string;
  };
};
export type NatsStatusQuery = {
  response: NatsStatusQuery$data;
  variables: NatsStatusQuery$variables;
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
    "name": "getNatsStatus",
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
    "name": "NatsStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "NatsStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "f232d0001a24e998f5d993786cfed5f6",
    "id": null,
    "metadata": {},
    "name": "NatsStatusQuery",
    "operationKind": "query",
    "text": "query NatsStatusQuery(\n  $connectionID: String!\n) {\n  controller: getNatsStatus(connectionID: $connectionID) {\n    name\n    version\n    status\n  }\n}\n"
  }
};
})();

(node as any).hash = "10628d273554398f8d127aaa764a94fd";

export default node;

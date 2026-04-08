/**
 * @generated SignedSource<<5cb2f0c3ecc377b812689c205dbf2f97>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type TelemetryComponentsQuery$variables = {
  contexts?: ReadonlyArray<string> | null | undefined;
};
export type TelemetryComponentsQuery$data = {
  readonly telemetryComps: ReadonlyArray<{
    readonly name: string;
    readonly spec: string;
    readonly status: string;
  } | null | undefined>;
};
export type TelemetryComponentsQuery = {
  response: TelemetryComponentsQuery$data;
  variables: TelemetryComponentsQuery$variables;
};

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "contexts"
  }
],
v1 = [
  {
    "alias": "telemetryComps",
    "args": [
      {
        "kind": "Variable",
        "name": "contexts",
        "variableName": "contexts"
      }
    ],
    "concreteType": "TelemetryComp",
    "kind": "LinkedField",
    "name": "fetchTelemetryComponents",
    "plural": true,
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
        "name": "spec",
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
    "name": "TelemetryComponentsQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "TelemetryComponentsQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "5c359a9d5b083eb5aba9b7bfe49663c9",
    "id": null,
    "metadata": {},
    "name": "TelemetryComponentsQuery",
    "operationKind": "query",
    "text": "query TelemetryComponentsQuery(\n  $contexts: [String!]\n) {\n  telemetryComps: fetchTelemetryComponents(contexts: $contexts) {\n    name\n    spec\n    status\n  }\n}\n"
  }
};
})();

(node as any).hash = "ec7907a52eb8925b7e063038d08beb85";

export default node;

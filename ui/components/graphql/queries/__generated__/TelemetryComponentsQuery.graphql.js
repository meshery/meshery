/**
 * @generated SignedSource<<56c44d77d934736ff49c3bb96b344e0f>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

var node = (function(){
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

node.hash = "ec7907a52eb8925b7e063038d08beb85";

export default node;

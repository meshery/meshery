/**
 * @generated SignedSource<<64a6e2242e655d185ae16b2fb0efa9fb>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type Severity = "alert" | "critical" | "debug" | "emergency" | "error" | "informational" | "warning" | "%future added value";
export type EventsSubscription$variables = Record<PropertyKey, never>;
export type EventsSubscription$data = {
  readonly event: {
    readonly actedUpon: string;
    readonly action: string;
    readonly category: string;
    readonly createdAt: any;
    readonly deletedAt: any | null | undefined;
    readonly description: string;
    readonly id: string;
    readonly metadata: any | null | undefined;
    readonly operationID: string;
    readonly severity: Severity;
    readonly status: string;
    readonly systemID: string;
    readonly updatedAt: any;
    readonly userID: string;
  };
};
export type EventsSubscription = {
  response: EventsSubscription$data;
  variables: EventsSubscription$variables;
};

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "alias": "event",
    "args": null,
    "concreteType": "Event",
    "kind": "LinkedField",
    "name": "subscribeEvents",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "id",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "userID",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "actedUpon",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "operationID",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "systemID",
        "storageKey": null
      },
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
        "name": "severity",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "action",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "category",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "description",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "metadata",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "createdAt",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "updatedAt",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "deletedAt",
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "EventsSubscription",
    "selections": (v0/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "EventsSubscription",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "dfb0d28db01886a2b8877f52db0fd0c6",
    "id": null,
    "metadata": {},
    "name": "EventsSubscription",
    "operationKind": "subscription",
    "text": "subscription EventsSubscription {\n  event: subscribeEvents {\n    id\n    userID\n    actedUpon\n    operationID\n    systemID\n    status\n    severity\n    action\n    category\n    description\n    metadata\n    createdAt\n    updatedAt\n    deletedAt\n  }\n}\n"
  }
};
})();

(node as any).hash = "b6a9380e1fcb5b6b417ef482bb98ab7b";

export default node;

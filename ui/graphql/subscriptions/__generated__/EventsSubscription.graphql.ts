/**
 * @generated SignedSource<<4156a3806335259d89f78a0ecba53d0b>>
 * @lightSyntaxTransform
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
    readonly owner: string;
    readonly severity: Severity;
    readonly status: string;
    readonly systemID: string;
    readonly updatedAt: any;
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
        "name": "owner",
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
    "selections": (v0/*:: as any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "EventsSubscription",
    "selections": (v0/*:: as any*/)
  },
  "params": {
    "cacheID": "1dd0f3928f73d92756772fd4c7093c31",
    "id": null,
    "metadata": {},
    "name": "EventsSubscription",
    "operationKind": "subscription",
    "text": "subscription EventsSubscription {\n  event: subscribeEvents {\n    id\n    owner\n    actedUpon\n    operationID\n    systemID\n    status\n    severity\n    action\n    category\n    description\n    metadata\n    createdAt\n    updatedAt\n    deletedAt\n  }\n}\n"
  }
};
})();

(node as any).hash = "a46175bd2d4f20226c798fe5c98548fc";

export default node;

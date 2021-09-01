/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorStatusEventsSubscriptionVariables = {||};
export type OperatorStatusEventsSubscriptionResponse = {|
  +operator: {|
    +status: Status,
    +version: string,
    +controllers: $ReadOnlyArray<{|
      +name: string,
      +version: string,
      +status: Status,
    |}>,
    +error: ?{|
      +code: string,
      +description: string,
    |},
  |}
|};
export type OperatorStatusEventsSubscription = {|
  variables: OperatorStatusEventsSubscriptionVariables,
  response: OperatorStatusEventsSubscriptionResponse,
|};
*/


/*
subscription OperatorStatusEventsSubscription {
  operator: listenToOperatorState {
    status
    version
    controllers {
      name
      version
      status
    }
    error {
      code
      description
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "status",
  "storageKey": null
},
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "version",
  "storageKey": null
},
v2 = [
  {
    "alias": "operator",
    "args": null,
    "concreteType": "OperatorStatus",
    "kind": "LinkedField",
    "name": "listenToOperatorState",
    "plural": false,
    "selections": [
      (v0/*: any*/),
      (v1/*: any*/),
      {
        "alias": null,
        "args": null,
        "concreteType": "OperatorControllerStatus",
        "kind": "LinkedField",
        "name": "controllers",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "name",
            "storageKey": null
          },
          (v1/*: any*/),
          (v0/*: any*/)
        ],
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "Error",
        "kind": "LinkedField",
        "name": "error",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "code",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "description",
            "storageKey": null
          }
        ],
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
    "name": "OperatorStatusEventsSubscription",
    "selections": (v2/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "OperatorStatusEventsSubscription",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "d185a04533119e60576876c0db6f85bf",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusEventsSubscription",
    "operationKind": "subscription",
    "text": "subscription OperatorStatusEventsSubscription {\n  operator: listenToOperatorState {\n    status\n    version\n    controllers {\n      name\n      version\n      status\n    }\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '1bd441d8731a86fa8508fab854a70c1a';

module.exports = node;

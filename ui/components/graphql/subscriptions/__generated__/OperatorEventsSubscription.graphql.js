/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorEventsSubscriptionVariables = {||};
export type OperatorEventsSubscriptionResponse = {|
  +listenToOperatorState: {|
    +status: ?Status,
    +error: ?{|
      +code: string,
      +description: string,
    |},
  |}
|};
export type OperatorEventsSubscription = {|
  variables: OperatorEventsSubscriptionVariables,
  response: OperatorEventsSubscriptionResponse,
|};
*/


/*
subscription OperatorEventsSubscription {
  listenToOperatorState {
    status
    error {
      code
      description
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": null,
    "args": null,
    "concreteType": "OperatorStatus",
    "kind": "LinkedField",
    "name": "listenToOperatorState",
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
    "name": "OperatorEventsSubscription",
    "selections": (v0/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "OperatorEventsSubscription",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "c90f4010bcc9c128c7d827d834e52bd6",
    "id": null,
    "metadata": {},
    "name": "OperatorEventsSubscription",
    "operationKind": "subscription",
    "text": "subscription OperatorEventsSubscription {\n  listenToOperatorState {\n    status\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a223a9f1efab40e805723e9e4ad5a6ce';

module.exports = node;

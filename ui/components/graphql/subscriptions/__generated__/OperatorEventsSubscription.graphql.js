/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "UNKNOWN" | "%future added value";
export type OperatorEventsSubscriptionVariables = {||};
export type OperatorEventsSubscriptionResponse = {|
  +listenToOperatorEvents: {|
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
  listenToOperatorEvents {
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
    "name": "listenToOperatorEvents",
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
    "cacheID": "2af0544985213b071fb9ad4f4506c899",
    "id": null,
    "metadata": {},
    "name": "OperatorEventsSubscription",
    "operationKind": "subscription",
    "text": "subscription OperatorEventsSubscription {\n  listenToOperatorEvents {\n    status\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '8e686bdd1aab4653b7d89df935c29d4c';

module.exports = node;

/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorStatusSubscriptionVariables = {||};
export type OperatorStatusSubscriptionResponse = {|
  +operator: {|
    +status: Status,
    +error: ?{|
      +code: string,
      +description: string,
    |},
  |}
|};
export type OperatorStatusSubscription = {|
  variables: OperatorStatusSubscriptionVariables,
  response: OperatorStatusSubscriptionResponse,
|};
*/


/*
subscription OperatorStatusSubscription {
  operator: listenToOperatorState {
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
    "alias": "operator",
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
    "name": "OperatorStatusSubscription",
    "selections": (v0/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "OperatorStatusSubscription",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "ee07f84c0887b3172099db3261faa84e",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription OperatorStatusSubscription {\n  operator: listenToOperatorState {\n    status\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '1946f24deab962396e84e4a0b9f1ec86';

module.exports = node;

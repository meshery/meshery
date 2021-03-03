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
  +listenToOperatorState: {|
    +status: ?Status,
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
    "cacheID": "408df65d57533618b763e13bbc6d0b2c",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription OperatorStatusSubscription {\n  listenToOperatorState {\n    status\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '495cc46f2653c9391ef3b7bb38457b1f';

module.exports = node;

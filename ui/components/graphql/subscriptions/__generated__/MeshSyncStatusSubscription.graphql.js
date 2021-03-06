/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type MeshSyncStatusSubscriptionVariables = {||};
export type MeshSyncStatusSubscriptionResponse = {|
  +listenToMeshSyncEvents: {|
    +name: ?string,
    +status: ?Status,
    +error: ?{|
      +code: string,
      +description: string,
    |},
  |}
|};
export type MeshSyncStatusSubscription = {|
  variables: MeshSyncStatusSubscriptionVariables,
  response: MeshSyncStatusSubscriptionResponse,
|};
*/


/*
subscription MeshSyncStatusSubscription {
  listenToMeshSyncEvents {
    name
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
    "concreteType": "OperatorControllerStatus",
    "kind": "LinkedField",
    "name": "listenToMeshSyncEvents",
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
    "name": "MeshSyncStatusSubscription",
    "selections": (v0/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "MeshSyncStatusSubscription",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "e330914ffe44d508e311f3dd7d625aa9",
    "id": null,
    "metadata": {},
    "name": "MeshSyncStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription MeshSyncStatusSubscription {\n  listenToMeshSyncEvents {\n    name\n    status\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'e0210af501e9bd9864ea6444ea6a5db2';

module.exports = node;

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
    +name: string,
    +status: Status,
    +version: string,
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
    version
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
        "kind": "ScalarField",
        "name": "version",
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
    "cacheID": "9a3d3b03529103cf9d737147ed7c2d7a",
    "id": null,
    "metadata": {},
    "name": "MeshSyncStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription MeshSyncStatusSubscription {\n  listenToMeshSyncEvents {\n    name\n    status\n    version\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f9a9715d6c56d418d17154fba8bda254';

module.exports = node;

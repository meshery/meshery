/**
 * @generated SignedSource<<451b18fba5c3680f231918922efe6cf9>>
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
    "name": "connectionIDs"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "eventTypes"
  }
],
v1 = [
  {
    "alias": "meshsyncevents",
    "args": [
      {
        "kind": "Variable",
        "name": "connectionIDs",
        "variableName": "connectionIDs"
      },
      {
        "kind": "Variable",
        "name": "eventTypes",
        "variableName": "eventTypes"
      }
    ],
    "concreteType": "MeshSyncEvent",
    "kind": "LinkedField",
    "name": "subscribeMeshSyncEvents",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "type",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "object",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "connectionID",
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
    "name": "MeshSyncEventsSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshSyncEventsSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "29cb26069b69c2858bb31abe76540d9f",
    "id": null,
    "metadata": {},
    "name": "MeshSyncEventsSubscription",
    "operationKind": "subscription",
    "text": "subscription MeshSyncEventsSubscription(\n  $connectionIDs: [String!]\n  $eventTypes: [MeshSyncEventType!]\n) {\n  meshsyncevents: subscribeMeshSyncEvents(connectionIDs: $connectionIDs, eventTypes: $eventTypes) {\n    type\n    object\n    connectionID\n  }\n}\n"
  }
};
})();

node.hash = "422c02f316080b9b8d43a562f696a3d4";

module.exports = node;

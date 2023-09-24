/**
 * @generated SignedSource<<b5265e02bc691cfba41a29225518a02f>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

var node = (function(){
var v0 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "eventTypes"
},
v1 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "k8scontextIDs"
},
v2 = [
  {
    "alias": "meshsyncevents",
    "args": [
      {
        "kind": "Variable",
        "name": "eventTypes",
        "variableName": "eventTypes"
      },
      {
        "kind": "Variable",
        "name": "k8scontextIDs",
        "variableName": "k8scontextIDs"
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
        "name": "contextId",
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": [
      (v0/*: any*/),
      (v1/*: any*/)
    ],
    "kind": "Fragment",
    "metadata": null,
    "name": "MeshSyncEventsSubscription",
    "selections": (v2/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v1/*: any*/),
      (v0/*: any*/)
    ],
    "kind": "Operation",
    "name": "MeshSyncEventsSubscription",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "587f7af94d634bb28ea91b1790e28e1f",
    "id": null,
    "metadata": {},
    "name": "MeshSyncEventsSubscription",
    "operationKind": "subscription",
    "text": "subscription MeshSyncEventsSubscription(\n  $k8scontextIDs: [String!]\n  $eventTypes: [MeshSyncEventType!]\n) {\n  meshsyncevents: subscribeMeshSyncEvents(k8scontextIDs: $k8scontextIDs, eventTypes: $eventTypes) {\n    type\n    object\n    contextId\n  }\n}\n"
  }
};
})();

node.hash = "8d4a467a34c136646e2af2ec3107a204";

module.exports = node;

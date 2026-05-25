/**
 * @generated SignedSource<<a76443003a12f66769cc0bb77e59630b>>
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
    "name": "filter"
  }
],
v1 = [
  {
    "alias": "addonsState",
    "args": [
      {
        "kind": "Variable",
        "name": "filter",
        "variableName": "filter"
      }
    ],
    "concreteType": "AddonList",
    "kind": "LinkedField",
    "name": "listenToAddonState",
    "plural": true,
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
        "name": "owner",
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
    "name": "AddonStatusSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AddonStatusSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "b7f42d3b9c1ab13f1d78c3048eb9a5db",
    "id": null,
    "metadata": {},
    "name": "AddonStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription AddonStatusSubscription(\n  $filter: ServiceMeshFilter\n) {\n  addonsState: listenToAddonState(filter: $filter) {\n    name\n    owner\n  }\n}\n"
  }
};
})();

node.hash = "7cef73e9bfdcc63d12dfe54d2a0f3fbf";

module.exports = node;

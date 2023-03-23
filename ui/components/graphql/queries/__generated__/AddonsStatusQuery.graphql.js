/**
 * @generated SignedSource<<1282ae0224439164af34f0d4fba4a71a>>
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
    "name": "getAvailableAddons",
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
    "name": "AddonsStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AddonsStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "49f96950f100465ccdab36e9903b9281",
    "id": null,
    "metadata": {},
    "name": "AddonsStatusQuery",
    "operationKind": "query",
    "text": "query AddonsStatusQuery(\n  $filter: ServiceMeshFilter\n) {\n  addonsState: getAvailableAddons(filter: $filter) {\n    name\n    owner\n  }\n}\n"
  }
};
})();

node.hash = "9cbf0a827a321dead7e3e6d0c2e9cbe7";

module.exports = node;

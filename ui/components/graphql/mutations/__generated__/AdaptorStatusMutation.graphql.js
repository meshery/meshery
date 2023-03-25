/**
 * @generated SignedSource<<bf8a810fe00984ff88e24dc09b2d6e9b>>
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
    "name": "input"
  }
],
v1 = [
  {
    "alias": "adaptorStatus",
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "kind": "ScalarField",
    "name": "changeAdaptorStatus",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "AdaptorStatusMutation",
    "selections": (v1/*: any*/),
    "type": "Mutation",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AdaptorStatusMutation",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "7120ac79bbdc5508fcdabea89a855ba5",
    "id": null,
    "metadata": {},
    "name": "AdaptorStatusMutation",
    "operationKind": "mutation",
    "text": "mutation AdaptorStatusMutation(\n  $input: AdaptorStatusInput\n) {\n  adaptorStatus: changeAdaptorStatus(input: $input)\n}\n"
  }
};
})();

node.hash = "794bade1e18539e87d0b561d2484bb6b";

module.exports = node;

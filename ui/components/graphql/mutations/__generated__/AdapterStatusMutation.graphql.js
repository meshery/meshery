/**
 * @generated SignedSource<<aff6f55dc556ea8dcf170a9dd4d5853b>>
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
    "alias": "adapterStatus",
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "kind": "ScalarField",
    "name": "changeAdapterStatus",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "AdapterStatusMutation",
    "selections": (v1/*: any*/),
    "type": "Mutation",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AdapterStatusMutation",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "910bac5b935be9f44edf2e0b45186c89",
    "id": null,
    "metadata": {},
    "name": "AdapterStatusMutation",
    "operationKind": "mutation",
    "text": "mutation AdapterStatusMutation(\n  $input: AdapterStatusInput\n) {\n  adapterStatus: changeAdapterStatus(input: $input)\n}\n"
  }
};
})();

node.hash = "3d97311fb835925f575991945093af8e";

module.exports = node;

/**
 * @generated SignedSource<<a4c6a02f54fde7994dda88f5a554953f>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type Status = "ENABLED" | "CONNECTED" | "DISABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type DeployNatsQuery$variables = {|
  k8scontextID: string,
|};
export type DeployNatsQuery$data = {|
  +connectToNats: Status,
|};
export type DeployNatsQuery = {|
  variables: DeployNatsQuery$variables,
  response: DeployNatsQuery$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "k8scontextID"
  }
],
v1 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextID",
        "variableName": "k8scontextID"
      }
    ],
    "kind": "ScalarField",
    "name": "connectToNats",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "DeployNatsQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "DeployNatsQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "76180bca9f0c5921831b5c004d7e2679",
    "id": null,
    "metadata": {},
    "name": "DeployNatsQuery",
    "operationKind": "query",
    "text": "query DeployNatsQuery(\n  $k8scontextID: String!\n) {\n  connectToNats(k8scontextID: $k8scontextID)\n}\n"
  }
};
})();

(node/*: any*/).hash = "36be030e0bed096e2cfaef1101bfaa27";

module.exports = ((node/*: any*/)/*: Query<
  DeployNatsQuery$variables,
  DeployNatsQuery$data,
>*/);

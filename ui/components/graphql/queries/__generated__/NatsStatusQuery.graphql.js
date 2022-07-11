/**
 * @generated SignedSource<<b88894883739e2e5e0a25ebe270ad416>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type Status = "ENABLED" | "CONNECTED" | "DISABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type NatsStatusQuery$variables = {|
  k8scontextID: string,
|};
export type NatsStatusQuery$data = {|
  +controller: {|
    +name: string,
    +version: string,
    +status: Status,
  |},
|};
export type NatsStatusQuery = {|
  variables: NatsStatusQuery$variables,
  response: NatsStatusQuery$data,
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
    "alias": "controller",
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextID",
        "variableName": "k8scontextID"
      }
    ],
    "concreteType": "OperatorControllerStatus",
    "kind": "LinkedField",
    "name": "getNatsStatus",
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
        "name": "version",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "status",
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
    "name": "NatsStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "NatsStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "e02d86123c1155db14fdabb73d68a596",
    "id": null,
    "metadata": {},
    "name": "NatsStatusQuery",
    "operationKind": "query",
    "text": "query NatsStatusQuery(\n  $k8scontextID: String!\n) {\n  controller: getNatsStatus(k8scontextID: $k8scontextID) {\n    name\n    version\n    status\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "4dc5b2d81433e73a9eab8bd581d16b46";

module.exports = ((node/*: any*/)/*: Query<
  NatsStatusQuery$variables,
  NatsStatusQuery$data,
>*/);

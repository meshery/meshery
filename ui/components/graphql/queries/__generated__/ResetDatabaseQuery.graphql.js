/**
 * @generated SignedSource<<fb4f2e789f74ac2e7d93af7785c3e673>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type Status = "ENABLED" | "CONNECTED" | "DISABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type ReSyncActions = {|
  clearDB: string,
  ReSync: string,
  hardReset: string,
|};
export type ResetDatabaseQuery$variables = {|
  selector: ReSyncActions,
  k8scontextID: string,
|};
export type ResetDatabaseQuery$data = {|
  +resetStatus: Status,
|};
export type ResetDatabaseQuery = {|
  variables: ResetDatabaseQuery$variables,
  response: ResetDatabaseQuery$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "k8scontextID"
},
v1 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "selector"
},
v2 = [
  {
    "alias": "resetStatus",
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextID",
        "variableName": "k8scontextID"
      },
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "kind": "ScalarField",
    "name": "resyncCluster",
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
    "name": "ResetDatabaseQuery",
    "selections": (v2/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v1/*: any*/),
      (v0/*: any*/)
    ],
    "kind": "Operation",
    "name": "ResetDatabaseQuery",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "37b0c2d517499c337bd4bdfa2ef79380",
    "id": null,
    "metadata": {},
    "name": "ResetDatabaseQuery",
    "operationKind": "query",
    "text": "query ResetDatabaseQuery(\n  $selector: ReSyncActions!\n  $k8scontextID: String!\n) {\n  resetStatus: resyncCluster(selector: $selector, k8scontextID: $k8scontextID)\n}\n"
  }
};
})();

(node/*: any*/).hash = "54a9344cc4d95023f5082936dc95d05d";

module.exports = ((node/*: any*/)/*: Query<
  ResetDatabaseQuery$variables,
  ResetDatabaseQuery$data,
>*/);

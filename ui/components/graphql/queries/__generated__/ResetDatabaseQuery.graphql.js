/**
 * @generated SignedSource<<b21d77de86999973940b6b6c9359aff2>>
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
|};
export type ResetDatabaseQuery$variables = {|
  selector: ReSyncActions,
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
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "selector"
  }
],
v1 = [
  {
    "alias": "resetStatus",
    "args": [
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
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "ResetDatabaseQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ResetDatabaseQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "5c0d3c41070a436ea5d9485d62823c5c",
    "id": null,
    "metadata": {},
    "name": "ResetDatabaseQuery",
    "operationKind": "query",
    "text": "query ResetDatabaseQuery(\n  $selector: ReSyncActions!\n) {\n  resetStatus: resyncCluster(selector: $selector)\n}\n"
  }
};
})();

(node/*: any*/).hash = "afcad6b0c394ebf36c5cba8588e39e74";

module.exports = ((node/*: any*/)/*: Query<
  ResetDatabaseQuery$variables,
  ResetDatabaseQuery$data,
>*/);

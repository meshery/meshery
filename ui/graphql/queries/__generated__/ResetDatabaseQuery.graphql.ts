/**
 * @generated SignedSource<<79475d75b004fdc1055ac1499ccaaf8b>>
 * @lightSyntaxTransform
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type ReSyncActions = {
  ReSync: string;
  clearDB: string;
  hardReset: string;
};
export type ResetDatabaseQuery$variables = {
  k8scontextID: string;
  selector: ReSyncActions;
};
export type ResetDatabaseQuery$data = {
  readonly resetStatus: Status;
};
export type ResetDatabaseQuery = {
  response: ResetDatabaseQuery$data;
  variables: ResetDatabaseQuery$variables;
};

const node: ConcreteRequest = (function(){
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
      (v0/*:: as any*/),
      (v1/*:: as any*/)
    ],
    "kind": "Fragment",
    "metadata": null,
    "name": "ResetDatabaseQuery",
    "selections": (v2/*:: as any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v1/*:: as any*/),
      (v0/*:: as any*/)
    ],
    "kind": "Operation",
    "name": "ResetDatabaseQuery",
    "selections": (v2/*:: as any*/)
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

(node as any).hash = "54a9344cc4d95023f5082936dc95d05d";

export default node;

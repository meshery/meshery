/**
 * @generated SignedSource<<a78340b6e85371d1661042de36ecf59e>>
 * @lightSyntaxTransform
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type MesheryController = "BROKER" | "MESHSYNC" | "OPERATOR" | "%future added value";
export type MesheryControllerStatus = "CONNECTED" | "DEPLOYED" | "DEPLOYING" | "ENABLED" | "NOTDEPLOYED" | "RUNNING" | "UNDEPLOYED" | "UNKOWN" | "%future added value";
export type MesheryControllersStatusSubscription$variables = {
  connectionIDs?: ReadonlyArray<string> | null | undefined;
};
export type MesheryControllersStatusSubscription$data = {
  readonly subscribeMesheryControllersStatus: ReadonlyArray<{
    readonly connectionID: string;
    readonly controller: MesheryController;
    readonly status: MesheryControllerStatus;
    readonly version: string;
  }>;
};
export type MesheryControllersStatusSubscription = {
  response: MesheryControllersStatusSubscription$data;
  variables: MesheryControllersStatusSubscription$variables;
};

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "connectionIDs"
  }
],
v1 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "connectionIDs",
        "variableName": "connectionIDs"
      }
    ],
    "concreteType": "MesheryControllersStatusListItem",
    "kind": "LinkedField",
    "name": "subscribeMesheryControllersStatus",
    "plural": true,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "connectionID",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "controller",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "status",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "version",
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*:: as any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "MesheryControllersStatusSubscription",
    "selections": (v1/*:: as any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*:: as any*/),
    "kind": "Operation",
    "name": "MesheryControllersStatusSubscription",
    "selections": (v1/*:: as any*/)
  },
  "params": {
    "cacheID": "534d4b9391826e634ad2ab05b1a30e75",
    "id": null,
    "metadata": {},
    "name": "MesheryControllersStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription MesheryControllersStatusSubscription(\n  $connectionIDs: [String!]\n) {\n  subscribeMesheryControllersStatus(connectionIDs: $connectionIDs) {\n    connectionID\n    controller\n    status\n    version\n  }\n}\n"
  }
};
})();

(node as any).hash = "57fdebbc7ad327ad4b1607c9fa7df9e2";

export default node;

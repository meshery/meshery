/**
 * @generated SignedSource<<f39d23088d57dc8c3b41374519736d28>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL_MESH" | "APP_MESH" | "CILIUM_SERVICE_MESH" | "CITRIX_SERVICE_MESH" | "CONSUL" | "INVALID_MESH" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORK_SERVICE_MESH" | "NGINX_SERVICE_MESH" | "OCTARINE" | "OPEN_SERVICE_MESH" | "TANZU" | "TRAEFIK_MESH" | "%future added value";
export type ServiceMeshFilter = {
  k8sClusterIDs?: ReadonlyArray<string> | null | undefined;
  type?: MeshType | null | undefined;
};
export type AddonsStatusQuery$variables = {
  filter?: ServiceMeshFilter | null | undefined;
};
export type AddonsStatusQuery$data = {
  readonly addonsState: ReadonlyArray<{
    readonly name: string;
    readonly owner: string;
  }>;
};
export type AddonsStatusQuery = {
  response: AddonsStatusQuery$data;
  variables: AddonsStatusQuery$variables;
};

const node: ConcreteRequest = (function(){
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

(node as any).hash = "9cbf0a827a321dead7e3e6d0c2e9cbe7";

export default node;

/**
 * @generated SignedSource<<85939655c2fe270f9d72e146bb9d1fa0>>
 * @lightSyntaxTransform
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type PageFilter = {
  from?: string | null | undefined;
  metrics?: string | null | undefined;
  order?: string | null | undefined;
  page: string;
  pageSize: string;
  populate?: ReadonlyArray<string | null | undefined> | null | undefined;
  search?: string | null | undefined;
  to?: string | null | undefined;
  updated_after?: string | null | undefined;
  visibility?: ReadonlyArray<string> | null | undefined;
};
export type K8sContextSubscription$variables = {
  selector: PageFilter;
};
export type K8sContextSubscription$data = {
  readonly k8sContext: {
    readonly contexts: ReadonlyArray<{
      readonly connectionId: string;
      readonly createdAt: string;
      readonly createdBy: string;
      readonly deploymentType: string;
      readonly id: string;
      readonly kubernetesServerId: string;
      readonly mesheryInstanceId: string;
      readonly name: string;
      readonly owner: string;
      readonly server: string;
      readonly updatedAt: string;
      readonly version: string;
    } | null | undefined>;
    readonly totalCount: number;
  };
};
export type K8sContextSubscription = {
  response: K8sContextSubscription$data;
  variables: K8sContextSubscription$variables;
};

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "selector"
  }
],
v1 = [
  {
    "alias": "k8sContext",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "K8sContextsPage",
    "kind": "LinkedField",
    "name": "subscribeK8sContext",
    "plural": false,
    "selections": [
      {
        "alias": "totalCount",
        "args": null,
        "kind": "ScalarField",
        "name": "total_count",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "K8sContext",
        "kind": "LinkedField",
        "name": "contexts",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "id",
            "storageKey": null
          },
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
            "name": "server",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "owner",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "createdBy",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "mesheryInstanceId",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "kubernetesServerId",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "deploymentType",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "updatedAt",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "createdAt",
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
            "name": "connectionId",
            "storageKey": null
          }
        ],
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
    "name": "K8sContextSubscription",
    "selections": (v1/*:: as any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*:: as any*/),
    "kind": "Operation",
    "name": "K8sContextSubscription",
    "selections": (v1/*:: as any*/)
  },
  "params": {
    "cacheID": "21e06ecbe7a8788897700273f61d677c",
    "id": null,
    "metadata": {},
    "name": "K8sContextSubscription",
    "operationKind": "subscription",
    "text": "subscription K8sContextSubscription(\n  $selector: PageFilter!\n) {\n  k8sContext: subscribeK8sContext(selector: $selector) {\n    totalCount: total_count\n    contexts {\n      id\n      name\n      server\n      owner\n      createdBy\n      mesheryInstanceId\n      kubernetesServerId\n      deploymentType\n      updatedAt\n      createdAt\n      version\n      connectionId\n    }\n  }\n}\n"
  }
};
})();

(node as any).hash = "77758342e08fc786e2cc7a7af4bf7492";

export default node;

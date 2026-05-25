/**
 * @generated SignedSource<<788e1f1b9e892b70447b1a5b3f09e742>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type NamespaceQuery$variables = {
  k8sClusterIDs?: ReadonlyArray<string> | null | undefined;
};
export type NamespaceQuery$data = {
  readonly namespaces: ReadonlyArray<{
    readonly namespace: string;
  }>;
};
export type NamespaceQuery = {
  response: NamespaceQuery$data;
  variables: NamespaceQuery$variables;
};

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "k8sClusterIDs"
  }
],
v1 = [
  {
    "alias": "namespaces",
    "args": [
      {
        "kind": "Variable",
        "name": "k8sClusterIDs",
        "variableName": "k8sClusterIDs"
      }
    ],
    "concreteType": "NameSpace",
    "kind": "LinkedField",
    "name": "getAvailableNamespaces",
    "plural": true,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "namespace",
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
    "name": "NamespaceQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "NamespaceQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "0e85918419a43e5791cd5a51e7f6896e",
    "id": null,
    "metadata": {},
    "name": "NamespaceQuery",
    "operationKind": "query",
    "text": "query NamespaceQuery(\n  $k8sClusterIDs: [String!]\n) {\n  namespaces: getAvailableNamespaces(k8sClusterIDs: $k8sClusterIDs) {\n    namespace\n  }\n}\n"
  }
};
})();

(node as any).hash = "04f74232907aa0ba765bd0f8db6c427c";

export default node;

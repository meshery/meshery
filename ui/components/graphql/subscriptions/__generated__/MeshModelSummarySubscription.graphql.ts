/**
 * @generated SignedSource<<cbc0fe6f2ae5a438c913eceb1c176790>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from 'relay-runtime';
export type MeshModelSummarySelector = {
  type: string;
};
export type MeshModelSummarySubscription$variables = {
  selector: MeshModelSummarySelector;
};
export type MeshModelSummarySubscription$data = {
  readonly meshmodelSummary: {
    readonly components: ReadonlyArray<{
      readonly count: number;
      readonly name: string;
    }> | null | undefined;
    readonly relationships: ReadonlyArray<{
      readonly count: number;
      readonly name: string;
    }> | null | undefined;
  };
};
export type MeshModelSummarySubscription = {
  response: MeshModelSummarySubscription$data;
  variables: MeshModelSummarySubscription$variables;
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
    "name": "count",
    "storageKey": null
  }
],
v2 = [
  {
    "alias": "meshmodelSummary",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "MeshModelSummary",
    "kind": "LinkedField",
    "name": "subscribeMeshModelSummary",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "MeshModelComponent",
        "kind": "LinkedField",
        "name": "components",
        "plural": true,
        "selections": (v1/*: any*/),
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "MeshModelRelationship",
        "kind": "LinkedField",
        "name": "relationships",
        "plural": true,
        "selections": (v1/*: any*/),
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
    "name": "MeshModelSummarySubscription",
    "selections": (v2/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshModelSummarySubscription",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "458b92d09a21e56a16ddbc316c2f1596",
    "id": null,
    "metadata": {},
    "name": "MeshModelSummarySubscription",
    "operationKind": "subscription",
    "text": "subscription MeshModelSummarySubscription(\n  $selector: MeshModelSummarySelector!\n) {\n  meshmodelSummary: subscribeMeshModelSummary(selector: $selector) {\n    components {\n      name\n      count\n    }\n    relationships {\n      name\n      count\n    }\n  }\n}\n"
  }
};
})();

(node as any).hash = "f9ce6171759cd696ee275c65c930974b";

export default node;

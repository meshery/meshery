/**
 * @generated SignedSource<<c9f98c08ee9ce19ee1f7e1d5951505ab>>
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
export type MeshModelSummaryQuery$variables = {
  selector: MeshModelSummarySelector;
};
export type MeshModelSummaryQuery$data = {
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
export type MeshModelSummaryQuery = {
  response: MeshModelSummaryQuery$data;
  variables: MeshModelSummaryQuery$variables;
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
    "name": "getMeshModelSummary",
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
    "name": "MeshModelSummaryQuery",
    "selections": (v2/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshModelSummaryQuery",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "fb8718ecc59096bb7a68d0479193aaab",
    "id": null,
    "metadata": {},
    "name": "MeshModelSummaryQuery",
    "operationKind": "query",
    "text": "query MeshModelSummaryQuery(\n  $selector: MeshModelSummarySelector!\n) {\n  meshmodelSummary: getMeshModelSummary(selector: $selector) {\n    components {\n      name\n      count\n    }\n    relationships {\n      name\n      count\n    }\n  }\n}\n"
  }
};
})();

(node as any).hash = "830ea0f72a52a2f3419b0796d9b3a562";

export default node;

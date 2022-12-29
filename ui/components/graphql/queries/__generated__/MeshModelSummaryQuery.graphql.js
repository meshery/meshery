/**
 * @generated SignedSource<<97cb8728db91c2cde8eeb0e8c9eaf6d3>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, Query } from 'relay-runtime';
export type MeshModelSummarySelector = {|
  type: string,
|};
export type MeshModelSummaryQuery$variables = {|
  selector: MeshModelSummarySelector,
|};
export type MeshModelSummaryQuery$data = {|
  +meshmodelSummary: {|
    +components: ?$ReadOnlyArray<{|
      +name: string,
      +count: number,
    |}>,
    +relationships: ?$ReadOnlyArray<{|
      +kind: string,
      +count: number,
    |}>,
  |},
|};
export type MeshModelSummaryQuery = {|
  variables: MeshModelSummaryQuery$variables,
  response: MeshModelSummaryQuery$data,
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
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "count",
  "storageKey": null
},
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
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "name",
            "storageKey": null
          },
          (v1/*: any*/)
        ],
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "MeshModelRelationship",
        "kind": "LinkedField",
        "name": "relationships",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "kind",
            "storageKey": null
          },
          (v1/*: any*/)
        ],
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
    "cacheID": "763e9b1c1381f34d10ec4d6acb511612",
    "id": null,
    "metadata": {},
    "name": "MeshModelSummaryQuery",
    "operationKind": "query",
    "text": "query MeshModelSummaryQuery(\n  $selector: MeshModelSummarySelector!\n) {\n  meshmodelSummary: getMeshModelSummary(selector: $selector) {\n    components {\n      name\n      count\n    }\n    relationships {\n      kind\n      count\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "58e0fe0f73c993bc3b8d94581034494e";

module.exports = ((node/*: any*/)/*: Query<
  MeshModelSummaryQuery$variables,
  MeshModelSummaryQuery$data,
>*/);

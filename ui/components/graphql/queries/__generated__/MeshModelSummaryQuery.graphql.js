/**
 * @generated SignedSource<<d3841b2609d7137b27d2a46fe4250cff>>
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
      +name: string,
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

(node/*: any*/).hash = "830ea0f72a52a2f3419b0796d9b3a562";

module.exports = ((node/*: any*/)/*: Query<
  MeshModelSummaryQuery$variables,
  MeshModelSummaryQuery$data,
>*/);

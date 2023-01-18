/**
 * @generated SignedSource<<162cc142fdd38d257fa9bf34317e7677>>
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
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "count",
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
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "MeshModelSummaryQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshModelSummaryQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "3591b13ea3bc9d4aa8b7adaa55ef3981",
    "id": null,
    "metadata": {},
    "name": "MeshModelSummaryQuery",
    "operationKind": "query",
    "text": "query MeshModelSummaryQuery(\n  $selector: MeshModelSummarySelector!\n) {\n  meshmodelSummary: getMeshModelSummary(selector: $selector) {\n    components {\n      name\n      count\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "12fdc86cb35a88e25632cd3b3d830a77";

module.exports = ((node/*: any*/)/*: Query<
  MeshModelSummaryQuery$variables,
  MeshModelSummaryQuery$data,
>*/);

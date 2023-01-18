/**
 * @generated SignedSource<<0967080ea24a31d0abac0519b441e3b2>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type MeshModelSummarySelector = {|
  type: string,
|};
export type MeshModelSummarySubscription$variables = {|
  selector: MeshModelSummarySelector,
|};
export type MeshModelSummarySubscription$data = {|
  +meshmodelSummary: {|
    +components: ?$ReadOnlyArray<{|
      +name: string,
      +count: number,
    |}>,
  |},
|};
export type MeshModelSummarySubscription = {|
  variables: MeshModelSummarySubscription$variables,
  response: MeshModelSummarySubscription$data,
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
    "name": "MeshModelSummarySubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MeshModelSummarySubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "e3df394d3f5405ea5cb8d232e5aa9fc2",
    "id": null,
    "metadata": {},
    "name": "MeshModelSummarySubscription",
    "operationKind": "subscription",
    "text": "subscription MeshModelSummarySubscription(\n  $selector: MeshModelSummarySelector!\n) {\n  meshmodelSummary: subscribeMeshModelSummary(selector: $selector) {\n    components {\n      name\n      count\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "df62dc502a31750c294f4bbabf29c4e0";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  MeshModelSummarySubscription$variables,
  MeshModelSummarySubscription$data,
>*/);

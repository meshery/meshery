/**
 * @generated SignedSource<<8dc2ae6951a9d6ff15b5c01022e30a71>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type MesheryController = "BROKER" | "OPERATOR" | "MESHSYNC" | "%future added value";
export type MesheryControllerStatus = "DEPLOYED" | "NOTDEPLOYED" | "DEPLOYING" | "UNKOWN" | "%future added value";
export type MesheryControllersStatusSubscription$variables = {|
  k8scontextIDs?: ?$ReadOnlyArray<string>,
|};
export type MesheryControllersStatusSubscription$data = {|
  +subscribeMesheryControllersStatus: $ReadOnlyArray<{|
    +contextId: string,
    +controller: MesheryController,
    +status: MesheryControllerStatus,
  |}>,
|};
export type MesheryControllersStatusSubscription = {|
  variables: MesheryControllersStatusSubscription$variables,
  response: MesheryControllersStatusSubscription$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "k8scontextIDs"
  }
],
v1 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextIDs",
        "variableName": "k8scontextIDs"
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
        "name": "contextId",
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
    "name": "MesheryControllersStatusSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "MesheryControllersStatusSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "2f51de3afecd08d2d6316973805b3890",
    "id": null,
    "metadata": {},
    "name": "MesheryControllersStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription MesheryControllersStatusSubscription(\n  $k8scontextIDs: [String!]\n) {\n  subscribeMesheryControllersStatus(k8scontextIDs: $k8scontextIDs) {\n    contextId\n    controller\n    status\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "5a93ffadfcfbf8ea14d9bac0fe6b50f3";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  MesheryControllersStatusSubscription$variables,
  MesheryControllersStatusSubscription$data,
>*/);

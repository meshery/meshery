/**
 * @generated SignedSource<<4ba67ecbe17a56aaa4d3b68dfaba890e>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type PageFilter = {|
  page: string,
  pageSize: string,
  order?: ?string,
  search?: ?string,
  from?: ?string,
  to?: ?string,
|};
export type K8sContextSubscription$variables = {|
  selector: PageFilter,
|};
export type K8sContextSubscription$data = {|
  +k8sContext: {|
    +total_count: number,
    +contexts: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +auth: any,
      +cluster: any,
      +server: string,
      +owner: string,
      +created_by: string,
      +meshery_instance_id: string,
      +kubernetes_server_id: string,
      +updated_at: string,
      +created_at: string,
    |}>,
  |},
|};
export type K8sContextSubscription = {|
  variables: K8sContextSubscription$variables,
  response: K8sContextSubscription$data,
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
        "alias": null,
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
            "name": "auth",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "cluster",
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
            "name": "created_by",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "meshery_instance_id",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "kubernetes_server_id",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "updated_at",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "created_at",
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
    "name": "K8sContextSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "K8sContextSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "0a9c00eaa40ee17df30b372f3698d6f4",
    "id": null,
    "metadata": {},
    "name": "K8sContextSubscription",
    "operationKind": "subscription",
    "text": "subscription K8sContextSubscription(\n  $selector: PageFilter!\n) {\n  k8sContext: subscribeK8sContext(selector: $selector) {\n    total_count\n    contexts {\n      id\n      name\n      auth\n      cluster\n      server\n      owner\n      created_by\n      meshery_instance_id\n      kubernetes_server_id\n      updated_at\n      created_at\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "87663c63ebdcff6e185259074bd22d7f";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  K8sContextSubscription$variables,
  K8sContextSubscription$data,
>*/);

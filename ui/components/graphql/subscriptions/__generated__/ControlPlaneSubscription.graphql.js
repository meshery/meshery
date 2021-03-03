/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type ControlPlaneFilter = {|
  type?: ?MeshType
|};
export type ControlPlaneSubscriptionVariables = {|
  filter?: ?ControlPlaneFilter
|};
export type ControlPlaneSubscriptionResponse = {|
  +listenToControlPlaneState: $ReadOnlyArray<{|
    +name: ?MeshType,
    +members: $ReadOnlyArray<{|
      +version: string,
      +component: string,
      +namespace: string,
    |}>,
  |}>
|};
export type ControlPlaneSubscription = {|
  variables: ControlPlaneSubscriptionVariables,
  response: ControlPlaneSubscriptionResponse,
|};
*/


/*
subscription ControlPlaneSubscription(
  $filter: ControlPlaneFilter
) {
  listenToControlPlaneState(filter: $filter) {
    name
    members {
      version
      component
      namespace
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "filter"
  }
],
v1 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "filter",
        "variableName": "filter"
      }
    ],
    "concreteType": "ControlPlane",
    "kind": "LinkedField",
    "name": "listenToControlPlaneState",
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
        "concreteType": "ControlPlaneMember",
        "kind": "LinkedField",
        "name": "members",
        "plural": true,
        "selections": [
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
            "name": "component",
            "storageKey": null
          },
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
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "ControlPlaneSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ControlPlaneSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "43ea5e5f46bb37e5ad389cd2db5b3043",
    "id": null,
    "metadata": {},
    "name": "ControlPlaneSubscription",
    "operationKind": "subscription",
    "text": "subscription ControlPlaneSubscription(\n  $filter: ControlPlaneFilter\n) {\n  listenToControlPlaneState(filter: $filter) {\n    name\n    members {\n      version\n      component\n      namespace\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '4de52f8e59b0d8768dba3ada886d029d';

module.exports = node;

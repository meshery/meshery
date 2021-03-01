/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type Status = "DISABLED" | "ENABLED" | "UNKNOWN" | "%future added value";
export type ControlPlaneFilter = {|
  type?: ?MeshType
|};
export type meshScanSubscriptionVariables = {|
  filter?: ?ControlPlaneFilter
|};
export type meshScanSubscriptionResponse = {|
  +listenToControlPlaneEvents: $ReadOnlyArray<{|
    +name: ?MeshType,
    +version: string,
    +members: $ReadOnlyArray<{|
      +component: string,
      +status: ?Status,
      +namespace: string,
    |}>,
  |}>
|};
export type meshScanSubscription = {|
  variables: meshScanSubscriptionVariables,
  response: meshScanSubscriptionResponse,
|};
*/


/*
subscription meshScanSubscription(
  $filter: ControlPlaneFilter
) {
  listenToControlPlaneEvents(filter: $filter) {
    name
    version
    members {
      component
      status
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
    "name": "listenToControlPlaneEvents",
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
        "name": "version",
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
            "name": "component",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "status",
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
    "name": "meshScanSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "meshScanSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "610f5d659d8674ba54168d715e9ec97b",
    "id": null,
    "metadata": {},
    "name": "meshScanSubscription",
    "operationKind": "subscription",
    "text": "subscription meshScanSubscription(\n  $filter: ControlPlaneFilter\n) {\n  listenToControlPlaneEvents(filter: $filter) {\n    name\n    version\n    members {\n      component\n      status\n      namespace\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '82a586a9835603dc4769821bfd5173e0';

module.exports = node;

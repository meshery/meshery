/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type AddonStatusSubscriptionVariables = {|
  selector?: ?MeshType
|};
export type AddonStatusSubscriptionResponse = {|
  +addonEvent: $ReadOnlyArray<{|
    +type: string,
    +status: ?Status,
    +config: {|
      +serviceName: string,
      +endpoint: string,
    |},
  |}>
|};
export type AddonStatusSubscription = {|
  variables: AddonStatusSubscriptionVariables,
  response: AddonStatusSubscriptionResponse,
|};
*/


/*
subscription AddonStatusSubscription(
  $selector: MeshType
) {
  addonEvent: listenToAddonState(selector: $selector) {
    type
    status
    config {
      serviceName
      endpoint
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "selector"
  }
],
v1 = [
  {
    "alias": "addonEvent",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "AddonList",
    "kind": "LinkedField",
    "name": "listenToAddonState",
    "plural": true,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "type",
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
        "concreteType": "AddonConfig",
        "kind": "LinkedField",
        "name": "config",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "serviceName",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "endpoint",
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
    "name": "AddonStatusSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AddonStatusSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "f249e651031639a7146c1ba3e0f3f592",
    "id": null,
    "metadata": {},
    "name": "AddonStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription AddonStatusSubscription(\n  $selector: MeshType\n) {\n  addonEvent: listenToAddonState(selector: $selector) {\n    type\n    status\n    config {\n      serviceName\n      endpoint\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '08f8b6677edf887e93a65b89731a649e';

module.exports = node;

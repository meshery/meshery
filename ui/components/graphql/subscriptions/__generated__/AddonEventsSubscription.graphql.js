/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type Status = "DISABLED" | "ENABLED" | "UNKNOWN" | "%future added value";
export type AddonEventsSubscriptionVariables = {|
  selector?: ?MeshType
|};
export type AddonEventsSubscriptionResponse = {|
  +listenToAddonEvents: $ReadOnlyArray<{|
    +type: string,
    +status: ?Status,
  |}>
|};
export type AddonEventsSubscription = {|
  variables: AddonEventsSubscriptionVariables,
  response: AddonEventsSubscriptionResponse,
|};
*/


/*
subscription AddonEventsSubscription(
  $selector: MeshType
) {
  listenToAddonEvents(selector: $selector) {
    type
    status
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
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "AddonList",
    "kind": "LinkedField",
    "name": "listenToAddonEvents",
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
    "name": "AddonEventsSubscription",
    "selections": (v1/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AddonEventsSubscription",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "de2fe45019d12d94969793bd803ec58d",
    "id": null,
    "metadata": {},
    "name": "AddonEventsSubscription",
    "operationKind": "subscription",
    "text": "subscription AddonEventsSubscription(\n  $selector: MeshType\n) {\n  listenToAddonEvents(selector: $selector) {\n    type\n    status\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ec53be7e58e9ad99e7c1ab34e3612d5d';

module.exports = node;

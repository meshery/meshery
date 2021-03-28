/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "CITRIX" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NGINXSM" | "NSM" | "OCTARINE" | "OSM" | "TRAEFIKMESH" | "%future added value";
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type AddonStatusMutationVariables = {|
  selector?: ?MeshType,
  targetStatus: Status,
|};
export type AddonStatusMutationResponse = {|
  +addonstate: Status
|};
export type AddonStatusMutation = {|
  variables: AddonStatusMutationVariables,
  response: AddonStatusMutationResponse,
|};
*/


/*
mutation AddonStatusMutation(
  $selector: MeshType
  $targetStatus: Status!
) {
  addonstate: changeAddonStatus(selector: $selector, targetStatus: $targetStatus)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "selector"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "targetStatus"
  }
],
v1 = [
  {
    "alias": "addonstate",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      },
      {
        "kind": "Variable",
        "name": "targetStatus",
        "variableName": "targetStatus"
      }
    ],
    "kind": "ScalarField",
    "name": "changeAddonStatus",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "AddonStatusMutation",
    "selections": (v1/*: any*/),
    "type": "Mutation",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AddonStatusMutation",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "bc81997a1e9bb3d2db83a11e0a4150b7",
    "id": null,
    "metadata": {},
    "name": "AddonStatusMutation",
    "operationKind": "mutation",
    "text": "mutation AddonStatusMutation(\n  $selector: MeshType\n  $targetStatus: Status!\n) {\n  addonstate: changeAddonStatus(selector: $selector, targetStatus: $targetStatus)\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '37539d93c0e9392117292e6633d6be1c';

module.exports = node;

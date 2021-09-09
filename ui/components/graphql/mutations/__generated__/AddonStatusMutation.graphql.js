/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL_MESH" | "APP_MESH" | "CITRIX_SERVICE_MESH" | "CONSUL" | "INVALID_MESH" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORK_SERVICE_MESH" | "NGINX_SERVICE_MESH" | "OCTARINE" | "OPEN_SERVICE_MESH" | "TANZU" | "TRAEFIK_MESH" | "%future added value";
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type AddonStatusInput = {|
  selector?: ?MeshType,
  targetStatus: Status,
|};
export type AddonStatusMutationVariables = {|
  input?: ?AddonStatusInput
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
  $input: AddonStatusInput
) {
  addonstate: changeAddonStatus(input: $input)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "input"
  }
],
v1 = [
  {
    "alias": "addonstate",
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
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
    "cacheID": "4f11c7d338b4a51b255fae5554424d74",
    "id": null,
    "metadata": {},
    "name": "AddonStatusMutation",
    "operationKind": "mutation",
    "text": "mutation AddonStatusMutation(\n  $input: AddonStatusInput\n) {\n  addonstate: changeAddonStatus(input: $input)\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'eaed3737aff98e45fc02f51c97601092';

module.exports = node;

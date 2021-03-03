/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type AddonsStatusQueryVariables = {|
  selector?: ?MeshType
|};
export type AddonsStatusQueryResponse = {|
  +addons: $ReadOnlyArray<{|
    +type: string,
    +status: ?Status,
    +config: {|
      +serviceName: string,
      +endpoint: string,
    |},
  |}>
|};
export type AddonsStatusQuery = {|
  variables: AddonsStatusQueryVariables,
  response: AddonsStatusQueryResponse,
|};
*/


/*
query AddonsStatusQuery(
  $selector: MeshType
) {
  addons: getAvailableAddons(selector: $selector) {
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
    "alias": "addons",
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "selector"
      }
    ],
    "concreteType": "AddonList",
    "kind": "LinkedField",
    "name": "getAvailableAddons",
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
    "name": "AddonsStatusQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "AddonsStatusQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "a0e1ea9e2b4df21ed735fcd2dd472dd7",
    "id": null,
    "metadata": {},
    "name": "AddonsStatusQuery",
    "operationKind": "query",
    "text": "query AddonsStatusQuery(\n  $selector: MeshType\n) {\n  addons: getAvailableAddons(selector: $selector) {\n    type\n    status\n    config {\n      serviceName\n      endpoint\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '4cec6d07aadb1179ffa543b0f83d3acc';

module.exports = node;

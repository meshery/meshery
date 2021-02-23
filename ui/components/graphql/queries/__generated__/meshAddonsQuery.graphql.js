/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type Status = "DISABLED" | "ENABLED" | "UNKNOWN" | "%future added value";
export type meshAddonsQueryVariables = {|
  selector?: ?MeshType
|};
export type meshAddonsQueryResponse = {|
  +addons: $ReadOnlyArray<{|
    +type: string,
    +status: ?Status,
    +config: {|
      +serviceName: string
    |},
  |}>
|};
export type meshAddonsQuery = {|
  variables: meshAddonsQueryVariables,
  response: meshAddonsQueryResponse,
|};
*/


/*
query meshAddonsQuery(
  $selector: MeshType
) {
  addons: getAvailableAddons(selector: $selector) {
    type
    status
    config {
      serviceName
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
    "name": "meshAddonsQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "meshAddonsQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "358de11af88770b3cf9e6c5899f357cb",
    "id": null,
    "metadata": {},
    "name": "meshAddonsQuery",
    "operationKind": "query",
    "text": "query meshAddonsQuery(\n  $selector: MeshType\n) {\n  addons: getAvailableAddons(selector: $selector) {\n    type\n    status\n    config {\n      serviceName\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'fdb358edef8fb4bf80b55772537eab80';

module.exports = node;

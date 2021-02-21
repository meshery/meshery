/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type AddonSelector = "GRAFANA" | "JAEGER" | "PROMETHEUS" | "ZIPKIN" | "%future added value";
export type MeshType = "ALL" | "CITRIXSM" | "CONSUL" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORKSM" | "NGINXSM" | "NONE" | "OCTARINE" | "OPENSERVICEMESH" | "TRAEFIK" | "%future added value";
export type Status = "DISABLED" | "ENABLED" | "UNKNOWN" | "%future added value";
export type meshAddonsQueryVariables = {|
  meshType?: ?MeshType
|};
export type meshAddonsQueryResponse = {|
  +getAvailableAddons: $ReadOnlyArray<{|
    +type: ?AddonSelector,
    +status: ?Status,
  |}>
|};
export type meshAddonsQuery = {|
  variables: meshAddonsQueryVariables,
  response: meshAddonsQueryResponse,
|};
*/


/*
query meshAddonsQuery(
  $meshType: MeshType
) {
  getAvailableAddons(selector: $meshType) {
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
    "name": "meshType"
  }
],
v1 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "selector",
        "variableName": "meshType"
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
    "cacheID": "20a323bc0410ad564a6278b8579b9f93",
    "id": null,
    "metadata": {},
    "name": "meshAddonsQuery",
    "operationKind": "query",
    "text": "query meshAddonsQuery(\n  $meshType: MeshType\n) {\n  getAvailableAddons(selector: $meshType) {\n    type\n    status\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd9572afc0f3d53ec5e45a99445693d24';

module.exports = node;

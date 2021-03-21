/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type NamespaceQueryVariables = {||};
export type NamespaceQueryResponse = {|
  +namespaces: $ReadOnlyArray<{|
    +namespace: string
  |}>
|};
export type NamespaceQuery = {|
  variables: NamespaceQueryVariables,
  response: NamespaceQueryResponse,
|};
*/


/*
query NamespaceQuery {
  namespaces: getAvailableNamespaces {
    namespace
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": "namespaces",
    "args": null,
    "concreteType": "NameSpace",
    "kind": "LinkedField",
    "name": "getAvailableNamespaces",
    "plural": true,
    "selections": [
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
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "NamespaceQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "NamespaceQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "46c4030fb0b28da4b5ce74853e237155",
    "id": null,
    "metadata": {},
    "name": "NamespaceQuery",
    "operationKind": "query",
    "text": "query NamespaceQuery {\n  namespaces: getAvailableNamespaces {\n    namespace\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '683a843465f9f33c9e162cdb23a14e05';

module.exports = node;

/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "CONNECTED" | "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type NatsStatusQueryVariables = {||};
export type NatsStatusQueryResponse = {|
  +controller: {|
    +name: string,
    +version: string,
    +status: Status,
  |}
|};
export type NatsStatusQuery = {|
  variables: NatsStatusQueryVariables,
  response: NatsStatusQueryResponse,
|};
*/


/*
query NatsStatusQuery {
  controller: getNatsStatus {
    name
    version
    status
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": "controller",
    "args": null,
    "concreteType": "OperatorControllerStatus",
    "kind": "LinkedField",
    "name": "getNatsStatus",
    "plural": false,
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
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "NatsStatusQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "NatsStatusQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "94428d644c1ee447c9b118e7771847a0",
    "id": null,
    "metadata": {},
    "name": "NatsStatusQuery",
    "operationKind": "query",
    "text": "query NatsStatusQuery {\n  controller: getNatsStatus {\n    name\n    version\n    status\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd558c76b74dcec56bf23feac4aebbb9f';

module.exports = node;

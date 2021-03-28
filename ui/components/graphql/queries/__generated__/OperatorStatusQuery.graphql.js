/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorStatusQueryVariables = {||};
export type OperatorStatusQueryResponse = {|
  +operator: ?{|
    +status: Status,
    +version: string,
    +controllers: $ReadOnlyArray<{|
      +name: string,
      +version: string,
      +status: Status,
    |}>,
    +error: ?{|
      +code: string,
      +description: string,
    |},
  |}
|};
export type OperatorStatusQuery = {|
  variables: OperatorStatusQueryVariables,
  response: OperatorStatusQueryResponse,
|};
*/


/*
query OperatorStatusQuery {
  operator: getOperatorStatus {
    status
    version
    controllers {
      name
      version
      status
    }
    error {
      code
      description
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "status",
  "storageKey": null
},
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "version",
  "storageKey": null
},
v2 = [
  {
    "alias": "operator",
    "args": null,
    "concreteType": "OperatorStatus",
    "kind": "LinkedField",
    "name": "getOperatorStatus",
    "plural": false,
    "selections": [
      (v0/*: any*/),
      (v1/*: any*/),
      {
        "alias": null,
        "args": null,
        "concreteType": "OperatorControllerStatus",
        "kind": "LinkedField",
        "name": "controllers",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "name",
            "storageKey": null
          },
          (v1/*: any*/),
          (v0/*: any*/)
        ],
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "Error",
        "kind": "LinkedField",
        "name": "error",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "code",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "description",
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
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "OperatorStatusQuery",
    "selections": (v2/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "OperatorStatusQuery",
    "selections": (v2/*: any*/)
  },
  "params": {
    "cacheID": "fe94be1b823f90868194ca09bfd4de28",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusQuery",
    "operationKind": "query",
    "text": "query OperatorStatusQuery {\n  operator: getOperatorStatus {\n    status\n    version\n    controllers {\n      name\n      version\n      status\n    }\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '9d8069b0d65a41383aba2700c5f837d1';

module.exports = node;

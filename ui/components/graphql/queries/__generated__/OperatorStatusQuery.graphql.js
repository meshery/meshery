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
    error {
      code
      description
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": "operator",
    "args": null,
    "concreteType": "OperatorStatus",
    "kind": "LinkedField",
    "name": "getOperatorStatus",
    "plural": false,
    "selections": [
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
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "OperatorStatusQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "ed65e93adf8349ea7d5a4e47f08bd16f",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusQuery",
    "operationKind": "query",
    "text": "query OperatorStatusQuery {\n  operator: getOperatorStatus {\n    status\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '1c2d4cc73e3bbacd3ad55c6f2fb8c4bd';

module.exports = node;

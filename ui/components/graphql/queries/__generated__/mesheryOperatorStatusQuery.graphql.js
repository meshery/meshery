/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type Status = "DISABLED" | "ENABLED" | "UNKNOWN" | "%future added value";
export type mesheryOperatorStatusQueryVariables = {||};
export type mesheryOperatorStatusQueryResponse = {|
  +operator: ?{|
    +status: ?Status,
    +error: ?{|
      +code: string,
      +description: string,
    |},
  |}
|};
export type mesheryOperatorStatusQuery = {|
  variables: mesheryOperatorStatusQueryVariables,
  response: mesheryOperatorStatusQueryResponse,
|};
*/


/*
query mesheryOperatorStatusQuery {
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
    "name": "mesheryOperatorStatusQuery",
    "selections": (v0/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "mesheryOperatorStatusQuery",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "0e7003f2d4816ecf4e6864a84333f5a8",
    "id": null,
    "metadata": {},
    "name": "mesheryOperatorStatusQuery",
    "operationKind": "query",
    "text": "query mesheryOperatorStatusQuery {\n  operator: getOperatorStatus {\n    status\n    error {\n      code\n      description\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '800b4e252bc857ec45413b673f1c8a0d';

module.exports = node;

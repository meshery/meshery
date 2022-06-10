/**
 * @generated SignedSource<<a71c67c0427ff681836b039233aa40d9>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type Status = "ENABLED" | "CONNECTED" | "DISABLED" | "PROCESSING" | "UNKNOWN" | "%future added value";
export type OperatorStatusSubscription$variables = {|
  k8scontextIDs?: ?$ReadOnlyArray<string>,
|};
export type OperatorStatusSubscription$data = {|
  +operator: ?{|
    +contextID: string,
    +operatorStatus: {|
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
    |},
  |},
|};
export type OperatorStatusSubscription = {|
  variables: OperatorStatusSubscription$variables,
  response: OperatorStatusSubscription$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "k8scontextIDs"
  }
],
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "status",
  "storageKey": null
},
v2 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "version",
  "storageKey": null
},
v3 = [
  {
    "alias": "operator",
    "args": [
      {
        "kind": "Variable",
        "name": "k8scontextIDs",
        "variableName": "k8scontextIDs"
      }
    ],
    "concreteType": "OperatorStatusPerK8sContext",
    "kind": "LinkedField",
    "name": "listenToOperatorState",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "contextID",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "OperatorStatus",
        "kind": "LinkedField",
        "name": "operatorStatus",
        "plural": false,
        "selections": [
          (v1/*: any*/),
          (v2/*: any*/),
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
              (v2/*: any*/),
              (v1/*: any*/)
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
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "OperatorStatusSubscription",
    "selections": (v3/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "OperatorStatusSubscription",
    "selections": (v3/*: any*/)
  },
  "params": {
    "cacheID": "8a1bb7537702e43181f395fdcaeceaa0",
    "id": null,
    "metadata": {},
    "name": "OperatorStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription OperatorStatusSubscription(\n  $k8scontextIDs: [String!]\n) {\n  operator: listenToOperatorState(k8scontextIDs: $k8scontextIDs) {\n    contextID\n    operatorStatus {\n      status\n      version\n      controllers {\n        name\n        version\n        status\n      }\n      error {\n        code\n        description\n      }\n    }\n  }\n}\n"
  }
};
})();

(node/*: any*/).hash = "5728620b5666bd13a414080f9d90778e";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  OperatorStatusSubscription$variables,
  OperatorStatusSubscription$data,
>*/);

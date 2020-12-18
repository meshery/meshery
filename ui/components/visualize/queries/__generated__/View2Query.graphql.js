/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type View2QueryVariables = {|
  namespaceID?: ?string,
  deploymentID?: ?string,
  showClusters: boolean,
  showNamespaces: boolean,
  showDeployments: boolean,
  showPods: boolean,
|};
export type View2QueryResponse = {|
  +cluster: ?$ReadOnlyArray<?{|
    +id?: string,
    +namespaces: ?$ReadOnlyArray<?{|
      +id?: string,
      +deployments: ?$ReadOnlyArray<?{|
        +id?: string,
        +pods: ?$ReadOnlyArray<?{|
          +id?: string,
          +name?: string,
        |}>,
      |}>,
    |}>,
  |}>
|};
export type View2Query = {|
  variables: View2QueryVariables,
  response: View2QueryResponse,
|};
*/


/*
query View2Query(
  $namespaceID: ID
  $deploymentID: ID
  $showClusters: Boolean!
  $showNamespaces: Boolean!
  $showDeployments: Boolean!
  $showPods: Boolean!
) {
  cluster {
    id @include(if: $showClusters)
    namespaces(namespaceid: $namespaceID) {
      id @include(if: $showNamespaces)
      deployments(deploymentid: $deploymentID) {
        id @include(if: $showDeployments)
        pods {
          id @include(if: $showPods)
          name @include(if: $showPods)
          id
        }
        id
      }
      id
    }
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "deploymentID"
},
v1 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "namespaceID"
},
v2 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "showClusters"
},
v3 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "showDeployments"
},
v4 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "showNamespaces"
},
v5 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "showPods"
},
v6 = [
  {
    "kind": "Variable",
    "name": "namespaceid",
    "variableName": "namespaceID"
  }
],
v7 = [
  {
    "kind": "Variable",
    "name": "deploymentid",
    "variableName": "deploymentID"
  }
],
v8 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v9 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "name",
  "storageKey": null
},
v10 = [
  (v8/*: any*/)
];
return {
  "fragment": {
    "argumentDefinitions": [
      (v0/*: any*/),
      (v1/*: any*/),
      (v2/*: any*/),
      (v3/*: any*/),
      (v4/*: any*/),
      (v5/*: any*/)
    ],
    "kind": "Fragment",
    "metadata": null,
    "name": "View2Query",
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Cluster",
        "kind": "LinkedField",
        "name": "cluster",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": (v6/*: any*/),
            "concreteType": "Namespace",
            "kind": "LinkedField",
            "name": "namespaces",
            "plural": true,
            "selections": [
              {
                "alias": null,
                "args": (v7/*: any*/),
                "concreteType": "Deployment",
                "kind": "LinkedField",
                "name": "deployments",
                "plural": true,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Pod",
                    "kind": "LinkedField",
                    "name": "pods",
                    "plural": true,
                    "selections": [
                      {
                        "condition": "showPods",
                        "kind": "Condition",
                        "passingValue": true,
                        "selections": [
                          (v8/*: any*/),
                          (v9/*: any*/)
                        ]
                      }
                    ],
                    "storageKey": null
                  },
                  {
                    "condition": "showDeployments",
                    "kind": "Condition",
                    "passingValue": true,
                    "selections": (v10/*: any*/)
                  }
                ],
                "storageKey": null
              },
              {
                "condition": "showNamespaces",
                "kind": "Condition",
                "passingValue": true,
                "selections": (v10/*: any*/)
              }
            ],
            "storageKey": null
          },
          {
            "condition": "showClusters",
            "kind": "Condition",
            "passingValue": true,
            "selections": (v10/*: any*/)
          }
        ],
        "storageKey": null
      }
    ],
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v1/*: any*/),
      (v0/*: any*/),
      (v2/*: any*/),
      (v4/*: any*/),
      (v3/*: any*/),
      (v5/*: any*/)
    ],
    "kind": "Operation",
    "name": "View2Query",
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Cluster",
        "kind": "LinkedField",
        "name": "cluster",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": (v6/*: any*/),
            "concreteType": "Namespace",
            "kind": "LinkedField",
            "name": "namespaces",
            "plural": true,
            "selections": [
              {
                "alias": null,
                "args": (v7/*: any*/),
                "concreteType": "Deployment",
                "kind": "LinkedField",
                "name": "deployments",
                "plural": true,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Pod",
                    "kind": "LinkedField",
                    "name": "pods",
                    "plural": true,
                    "selections": [
                      (v8/*: any*/),
                      {
                        "condition": "showPods",
                        "kind": "Condition",
                        "passingValue": true,
                        "selections": [
                          (v9/*: any*/)
                        ]
                      }
                    ],
                    "storageKey": null
                  },
                  (v8/*: any*/)
                ],
                "storageKey": null
              },
              (v8/*: any*/)
            ],
            "storageKey": null
          },
          (v8/*: any*/)
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "531f230ac57d6ba697b2cdc90c9baffb",
    "id": null,
    "metadata": {},
    "name": "View2Query",
    "operationKind": "query",
    "text": "query View2Query(\n  $namespaceID: ID\n  $deploymentID: ID\n  $showClusters: Boolean!\n  $showNamespaces: Boolean!\n  $showDeployments: Boolean!\n  $showPods: Boolean!\n) {\n  cluster {\n    id @include(if: $showClusters)\n    namespaces(namespaceid: $namespaceID) {\n      id @include(if: $showNamespaces)\n      deployments(deploymentid: $deploymentID) {\n        id @include(if: $showDeployments)\n        pods {\n          id @include(if: $showPods)\n          name @include(if: $showPods)\n          id\n        }\n        id\n      }\n      id\n    }\n    id\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '27430c38a2578255abc6044fb0ec1eee';

module.exports = node;

/**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MeshType = "ALL_MESH" | "APP_MESH" | "CITRIX_SERVICE_MESH" | "CONSUL" | "INVALID_MESH" | "ISTIO" | "KUMA" | "LINKERD" | "NETWORK_SERVICE_MESH" | "NGINX_SERVICE_MESH" | "OCTARINE" | "OPEN_SERVICE_MESH" | "TANZU" | "TRAEFIK_MESH" | "%future added value";
export type ControlPlaneFilter = {|
  type?: ?MeshType
|};
export type ControlPlaneSubscriptionVariables = {|
  filter?: ?ControlPlaneFilter
|};
export type ControlPlaneSubscriptionResponse = {|
  +controlPlanesState: $ReadOnlyArray<{|
    +name: string,
    +members: $ReadOnlyArray<{|
      +name: string,
      +version: string,
      +component: string,
      +namespace: string,
      +data_planes: ?$ReadOnlyArray<?{|
        +name: string,
        +image: string,
        +ports: ?$ReadOnlyArray<?{|
          +name: ?string,
          +containerPort: number,
          +protocol: string,
        |}>,
        +resources: ?{|
          +limits: ?{|
            +cpu: ?string,
            +memory: ?string,
          |},
          +requests: ?{|
            +cpu: ?string,
            +memory: ?string,
          |},
        |},
      |}>,
    |}>,
  |}>
|};
export type ControlPlaneSubscription = {|
  variables: ControlPlaneSubscriptionVariables,
  response: ControlPlaneSubscriptionResponse,
|};
*/


/*
subscription ControlPlaneSubscription(
  $filter: ControlPlaneFilter
) {
  controlPlanesState: listenToControlPlaneState(filter: $filter) {
    name
    members {
      name
      version
      component
      namespace
      data_planes {
        name
        image
        ports {
          name
          containerPort
          protocol
        }
        resources {
          limits {
            cpu
            memory
          }
          requests {
            cpu
            memory
          }
        }
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "filter"
  }
],
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "name",
  "storageKey": null
},
v2 = [
  {
    "alias": null,
    "args": null,
    "kind": "ScalarField",
    "name": "cpu",
    "storageKey": null
  },
  {
    "alias": null,
    "args": null,
    "kind": "ScalarField",
    "name": "memory",
    "storageKey": null
  }
],
v3 = [
  {
    "alias": "controlPlanesState",
    "args": [
      {
        "kind": "Variable",
        "name": "filter",
        "variableName": "filter"
      }
    ],
    "concreteType": "ControlPlane",
    "kind": "LinkedField",
    "name": "listenToControlPlaneState",
    "plural": true,
    "selections": [
      (v1/*: any*/),
      {
        "alias": null,
        "args": null,
        "concreteType": "ControlPlaneMember",
        "kind": "LinkedField",
        "name": "members",
        "plural": true,
        "selections": [
          (v1/*: any*/),
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
            "name": "component",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "namespace",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "concreteType": "Container",
            "kind": "LinkedField",
            "name": "data_planes",
            "plural": true,
            "selections": [
              (v1/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "image",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "concreteType": "Container_Port",
                "kind": "LinkedField",
                "name": "ports",
                "plural": true,
                "selections": [
                  (v1/*: any*/),
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "containerPort",
                    "storageKey": null
                  },
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "protocol",
                    "storageKey": null
                  }
                ],
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "concreteType": "Container_Resource",
                "kind": "LinkedField",
                "name": "resources",
                "plural": false,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Resource",
                    "kind": "LinkedField",
                    "name": "limits",
                    "plural": false,
                    "selections": (v2/*: any*/),
                    "storageKey": null
                  },
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Resource",
                    "kind": "LinkedField",
                    "name": "requests",
                    "plural": false,
                    "selections": (v2/*: any*/),
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
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "ControlPlaneSubscription",
    "selections": (v3/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ControlPlaneSubscription",
    "selections": (v3/*: any*/)
  },
  "params": {
    "cacheID": "f33b684c926f5b6afbafe460140230ef",
    "id": null,
    "metadata": {},
    "name": "ControlPlaneSubscription",
    "operationKind": "subscription",
    "text": "subscription ControlPlaneSubscription(\n  $filter: ControlPlaneFilter\n) {\n  controlPlanesState: listenToControlPlaneState(filter: $filter) {\n    name\n    members {\n      name\n      version\n      component\n      namespace\n      data_planes {\n        name\n        image\n        ports {\n          name\n          containerPort\n          protocol\n        }\n        resources {\n          limits {\n            cpu\n            memory\n          }\n          requests {\n            cpu\n            memory\n          }\n        }\n      }\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '9db33629b704e6233c09375d7c89f7fe';

module.exports = node;

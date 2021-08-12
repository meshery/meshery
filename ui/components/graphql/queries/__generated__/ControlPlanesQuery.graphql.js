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
export type ControlPlanesQueryVariables = {|
  filter?: ?ControlPlaneFilter
|};
export type ControlPlanesQueryResponse = {|
  +controlPlanesState: $ReadOnlyArray<{|
    +name: string,
    +members: $ReadOnlyArray<{|
      +name: string,
      +version: string,
      +component: string,
      +namespace: string,
      +data_planes: ?$ReadOnlyArray<{|
        +name: string,
        +image: string,
        +status: ?{|
          +ready: boolean,
          +started: boolean,
          +state: ?any,
        |},
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
export type ControlPlanesQuery = {|
  variables: ControlPlanesQueryVariables,
  response: ControlPlanesQueryResponse,
|};
*/


/*
query ControlPlanesQuery(
  $filter: ControlPlaneFilter
) {
  controlPlanesState: getControlPlanes(filter: $filter) {
    name
    members {
      name
      version
      component
      namespace
      data_planes {
        name
        image
        status {
          ready
          started
          state
        }
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
    "name": "getControlPlanes",
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
                "concreteType": "Container_Status",
                "kind": "LinkedField",
                "name": "status",
                "plural": false,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "ready",
                    "storageKey": null
                  },
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "started",
                    "storageKey": null
                  },
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "state",
                    "storageKey": null
                  }
                ],
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
    "name": "ControlPlanesQuery",
    "selections": (v3/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "ControlPlanesQuery",
    "selections": (v3/*: any*/)
  },
  "params": {
    "cacheID": "0cd0de0e783c2d752355d1d82b05e238",
    "id": null,
    "metadata": {},
    "name": "ControlPlanesQuery",
    "operationKind": "query",
    "text": "query ControlPlanesQuery(\n  $filter: ControlPlaneFilter\n) {\n  controlPlanesState: getControlPlanes(filter: $filter) {\n    name\n    members {\n      name\n      version\n      component\n      namespace\n      data_planes {\n        name\n        image\n        status {\n          ready\n          started\n          state\n        }\n        ports {\n          name\n          containerPort\n          protocol\n        }\n        resources {\n          limits {\n            cpu\n            memory\n          }\n          requests {\n            cpu\n            memory\n          }\n        }\n      }\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '6638721aef7c56a632996157605013cb';

module.exports = node;

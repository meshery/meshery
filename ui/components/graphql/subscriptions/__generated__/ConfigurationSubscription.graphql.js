/**
 * @generated SignedSource<<393dc75182410b818fda16b7bb4f5a56>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

var node = (function(){
var v0 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "applicationSelector"
},
v1 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "filterSelector"
},
v2 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "patternSelector"
},
v3 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "page",
  "storageKey": null
},
v4 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "page_size",
  "storageKey": null
},
v5 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "total_count",
  "storageKey": null
},
v6 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v7 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "name",
  "storageKey": null
},
v8 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "user_id",
  "storageKey": null
},
v9 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "visibility",
  "storageKey": null
},
v10 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "created_at",
  "storageKey": null
},
v11 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "updated_at",
  "storageKey": null
},
v12 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "catalog_data",
  "storageKey": null
},
v13 = [
  {
    "alias": "configuration",
    "args": [
      {
        "kind": "Variable",
        "name": "applicationSelector",
        "variableName": "applicationSelector"
      },
      {
        "kind": "Variable",
        "name": "filterSelector",
        "variableName": "filterSelector"
      },
      {
        "kind": "Variable",
        "name": "patternSelector",
        "variableName": "patternSelector"
      }
    ],
    "concreteType": "ConfigurationPage",
    "kind": "LinkedField",
    "name": "subscribeConfiguration",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "ApplicationPage",
        "kind": "LinkedField",
        "name": "applications",
        "plural": false,
        "selections": [
          (v3/*: any*/),
          (v4/*: any*/),
          (v5/*: any*/),
          {
            "alias": null,
            "args": null,
            "concreteType": "ApplicationResult",
            "kind": "LinkedField",
            "name": "applications",
            "plural": true,
            "selections": [
              (v6/*: any*/),
              (v7/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "application_file",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "concreteType": "NullString",
                "kind": "LinkedField",
                "name": "type",
                "plural": false,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "String",
                    "storageKey": null
                  },
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "Valid",
                    "storageKey": null
                  }
                ],
                "storageKey": null
              },
              (v8/*: any*/),
              (v9/*: any*/),
              (v10/*: any*/),
              (v11/*: any*/)
            ],
            "storageKey": null
          }
        ],
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "PatternPageResult",
        "kind": "LinkedField",
        "name": "patterns",
        "plural": false,
        "selections": [
          (v3/*: any*/),
          (v4/*: any*/),
          (v5/*: any*/),
          {
            "alias": null,
            "args": null,
            "concreteType": "PatternResult",
            "kind": "LinkedField",
            "name": "patterns",
            "plural": true,
            "selections": [
              (v6/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "pattern_file",
                "storageKey": null
              },
              (v9/*: any*/),
              (v12/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "canSupport",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "errmsg",
                "storageKey": null
              },
              (v10/*: any*/),
              (v11/*: any*/)
            ],
            "storageKey": null
          }
        ],
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "FilterPage",
        "kind": "LinkedField",
        "name": "filters",
        "plural": false,
        "selections": [
          (v3/*: any*/),
          (v4/*: any*/),
          (v5/*: any*/),
          {
            "alias": null,
            "args": null,
            "concreteType": "FilterResult",
            "kind": "LinkedField",
            "name": "filters",
            "plural": true,
            "selections": [
              (v6/*: any*/),
              (v7/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "filter_file",
                "storageKey": null
              },
              (v9/*: any*/),
              (v12/*: any*/),
              (v8/*: any*/),
              (v10/*: any*/),
              (v11/*: any*/)
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
    "argumentDefinitions": [
      (v0/*: any*/),
      (v1/*: any*/),
      (v2/*: any*/)
    ],
    "kind": "Fragment",
    "metadata": null,
    "name": "ConfigurationSubscription",
    "selections": (v13/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v0/*: any*/),
      (v2/*: any*/),
      (v1/*: any*/)
    ],
    "kind": "Operation",
    "name": "ConfigurationSubscription",
    "selections": (v13/*: any*/)
  },
  "params": {
    "cacheID": "f306c1d90559367e4a4c22d2a59f40b7",
    "id": null,
    "metadata": {},
    "name": "ConfigurationSubscription",
    "operationKind": "subscription",
    "text": "subscription ConfigurationSubscription(\n  $applicationSelector: PageFilter!\n  $patternSelector: PageFilter!\n  $filterSelector: PageFilter!\n) {\n  configuration: subscribeConfiguration(applicationSelector: $applicationSelector, patternSelector: $patternSelector, filterSelector: $filterSelector) {\n    applications {\n      page\n      page_size\n      total_count\n      applications {\n        id\n        name\n        application_file\n        type {\n          String\n          Valid\n        }\n        user_id\n        visibility\n        created_at\n        updated_at\n      }\n    }\n    patterns {\n      page\n      page_size\n      total_count\n      patterns {\n        id\n        name\n        user_id\n        pattern_file\n        visibility\n        catalog_data\n        canSupport\n        errmsg\n        created_at\n        updated_at\n      }\n    }\n    filters {\n      page\n      page_size\n      total_count\n      filters {\n        id\n        name\n        filter_file\n        visibility\n        catalog_data\n        user_id\n        created_at\n        updated_at\n      }\n    }\n  }\n}\n"
  }
};
})();

node.hash = "e253f7fac4ca190bc348f864eec5525b";

module.exports = node;

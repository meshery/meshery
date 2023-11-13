/**
 * @generated SignedSource<<d4cdc76aba325f0e522811cde51a1da3>>
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

var node = (function(){
var v0 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "filterSelector"
},
v1 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "patternSelector"
},
v2 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "page",
  "storageKey": null
},
v3 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "page_size",
  "storageKey": null
},
v4 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "total_count",
  "storageKey": null
},
v5 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v6 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "name",
  "storageKey": null
},
v7 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "user_id",
  "storageKey": null
},
v8 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "visibility",
  "storageKey": null
},
v9 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "catalog_data",
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
v12 = [
  {
    "alias": "configuration",
    "args": [
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
        "concreteType": "PatternPageResult",
        "kind": "LinkedField",
        "name": "patterns",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
          {
            "alias": null,
            "args": null,
            "concreteType": "PatternResult",
            "kind": "LinkedField",
            "name": "patterns",
            "plural": true,
            "selections": [
              (v5/*: any*/),
              (v6/*: any*/),
              (v7/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "pattern_file",
                "storageKey": null
              },
              (v8/*: any*/),
              (v9/*: any*/),
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
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
          {
            "alias": null,
            "args": null,
            "concreteType": "FilterResult",
            "kind": "LinkedField",
            "name": "filters",
            "plural": true,
            "selections": [
              (v5/*: any*/),
              (v6/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "filter_file",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "filter_resource",
                "storageKey": null
              },
              (v8/*: any*/),
              (v9/*: any*/),
              (v7/*: any*/),
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
      (v1/*: any*/)
    ],
    "kind": "Fragment",
    "metadata": null,
    "name": "ConfigurationSubscription",
    "selections": (v12/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v1/*: any*/),
      (v0/*: any*/)
    ],
    "kind": "Operation",
    "name": "ConfigurationSubscription",
    "selections": (v12/*: any*/)
  },
  "params": {
    "cacheID": "8d81d5dd33c4ea1b76228108d5cf7f2a",
    "id": null,
    "metadata": {},
    "name": "ConfigurationSubscription",
    "operationKind": "subscription",
    "text": "subscription ConfigurationSubscription(\n  $patternSelector: PageFilter!\n  $filterSelector: PageFilter!\n) {\n  configuration: subscribeConfiguration(patternSelector: $patternSelector, filterSelector: $filterSelector) {\n    patterns {\n      page\n      page_size\n      total_count\n      patterns {\n        id\n        name\n        user_id\n        pattern_file\n        visibility\n        catalog_data\n        canSupport\n        errmsg\n        created_at\n        updated_at\n      }\n    }\n    filters {\n      page\n      page_size\n      total_count\n      filters {\n        id\n        name\n        filter_file\n        filter_resource\n        visibility\n        catalog_data\n        user_id\n        created_at\n        updated_at\n      }\n    }\n  }\n}\n"
  }
};
})();

node.hash = "28ee13b2d3a5254fe8c7933ab7835210";

module.exports = node;

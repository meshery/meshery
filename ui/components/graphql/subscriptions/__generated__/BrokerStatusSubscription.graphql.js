/**
 * @generated SignedSource<<b5bafefd8ea424078209f23685fb6ae5>>
 * @flow
 * @lightSyntaxTransform
 * @nogrep
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest, GraphQLSubscription } from 'relay-runtime';
export type BrokerStatusSubscription$variables = {||};
export type BrokerStatusSubscription$data = {|
  +subscribeBrokerConnection: boolean,
|};
export type BrokerStatusSubscription = {|
  variables: BrokerStatusSubscription$variables,
  response: BrokerStatusSubscription$data,
|};
*/

var node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "alias": null,
    "args": null,
    "kind": "ScalarField",
    "name": "subscribeBrokerConnection",
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "BrokerStatusSubscription",
    "selections": (v0/*: any*/),
    "type": "Subscription",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "BrokerStatusSubscription",
    "selections": (v0/*: any*/)
  },
  "params": {
    "cacheID": "3eaadbc18fa2fc8430e387a3059fafa9",
    "id": null,
    "metadata": {},
    "name": "BrokerStatusSubscription",
    "operationKind": "subscription",
    "text": "subscription BrokerStatusSubscription {\n  subscribeBrokerConnection\n}\n"
  }
};
})();

(node/*: any*/).hash = "c3d1793f9ca896edd7ce6ec58cb79a59";

module.exports = ((node/*: any*/)/*: GraphQLSubscription<
  BrokerStatusSubscription$variables,
  BrokerStatusSubscription$data,
>*/);

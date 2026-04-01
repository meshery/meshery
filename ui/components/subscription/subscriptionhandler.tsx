import { fnMapping } from './helpers';

/**
 * A GraphQL subscription controller
 */
export class GQLSubscription {
  subscription;
  state;
  callbackFunction;
  connectionIDs;
  type;

  constructor({ type, connectionIDs, callbackFunction }) {
    this.type = type;
    this.connectionIDs = connectionIDs || [];
    this.callbackFunction = callbackFunction;
  }

  setState = (newState) => {
    this.state = newState;
  };

  initSubscription = () => {
    if (!this.type) {
      throw new Error('Subscription Type is Empty, initialise a subscription in a constructor');
    }

    this.updateSubscription(this.connectionIDs);
  };

  doCallback = (data) => {
    // assuming the data is in stream and contextId being the key to the data
    data = data[fnMapping[this.type].eventName];
    if (fnMapping[this.type]?.comparatorFn(this.state, data)) {
      // returns true if data is updated
      this.setState(fnMapping[this.type].mergeFn(this.state, data));

      this.callbackFunction(this.state);
    }
  };

  flushSubscription = () => {
    if (this.subscription) {
      this.subscription.dispose();
    }
  };

  updateSubscription = (connectionIDs) => {
    this.flushSubscription();
    this.subscription = fnMapping[this.type].subscriptionFn(this.doCallback, connectionIDs);
  };
}

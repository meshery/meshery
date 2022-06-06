import _ from "lodash"

/**
 * A GraphQL subscription controller
 */
export class GQLSubscription {
  subscriptionFn;
  subscription;
  state;
  callbackFunctions;
  contextIds;

  constructor({ subscriptionFn, contextIds }) {
    this.subscriptionFn = subscriptionFn;
    this.contextIds = contextIds || [];
    this.callbackFunctions = [];
  }

  setState = (newState) => {
    state = newState;
  }

  initSubscription = () => {
    if (!this.subscriptionFn)  {
      throw new Error("Subscription function is Empty, initialise a subscription in a constructor")
    }

    // if already done
    this.updateSubscription(this.contextIds)
  }

  addCallback = (fn, caller) => {
    this.callbackFunctions.push({ fn, caller });
  }

  removeCallback = (caller)  => {
    this.callbackFunctions.filter(fn => fn.caller !== caller);
  }

  doCallback = (data) => {
    if (!_.isEqual(this.state, data)) {
      this.setState(data);
      this.callbackFunctions.forEach(fn => fn.fn(data));
    }
  }

  flushSubscription = () => {
    if (this.subscription) {
      this.subscription.destroy();
    }
  }

  updateSubscription = (contextIds) => {
    this.flushSubscription();
    this.subscription = this.subscriptionFn(this.doCallback, contextIds);
  }
}

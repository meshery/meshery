import _ from 'lodash';

export function isMeshSyncSubscriptionDataUpdated(currentState, newData) {
  if (!currentState) {
    return true;
  }

  const oldData = currentState.find((data) => data.contextID === newData.contextID);
  if (!oldData) {
    return true;
  }

  return !_.isEqual(oldData, newData);
}

export function isOperatorStateSubscriptionDataUpdated(currentState, newData) {
  if (!currentState) {
    return true;
  }

  const oldData = currentState.find((data) => data.contextID === newData.contextID);
  if (!oldData) {
    return true;
  }

  return !_.isEqual(oldData, newData);
}

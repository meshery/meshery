import _ from 'lodash';
import { isControllerObjectEqual } from './helpers';

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

export function isMesheryControllerStateSubscriptionDataUpdated(currentState, newData) {
  if (!currentState) {
    return true;
  }

  const oldData = currentState.filter((data) => {
    for (let newControllerData of newData) {
      if (isControllerObjectEqual(newControllerData, data)) {
        return true;
      }
    }
    return false;
  });
  if (!oldData) {
    return true;
  }
  return !_.isEqual(oldData, newData);
}

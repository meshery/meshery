import { isControllerObjectEqual } from './helpers';

export function mergeMeshSyncSubscription(currentState, newData) {
  if (!currentState) {
    return [newData];
  }

  return [...currentState.filter((data) => data.contextID !== newData.contextID), newData];
}

export function mergeOperatorStateSubscription(currentState, newData) {
  // the merger function for both are same
  return mergeMeshSyncSubscription(currentState, newData);
}

export function mergeMesheryController(currentState, newData) {
  if (!currentState) {
    return newData;
  }
  const merged = [
    ...currentState.filter((data) => {
      for (let newControllerData of newData) {
        if (isControllerObjectEqual(newControllerData, data)) {
          return false;
        }
      }
      return true;
    }),
    ...newData,
  ];
  return merged;
}

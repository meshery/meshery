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
  console.log('CONTROLLER BEFORE MERGE: ', currentState, newData);
  const merged = [
    ...currentState.filter((data) => {
      for (let newControllerData of newData) {
        if (isControllerObjectEqual(newControllerData, data)) {
          console.log(
            'CONTROLLER INSIDE MERGE LOOP: ',
            newControllerData.controller,
            data.controller,
            newControllerData?.status,
            data?.status,
          );
          return false;
        }
      }
      return true;
    }),
    ...newData,
  ];
  console.log('CONTROLLER MERGED FUNCTION EXIT: ', merged);
  return merged;
}

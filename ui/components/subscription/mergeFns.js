export function mergeMeshSyncSubscription(currentState, newData) {
  if (!currentState) {
    return [newData];
  }

  return [...currentState.filter(data => data.contextID !== newData.contextID), newData];
}

export function mergeOperatorStateSubscription(currentState, newData) {
  // the merger function for both are same
  return mergeMeshSyncSubscription(currentState, newData);
}
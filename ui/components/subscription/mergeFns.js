export function mergeMeshSyncSubscription(currentState, newData) {
  if (!currentState) {
    return [newData];
  }

  return [...currentState.filter(data => data.contextID !== newData.contextID), newData];
}

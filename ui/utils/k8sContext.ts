type K8sContextIdentity = {
  id?: string;
  connectionId?: string;
};

/**
 * Returns true when a kubernetes context row corresponds to the given connection id.
 * Context ids are hex strings; connection ids are canonical UUID strings.
 */
export function k8sContextMatchesConnectionId(
  ctx: K8sContextIdentity,
  connectionId: string,
): boolean {
  if (!ctx || !connectionId) {
    return false;
  }

  const normalizedConnectionId = connectionId.replace(/-/g, '');

  return (
    ctx.connectionId === connectionId ||
    ctx.id === connectionId ||
    ctx.id === normalizedConnectionId
  );
}

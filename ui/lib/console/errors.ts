/**
 * Turning a failed console request into something a human can act on.
 *
 * The schemas base query parses the server's MeshKit error envelope onto
 * `error.meshkit`, so the short description, the code, and the suggested
 * remediation are all available. Rendering "something went wrong" over the top
 * of that throws away everything the server took the trouble to say.
 */

/** The MeshKit envelope the schemas base query attaches to a failed query. */
interface MeshkitErrorLike {
  message?: string;
  code?: string;
  probableCause?: string[];
  suggestedRemediation?: string[];
}

interface QueryErrorLike {
  status?: number | string;
  meshkit?: MeshkitErrorLike;
  error?: string;
}

export interface ConsoleErrorDescription {
  message: string;
  /** Extra context, e.g. the server's suggested remediation. */
  detail?: string;
  /** MeshKit code, shown so a user can quote it in a bug report. */
  code?: string;
}

/**
 * Describes why resolving a target's capabilities failed.
 *
 * A 404 gets a purpose-written message rather than the server's: the resource is
 * addressed from a MeshSync-cached view, so "not found" almost always means the
 * view is showing something the cluster has since deleted — which is a very
 * different thing from a typo, and worth saying out loud.
 */
export const describeConsoleError = (
  error: unknown,
  targetName?: string,
): ConsoleErrorDescription => {
  const queryError = (error ?? {}) as QueryErrorLike;
  const meshkit = queryError.meshkit;

  if (queryError.status === 404) {
    const subject = targetName ? `"${targetName}"` : 'This resource';
    return {
      message: `${subject} no longer exists in the cluster.`,
      detail:
        'It is still listed here because this view is served from MeshSync’s cache, which has not caught up yet.',
      code: meshkit?.code,
    };
  }

  if (queryError.status === 'FETCH_ERROR') {
    return { message: 'Could not reach the Meshery server.' };
  }

  if (meshkit?.message) {
    return {
      message: meshkit.message,
      detail: meshkit.suggestedRemediation?.[0] ?? meshkit.probableCause?.[0],
      code: meshkit.code,
    };
  }

  return { message: 'Could not determine what this resource supports.' };
};

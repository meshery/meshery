import { useCallback } from 'react';
import { useNotificationHandlers } from './useNotification';

interface NotifiedMutationOptions {
  successMessage?: string | ((result: unknown) => string);
  errorMessage?: string | ((error: unknown) => string);
}

/**
 * Wraps an RTK Query mutation trigger with automatic success/error notifications.
 * Eliminates the repetitive `.then(notify).catch(notify)` pattern.
 *
 * Usage:
 * ```
 * const [triggerDelete] = useDeleteDesignMutation();
 * const deleteDesign = useNotifiedMutation(triggerDelete, {
 *   successMessage: 'Design deleted',
 *   errorMessage: 'Failed to delete design',
 * });
 *
 * // Later:
 * await deleteDesign({ id: '123' });
 * ```
 */
export function useNotifiedMutation<TArg, TResult>(
  trigger: (arg: TArg) => { unwrap: () => Promise<TResult> },
  options: NotifiedMutationOptions = {},
) {
  const { handleSuccess, handleError } = useNotificationHandlers();

  return useCallback(
    async (arg: TArg) => {
      try {
        const result = await trigger(arg).unwrap();
        if (options.successMessage) {
          const msg =
            typeof options.successMessage === 'function'
              ? options.successMessage(result)
              : options.successMessage;
          handleSuccess(msg);
        }
        return result;
      } catch (error) {
        const msg =
          typeof options.errorMessage === 'function'
            ? options.errorMessage(error)
            : options.errorMessage || 'An error occurred';
        handleError(msg);
        throw error;
      }
    },
    [trigger, options.successMessage, options.errorMessage, handleSuccess, handleError],
  );
}

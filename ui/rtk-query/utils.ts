import { store } from '../store';

/**
 * Builds the callback form of an `enhanceEndpoints` partial that appends local
 * cache tags to a schemas-generated endpoint. The callback is handed the live
 * generated definition, so the tags schemas declares are kept; the object form
 * of `enhanceEndpoints` would `Object.assign` over them.
 * @param {string} operationId - The generated operationId being enhanced, used in the failure message.
 * @param {...(string|Object)} localTags - Tag descriptions to invalidate in addition to the generated ones.
 * @returns {(definition: Object) => void} - An `enhanceEndpoints` endpoint callback.
 * @throws {Error} If the operationId is absent from the generated client.
 */
export const appendInvalidatesTags =
  (operationId, ...localTags) =>
  (definition) => {
    if (!definition) {
      throw new Error(
        `RTK Query endpoint "${operationId}" is not present on the @meshery/schemas generated client, ` +
          `so its local cache tags could not be attached. @meshery/schemas most likely renamed or removed ` +
          `the operation - update the enhanceEndpoints call in ui/rtk-query to the current operationId.`,
      );
    }

    const generatedTags = definition.invalidatesTags;
    definition.invalidatesTags = (result, error, arg, meta) => [
      ...(typeof generatedTags === 'function'
        ? generatedTags(result, error, arg, meta)
        : (generatedTags ?? [])),
      ...localTags,
    ];
  };

/**
 * Initiates a query using specified query and variables via store.dispatch.
 * @param {Object} query - The query object containing the initiate function.
 * @param {any} variables - Variables to be passed to the query initiate function.
 * @returns {Promise<Object>} - A Promise resolving with an object containing query execution results.
 */
export const initiateQuery = async (query, variables) => {
  try {
    return await store.dispatch(query.initiate(variables));
  } catch (error) {
    // Return an object with error details if there's an exception
    return {
      data: null,
      error: error,
      isFetching: false,
      isSuccess: false,
      isLoading: false,
      isError: true,
    };
  }
};

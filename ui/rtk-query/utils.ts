import { store } from '../store';

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

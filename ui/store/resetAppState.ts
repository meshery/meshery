import { api } from '../rtk-query';
import { subscriptionClient } from '../lib/relayEnvironment';
import Cookies from 'universal-cookie';

/**
 * Resets all client-side state on logout.
 * - Resets RTK Query cache
 * - Disposes all GraphQL subscriptions via the shared WS client
 * - Clears cookies and session storage
 */
export const resetAppState = () => (dispatch) => {
  // Reset RTK Query cache
  dispatch(api.util.resetApiState());

  // Terminate all active GraphQL subscriptions
  if (subscriptionClient) {
    subscriptionClient.terminate();
  }

  // Clear session storage
  sessionStorage.removeItem('currentOrg');
  sessionStorage.removeItem('currentWorkspace');
  sessionStorage.removeItem('keys');

  // Clear cookies
  const cookies = new Cookies();
  cookies.remove('registered', { path: '/' });
};

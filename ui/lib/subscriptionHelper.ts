import { requestSubscription } from 'react-relay';
import { createRelayEnvironment } from './relayEnvironment';

type OnNext = (data: unknown) => void;
type OnError = (error: Error) => void;

interface SubscriptionOptions {
  subscription: unknown;
  variables?: Record<string, unknown>;
  onNext: OnNext;
  onError?: OnError;
  subscriptionName?: string;
}

/**
 * Creates a Relay subscription with consistent error handling.
 * Uses the singleton Relay environment and logs errors with the subscription name.
 * Callers can provide a custom onError to surface errors to the UI.
 */
export function createSubscription({
  subscription,
  variables = {},
  onNext,
  onError,
  subscriptionName = 'Unknown',
}: SubscriptionOptions) {
  const environment = createRelayEnvironment({});

  const defaultOnError: OnError = (error) => {
    console.error(`[GraphQL Subscription: ${subscriptionName}] Error:`, error);
  };

  return requestSubscription(environment, {
    subscription,
    variables,
    onNext,
    onError: onError || defaultOnError,
  });
}

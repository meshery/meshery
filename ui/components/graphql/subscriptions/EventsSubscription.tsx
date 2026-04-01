import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

const eventsSubscription = graphql`
  subscription EventsSubscription {
    event: subscribeEvents {
      id
      userID
      actedUpon
      operationID
      systemID
      status
      severity
      action
      category
      description
      metadata
      createdAt
      updatedAt
      deletedAt
    }
  }
`;

export default function subscribeEvents(dataCB, onError) {
  const environment = createRelayEnvironment({});

  return requestSubscription(environment, {
    subscription: eventsSubscription,
    onNext: dataCB,
    onError:
      onError || ((error) => console.error('An error occurred in subscription to events:', error)),
  });
}

import { graphql } from 'react-relay';
import { createSubscription } from '../../../lib/subscriptionHelper';

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

export default function subscribeEvents(dataCB, onError?) {
  return createSubscription({
    subscription: eventsSubscription,
    onNext: dataCB,
    onError,
    subscriptionName: 'Events',
  });
}

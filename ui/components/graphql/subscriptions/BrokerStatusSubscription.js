import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

const brokerStatusSubscription = graphql`
  subscription BrokerStatusSubscription {
    subscribeBrokerConnection
  }
`;

export default function subscribeBrokerStatusEvents(dataCB) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription: brokerStatusSubscription,
    variables: {},
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}

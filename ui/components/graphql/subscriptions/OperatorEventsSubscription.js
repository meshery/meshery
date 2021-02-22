import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const operatorEventsSubscription = graphql`
  subscription OperatorEventsSubscription {
    listenToOperatorEvents {
      status
      error {
        code
        description
      }
    }
  }
`;

export default function subscribeOperatorEvents(dataCB) {
  requestSubscription(environment, {
    subscription: operatorEventsSubscription,
    variables: {},
    onNext: dataCB,
    updater: updater,
    onError: (error) => console.log(`An error occured:`, error),
  });
}

import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const operatorStatusSubscription = graphql`
  subscription OperatorStatusSubscription {
    operator: listenToOperatorState {
      status
      version
      controllers {
          name
          version
          status
      }
      error {
        code
        description
      }
    }
  }
`;

export default function subscribeOperatorStatusEvents(dataCB) {
  requestSubscription(environment, {
    subscription: operatorStatusSubscription,
    variables: {},
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}

import environment from "@/app/relayEnvironment";
import { graphql, requestSubscription } from "react-relay";

const operatorStatusSubscription = graphql`
  subscription OperatorStatusEventsSubscription {
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

export default function subscribeOperatorStatusEvents(dataCB, errorCB) {
  requestSubscription(environment, {
    subscription: operatorStatusSubscription,
    variables: {},
    onNext: dataCB,
    onError: errorCB,
  });
}

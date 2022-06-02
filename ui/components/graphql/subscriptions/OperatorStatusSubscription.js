import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const operatorStatusSubscription = graphql`
  subscription OperatorStatusSubscription($k8scontextID: String!) {
    operator: listenToOperatorState(k8scontextID: $k8scontextID) {
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
  return  requestSubscription(environment, {
    subscription : operatorStatusSubscription,
    variables : {},
    onNext : dataCB,
    onError : (error) => console.log(`An error occured:`, error),
  });
}

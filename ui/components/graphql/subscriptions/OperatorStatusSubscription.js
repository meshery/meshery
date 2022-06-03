import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const operatorStatusSubscription = graphql`
  subscription OperatorStatusSubscription($k8scontextIDs: [String!]) {
    operator: listenToOperatorState(k8scontextIDs: $k8scontextIDs) {
      contextID
      operatorStatus {
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
  }
`;

export default function subscribeOperatorStatusEvents(dataCB, contextIds) {
  return  requestSubscription(environment, {
    subscription : operatorStatusSubscription,
    variables : { k8scontextIDs : contextIds },
    onNext : dataCB,
    onError : (error) => console.log(`An error occured:`, error),
  });
}

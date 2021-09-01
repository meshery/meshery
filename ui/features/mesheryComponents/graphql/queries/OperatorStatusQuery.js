import { graphql, fetchQuery } from "react-relay";
import environment from "../../../../app/relayEnvironment";

export default function fetchMesheryOperatorStatus() {
  const query = graphql`
    query OperatorStatusQuery {
      operator: getOperatorStatus {
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

  return fetchQuery(environment, query);
}

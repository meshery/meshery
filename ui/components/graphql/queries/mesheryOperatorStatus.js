import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchMesheryOperatorStatus() {
  const query = graphql`
        query mesheryOperatorStatusQuery{
          operator: getOperatorStatus {
            status
            error {
              code
              description
            }
          }
        }
    `;

  return fetchQuery(environment, query)
}

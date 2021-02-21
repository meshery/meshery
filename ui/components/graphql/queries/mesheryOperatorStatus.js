import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchMesheryOperatorStatus() {
  const query = graphql`
        query mesheryOperatorStatusQuery{
          getOperatorStatus
        }
    `;

  return fetchQuery(environment, query)
}

import { graphql, fetchQuery } from "react-relay";
import environment from "../environment";

export function fetchMesheryOperatorStatus() {
  const query = graphql`
        query mesheryOperatorStatusQuery{
          getOperatorStatus
        }
    `;

  return fetchQuery(environment, query)
}

import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchMesheryOperatorStatus(variables) {
  const vars = { k8scontextID : variables.k8scontextID };

  const query = graphql`
        query OperatorStatusQuery($k8scontextID: String!) {
          operator: getOperatorStatus(k8scontextID: $k8scontextID) {
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

  return fetchQuery(environment, query, vars)
}

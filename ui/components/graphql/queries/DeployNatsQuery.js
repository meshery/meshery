import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function connectToNATS() {

  const query = graphql`
    query DeployNatsQuery($k8scontextID: String!)  {
      connectToNats (k8scontextID: $k8scontextID)
    }
  `;

  return fetchQuery(environment, query);
}

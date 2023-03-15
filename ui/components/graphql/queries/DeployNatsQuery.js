import { fetchQuery, graphql } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

export default function connectToNATS() {
  const environment = createRelayEnvironment({});

  const DeployNatsQueryNode = graphql`
    query DeployNatsQuery($k8scontextID: String!)  {
      connectToNats (k8scontextID: $k8scontextID)
    }
  `;

  return fetchQuery(environment, DeployNatsQueryNode);
}

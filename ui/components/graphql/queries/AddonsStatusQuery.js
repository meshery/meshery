import { fetchQuery, graphql } from "relay-runtime";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

export default function fetchAvailableAddons(variables) {
  const environment = createRelayEnvironment({});
  const vars = { filter : variables }

  const AddonsStatusQueryNode = graphql`
  query AddonsStatusQuery($filter: ServiceMeshFilter) {
    addonsState: getAvailableAddons(filter: $filter) {
      name
      owner
      }
  }
`
  return fetchQuery(environment, AddonsStatusQueryNode, vars);
}
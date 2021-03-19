import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchAvailableNamespaces() {
  const query = graphql`
        query NamespaceQuery {
          namespaces: getAvailableNamespaces {
            namespace
          }
        }
    `;

  return fetchQuery(environment, query)
}

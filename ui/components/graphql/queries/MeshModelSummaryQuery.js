import { graphql ,fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchMeshModelSummary(selector) {
  const vars = {
    selector : selector
  }

  const query = graphql`
    query MeshModelSummaryQuery($selector: MeshModelSummarySelector!) {
      meshmodelSummary: getMeshModelSummary(selector: $selector) {
        components {
          name,
          count
        },
        relationships {
          name,
          count
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars)
}

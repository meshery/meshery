import environment from "../../../lib/relayEnvironment";
import { fetchQuery, graphql } from "relay-runtime";

function meshScanQuery(variables) {
  let filter_variables = {
    "filter": {
      "type": variables.type,
    }
  }
  let query = graphql`
      query meshScanQuery(
        $filter : ControlPlaneFilter
        ) {
        getControlPlanes (
          filter: $filter
        ) {
          name
          version
          members {
            component
            status
            namespace
          }
        }
      }
    `;
  const data = fetchQuery(environment, query, filter_variables);
  return data;
}

export default meshScanQuery;

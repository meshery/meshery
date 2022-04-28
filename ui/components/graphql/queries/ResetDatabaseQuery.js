import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";


export default function resetDatabase(variables) {
  const vars = { selector : variables.selector }

  const query = graphql`
        query ResetDatabaseQuery($selector: ReSyncActions!) {
           resetStatus: resyncCluster(selector: $selector) 
        }
    `;

  return fetchQuery(environment, query, vars)

}

import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchPerformanceResults(variables) {
  const vars = { selector : variables.selector,
    profileID : variables.profileID }

  const query = graphql`
        query PerformanceResultQuery($selector: PageFilter!, $profileID: String!) {
                fetchResults(selector: $selector, profileID: $profileID) {
                    page
                    page_size
                    total_count
                    results {
                        meshery_id
                        name
                        mesh
                        performance_profile
                        test_id
                        server_metrics
                        test_start_time
                        created_at
                        user_id
                        updated_at
                        runner_results
                    }
                }
            }
        `;

  return fetchQuery(environment, query, vars)
}

import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchAllResults(variables) {
  const environment = createRelayEnvironment({});
  const vars = { selector: variables.selector };

  const FetchAllResultsQueryNode = graphql`
    query FetchAllResultsQuery($selector: PageFilter!) {
      fetchAllResults(selector: $selector) {
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

  return fetchQuery(environment, FetchAllResultsQueryNode, vars);
}

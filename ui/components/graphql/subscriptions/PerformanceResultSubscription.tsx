import { graphql } from 'react-relay';
import { createSubscription } from '../../../lib/subscriptionHelper';

const performanceResultSubscription = graphql`
  subscription PerformanceResultSubscription($selector: PageFilter!, $profileID: String!) {
    subscribePerfResults(selector: $selector, profileID: $profileID) {
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

export default function subscribePerformanceResults(dataCB, variables) {
  return createSubscription({
    subscription: performanceResultSubscription,
    variables,
    onNext: dataCB,
    subscriptionName: 'PerformanceResult',
  });
}

import { graphql } from 'react-relay';
import { createSubscription } from '../../../lib/subscriptionHelper';

const performanceProfilesSubscription = graphql`
  subscription PerformanceProfilesSubscription($selector: PageFilter!) {
    subscribePerfProfiles(selector: $selector) {
      page
      page_size
      total_count
      profiles {
        concurrent_request
        created_at
        duration
        endpoints
        id
        last_run
        load_generators
        name
        qps
        total_results
        updated_at
        user_id
        request_body
        request_cookies
        request_headers
        content_type
        service_mesh
        metadata
      }
    }
  }
`;

export default function subscribePerformanceProfiles(dataCB, variables) {
  return createSubscription({
    subscription: performanceProfilesSubscription,
    variables,
    onNext: dataCB,
    subscriptionName: 'PerformanceProfiles',
  });
}

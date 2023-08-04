import { graphql, requestSubscription } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

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
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription : performanceProfilesSubscription,
    variables : variables,
    onNext : dataCB,
    onError : (error) => console.log(`requestSubscription error:`, error),
  });
}
import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

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
export default function subscribePerformanceProfiles(dataCB, variables) {
  return requestSubscription(environment, {
    subscription : performanceResultSubscription,
    variables : variables,
    onNext : dataCB,
    onError : (error) => console.log(`requestSubscription error:`, error),
  });
}
import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchPerformanceProfiles(variables) {
  const vars = { selector : variables.selector }

  const query = graphql`
        query PerformanceProfilesQuery($selector: PageFilter!) {
          getPerformanceProfiles(selector: $selector) {
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
            }
          }
        }
  `;

  return fetchQuery(environment, query, vars)
}

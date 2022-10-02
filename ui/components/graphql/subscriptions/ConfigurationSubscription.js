import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";


const configurationSubscription = graphql`
  subscription ConfigurationSubscription($applicationSelector: PageFilter!, $patternSelector: PageFilter!, $filterSelector: PageFilter!) {
    configuration: subscribeConfiguration(applicationSelector: $applicationSelector, patternSelector: $patternSelector, filterSelector: $filterSelector) {
      applications {
        page
        page_size
        total_count
        applications {
          id
          name
          application_file
          type {
            String
            Valid
          }
          user_id
          created_at
          updated_at
        }
      }
      patterns {
        page
        page_size
        total_count
        patterns {
          id
          name
          user_id
          pattern_file
          visibility
          catalog_data
          canSupport
          errmsg
          created_at
          updated_at
        }
      }
      filters {
        page
        page_size
        total_count
        filters {
          id
          name
          filter_file
          visibility
          catalog_data
          user_id
          created_at
          updated_at
        }
      }
    }
  }
`;

export default function ConfigurationSubscription(onNext, variables) {
  return requestSubscription(environment, {
    subscription : configurationSubscription,
    variables : variables,
    onNext : onNext,
    onError : error => console.log("ERROR OCCURED IN CONFIGURATION SUBCRIPTION", error)
  });
}

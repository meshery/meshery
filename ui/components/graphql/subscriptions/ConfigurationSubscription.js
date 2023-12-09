import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

const configurationSubscription = graphql`
  subscription ConfigurationSubscription(
    $patternSelector: PageFilter!
    $filterSelector: PageFilter!
  ) {
    configuration: subscribeConfiguration(
      patternSelector: $patternSelector
      filterSelector: $filterSelector
    ) {
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
          type {
            String
            Valid
          }
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
          filter_resource
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
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription: configurationSubscription,
    variables: variables,
    onNext: onNext,
    onError: (error) => console.log('ERROR OCCURED IN CONFIGURATION SUBCRIPTION', error),
  });
}

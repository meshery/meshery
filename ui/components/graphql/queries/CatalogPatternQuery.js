import { fetchQuery, graphql } from 'relay-runtime';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchCatalogPattern(variables) {
  const environment = createRelayEnvironment({});

  const CatalogPatternQueryNode = graphql`
    query CatalogPatternQuery($selector: CatalogSelector!) {
      catalogPatterns: fetchPatternCatalogContent(selector: $selector) {
        id
        name
        user_id
        pattern_file
        visibility
        catalog_data
        created_at
        updated_at
      }
    }
  `;

  return fetchQuery(environment, CatalogPatternQueryNode, variables);
}

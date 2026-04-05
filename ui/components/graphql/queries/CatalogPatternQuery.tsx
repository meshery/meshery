import { fetchQuery, graphql } from 'relay-runtime';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

interface CatalogSelector {
  [key: string]: unknown;
}

interface CatalogPatternVariables {
  selector: CatalogSelector;
}

export default function fetchCatalogPattern(variables: CatalogPatternVariables) {
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

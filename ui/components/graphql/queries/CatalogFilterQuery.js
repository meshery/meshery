import { fetchQuery, graphql } from 'relay-runtime';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchCatalogFilter(variables) {
  const environment = createRelayEnvironment({});

  const CatalogFilterQueryNode = graphql`
    query CatalogFilterQuery($selector: CatalogSelector!) {
      catalogFilters: fetchFilterCatalogContent(selector: $selector) {
        id
        name
        user_id
        filter_file
        filter_resource
        visibility
        catalog_data
        created_at
        updated_at
      }
    }
  `;

  return fetchQuery(environment, CatalogFilterQueryNode, variables);
}

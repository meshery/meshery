import { fetchQuery, graphql } from "relay-runtime";
import environment from "../../../lib/relayEnvironment";

export default function fetchCatalogFilter(variables) {
  const query = graphql`
      query CatalogFilterQuery($selector: CatalogSelector!) {
        catalogFilters: fetchFilterCatalogContent(selector: $selector) {
            id
            name
            user_id
            filter_file
            visibility
            catalog_data
            created_at
            updated_at
        }
      }
    `;

  return fetchQuery(environment, query, variables)
}
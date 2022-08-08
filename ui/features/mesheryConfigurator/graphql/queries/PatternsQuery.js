import environment from "@/app/relayEnvironment";
import { graphql, fetchQuery } from "react-relay";


export default function fetchPatternsQuery(variables) {
    const vars = { selector : variables.selector }

    const query = graphql`
     query PatternsQuery($selector: PageFilter!) {
      getPatterns(selector: $selector) {
        page
        page_size
        patterns{
            canSupport
            created_at
            errmsg
            id
            location{
                branch
                host
                path
                type
            }
            name
            pattern_file
            updated_at
            user_id
        }
        total_count
      } 
    }
    `;

    return fetchQuery(environment, query, vars)
}
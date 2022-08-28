import environment from "@/app/relayEnvironment";
import { graphql, fetchQuery } from "react-relay";


export default function fetchApplicationsQuery(variables) {

    const vars = { selector : variables.selector }

    const query = graphql`
    query ApplicationsQuery($selector: PageFilter!) {
        fetchApplications(selector: $selector) {
         page
         page_size
         total_count
         applications{
            id
            name
            user_id
            application_file
            location{
                branch
                host
                path
                type
            }
            created_at
            updated_at
         }
        }
    }
    `;

    return fetchQuery(environment, query, vars)
}
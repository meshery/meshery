import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function NatsStatusQuery() {

  const query = graphql`
    query NatsStatusQuery {
      controller: getNatsStatus{
            name
            version
            status
        }
    }
  `;

  return fetchQuery(environment, query);
}

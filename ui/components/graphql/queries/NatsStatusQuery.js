import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function NatsStatusQuery(vars) {

  const query = graphql`
    query NatsStatusQuery($k8scontextID: String!) {
      controller: getNatsStatus(k8scontextID: $k8scontextID){
            name
            version
            status
        }
    }
  `;

  return fetchQuery(environment, query, vars);
}

import { graphql, fetchQuery } from "react-relay";
import environment from "../environment";

export function fetchAvailableAddons(variables) {
  const vars = {
    selector: variables.serviceMesh
  }

  const query = graphql`
        query getAvailableAddonsQuery{
            getAvailableAddons($selector: )
        }
    `;

  return fetchQuery(environment, query, vars)
}

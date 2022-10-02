import { graphql, requestSubscription } from "relay-runtime";
import environment from "../../../lib/relayEnvironment";

const k8sContextSubscription = graphql`
subscription K8sContextSubscription($selector: PageFilter!) {
    k8sContext: subscribeK8sContext(selector: $selector) {
      total_count
      contexts {
        id 
        name 
        auth
        cluster
        server
        owner
        created_by
        meshery_instance_id
        kubernetes_server_id
        updated_at
        created_at 
      }  
    }
}
`;

export default function subscribeK8sContext(dataCB, variables) {
  return requestSubscription(environment, {
    subscription : k8sContextSubscription,
    variables : variables,
    onNext : dataCB,
    onError : (error) => console.log("K8sContextSubscription: An error occured:", error)
  })
}
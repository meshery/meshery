import { graphql } from 'relay-runtime';
import { createSubscription } from '../../../lib/subscriptionHelper';

const k8sContextSubscription = graphql`
  subscription K8sContextSubscription($selector: PageFilter!) {
    k8sContext: subscribeK8sContext(selector: $selector) {
      total_count
      contexts {
        id
        name
        server
        owner
        created_by
        meshery_instance_id
        kubernetes_server_id
        deployment_type
        updated_at
        created_at
        version
        connection_id
      }
    }
  }
`;

export default function subscribeK8sContext(dataCB, variables) {
  return createSubscription({
    subscription: k8sContextSubscription,
    variables,
    onNext: dataCB,
    subscriptionName: 'K8sContext',
  });
}

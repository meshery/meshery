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
        createdBy
        mesheryInstanceId
        kubernetesServerId
        deploymentType
        updatedAt
        createdAt
        version
        connectionId
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

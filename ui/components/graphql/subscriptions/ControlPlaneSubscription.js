import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const controlPlaneSubscription = graphql`
  subscription ControlPlaneSubscription($filter: ControlPlaneFilter) {
    controlPlanesState: listenToControlPlaneState(filter: $filter) {
      name
      members {
        name
        version
        component
        namespace
        data_planes {
          name
          image
          status {
            ready
            started
            state
            # state {
            #   waiting {
            #     reason
            #     message
            #   }
            #   running {
            #     startedAt
            #   }
            #   terminated {
            #     reason
            #     message
            #     startedAt
            #     finishedAt
            #     containerID
            #   }
            # }
          }
          ports {
            name
            containerPort
            protocol
          }
          resources {
            limits {
              cpu
              memory
            }
            requests {
              cpu
              memory
            }
          }
        }
      }
    }
  }
`;

export default function subscribeControlPlaneEvents(dataCB, variables) {
  requestSubscription(environment, {
    subscription: controlPlaneSubscription,
    variables: {
      filter: variables
    },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}

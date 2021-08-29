import { dataPlaneSubscription } from "./DataPlanesSubscription";
import { controlPlaneSubscription } from "./ControlPlaneSubscription";
import { requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function subscribeServiceMeshEvents(dataCB, variables) {
  requestSubscription(environment, {
    subscription : dataPlaneSubscription,
    variables : {
      filter : variables,
    },
    onNext : (controlPlanesRes) => {
      dataCB(controlPlanesRes, null)
      requestSubscription(environment, {
        subscription : controlPlaneSubscription,
        variables : {
          filter : variables,
        },
        onNext : (dataPlanesRes) => dataCB(controlPlanesRes, dataPlanesRes),
        onError : (error) => console.log(`An error occured:`, error),
      });
    },
    onError : (error) => console.log(`An error occured:`, error),
  });

}

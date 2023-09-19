import { dataPlaneSubscription } from './DataPlanesSubscription';
import { controlPlaneSubscription } from './ControlPlaneSubscription';
import { requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function subscribeServiceMeshEvents(dataCB, variables) {
  const environment = createRelayEnvironment({});
  let controlRes = null;

  requestSubscription(environment, {
    subscription: controlPlaneSubscription,
    variables: {
      filter: variables,
    },
    onNext: (res) => {
      // decrease state change in dashboard

      // dataCB(res, null)
      controlRes = res;
    },
    onError: (error) => console.log(`An error occured:`, error),
  });

  requestSubscription(environment, {
    subscription: dataPlaneSubscription,
    variables: {
      filter: variables,
    },
    onNext: (dataRes) => {
      if (controlRes) dataCB(controlRes, dataRes);
    },
    onError: (error) => {
      if (controlRes) dataCB(controlRes, null);
      console.log(`An error occured:`, error);
    },
  });
}

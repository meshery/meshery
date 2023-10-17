/* eslint-disable no-unused-vars */
import MesheryOperatorIcon from '../icons/MesheryOperatorIcon.js';
import fetchMesheryOperatorStatus from '../../graphql/queries/OperatorStatusQuery';
import subscribeOperatorStatusEvents from '../../graphql/subscriptions/OperatorStatusSubscription';
import ServiceCard from '../ServiceCard';
import { CircularProgress, Grid } from '@material-ui/core';
import MesheryOperatorDataPanel from '../DataPanels/MesheryOperator';
import { useEffect, useState } from 'react';
import {
  getOperatorStatusFromQueryResult,
  isMesheryOperatorConnected,
} from '../helpers/mesheryOperator.js';

// Connection Wizard
// TODO: bind to contextID prop, leaving due to no use in current UI
const MesheryOperatorScreen = ({ setStepStatus }) => {
  const [operatorInformation, setOperatorInformation] = useState({
    operatorInstalled: false,
    NATSInstalled: false,
    meshSyncInstalled: false,
    operatorVersion: 'N/A',
    meshSyncVersion: 'N/A',
    NATSVersion: 'N/A',
  });
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const mesheryOperatorinfo = {
    name: 'Meshery Operator',
    logoComponent: MesheryOperatorIcon,
    configComp: <div />,
    operatorInformation,
  };

  const showDataPanel = () => isMesheryOperatorConnected(operatorInformation);

  useEffect(() => {
    setStepStatus((prev) => ({ ...prev, operator: isConnected }));
  }, [isConnected]);

  useEffect(() => {
    subscribeOperatorStatusEvents(setOperatorState);
    setIsLoading(true);
    fetchMesheryOperatorStatus().subscribe({
      next: (res) => {
        setIsLoading(false);
        setOperatorState(res);
      },
      error: (err) => setIsLoading(false),
    });
  }, []);

  useEffect(() => {
    setIsConnected(isMesheryOperatorConnected(operatorInformation));
  }, [operatorInformation]);

  const setOperatorState = (res) => {
    const [isReachable, operatorInformation] = getOperatorStatusFromQueryResult(res);
    setOperatorInformation(operatorInformation);
  };

  return (
    <Grid item xs={12} container justify="center" alignItems="flex-start">
      <Grid
        item
        container
        justify="center"
        alignItems="flex-start"
        lg={6}
        sm={12}
        md={12}
        style={{ paddingLeft: '1rem' }}
      >
        <ServiceCard
          serviceInfo={mesheryOperatorinfo}
          isConnected={isConnected}
          style={{ paddingRight: '1rem' }}
        />
      </Grid>
      <Grid item lg={6} sm={12} md={12} container justify="center">
        {isLoading ? (
          <CircularProgress />
        ) : (
          showDataPanel() && <MesheryOperatorDataPanel operatorInformation={operatorInformation} />
        )}
      </Grid>
    </Grid>
  );
};

export default MesheryOperatorScreen;

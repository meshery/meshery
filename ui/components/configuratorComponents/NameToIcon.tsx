import {
  AddCircle,
  BuildRounded,
  DirectionsCar,
  Explore,
  FileCopy,
  Filter,
  ListAlt,
  Lock,
  SimCard,
  SupervisedUserCircle,
} from '@mui/icons-material';
import TouchApp from '@mui/icons-material';
import React from 'react';

/**
 * nameToIcon returns icons for the rjsf form titles
 *
 * @param {string} name Tooltip name
 * @param {string} color default value is primary
 * @returns CustomIconButton
 */
export default function NameToIcon({ name, color = '#607D8B', ...other }) {
  const CustomIcon = ({ Icon }) => <Icon style={{ color }} {...other} />;

  switch (name) {
    //core
    case 'Application':
      return <CustomIcon Icon={TouchApp} />;
    case 'Kubernetes Service':
      return <CustomIcon Icon={BuildRounded} />;
    // Istio
    case 'AuthorizationPolicy':
      return <CustomIcon Icon={Lock} />;
    case 'DestinationRule':
      return <CustomIcon Icon={Explore} />;
    case 'EnvoyFilter':
      return <CustomIcon Icon={Filter} />;
    case 'Gateway':
      return <CustomIcon Icon={ListAlt} />;
    case 'PeerAuthentication':
      return <CustomIcon Icon={FileCopy} />;
    case 'Sidecar':
      return <CustomIcon Icon={DirectionsCar} />;
    case 'VirtualService':
      return <CustomIcon Icon={SupervisedUserCircle} />;
    case 'WorkloadEntry':
      return <CustomIcon Icon={SimCard} />;
    // default
    default:
      return <CustomIcon Icon={AddCircle} />;
  }
}

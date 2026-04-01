import {
  AddCircle,
  BuildRounded,
  DirectionsCar,
  Filter,
  SimCard,
  SupervisedUserCircle,
  TouchApp,
} from '@mui/icons-material';
import ExploreIcon from '@mui/icons-material/Explore';
import FileCopyIcon from '@mui/icons-material/FileCopy';
import ListAltIcon from '@mui/icons-material/ListAlt';
import LockIcon from '@mui/icons-material/Lock';
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
      return <CustomIcon Icon={LockIcon} />;
    case 'DestinationRule':
      return <CustomIcon Icon={ExploreIcon} />;
    case 'EnvoyFilter':
      return <CustomIcon Icon={Filter} />;
    case 'Gateway':
      return <CustomIcon Icon={ListAltIcon} />;
    case 'PeerAuthentication':
      return <CustomIcon Icon={FileCopyIcon} />;
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

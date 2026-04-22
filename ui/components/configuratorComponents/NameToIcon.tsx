import {
  AddCircleIcon,
  BuildRoundedIcon,
  DirectionsCarIcon,
  ExploreIcon,
  FileCopyIcon,
  FilterIcon,
  ListAltIcon,
  LockIcon,
  SimCardIcon,
  SupervisedUserCircleIcon,
  TouchAppIcon,
} from '@sistent/sistent';
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
      return <CustomIcon Icon={TouchAppIcon} />;
    case 'Kubernetes Service':
      return <CustomIcon Icon={BuildRoundedIcon} />;
    // Istio
    case 'AuthorizationPolicy':
      return <CustomIcon Icon={LockIcon} />;
    case 'DestinationRule':
      return <CustomIcon Icon={ExploreIcon} />;
    case 'EnvoyFilter':
      return <CustomIcon Icon={FilterIcon} />;
    case 'Gateway':
      return <CustomIcon Icon={ListAltIcon} />;
    case 'PeerAuthentication':
      return <CustomIcon Icon={FileCopyIcon} />;
    case 'Sidecar':
      return <CustomIcon Icon={DirectionsCarIcon} />;
    case 'VirtualService':
      return <CustomIcon Icon={SupervisedUserCircleIcon} />;
    case 'WorkloadEntry':
      return <CustomIcon Icon={SimCardIcon} />;
    // default
    default:
      return <CustomIcon Icon={AddCircleIcon} />;
  }
}

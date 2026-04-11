import React from 'react';
import SearchIcon from '../../assets/icons/search';
import AddIcon from '../../assets/icons/AddIcon';
import ConfigurationIcon from '../../assets/icons/ConfigurationIcon';
import { DeploymentSelectorIcon } from '../../assets/icons/DeploymentSelectorIcon';
import FilterIcon from '../../assets/icons/filter';
import ComponentIcon from '../../assets/icons/Component';
import OrgIcon from '../../assets/icons/OrgIcon';
import TipsIcon from '../../assets/icons/Tipsicon';
import CopyIcon from '../../assets/icons/CopyIcon';
import DocumentIcon from '../../assets/icons/DocumentIcon';
import CredentialIcon from '../../assets/icons/CredentialIcon';

/**
 * nameToIcon returns icons for the rjsf form titles
 *
 * @param {string} name Tooltip name
 * @param {string} color default value is primary
 * @returns CustomIconButton
 */
export default function NameToIcon({ name, ...other }) {
  const CustomIcon = ({ Icon }) => <Icon {...other} />;

  switch (name) {
    //core
    case 'Application':
      return <CustomIcon Icon={TipsIcon} />;
    case 'Kubernetes Service':
      return <CustomIcon Icon={ConfigurationIcon} />;
    // Istio
    case 'AuthorizationPolicy':
      return <CustomIcon Icon={CredentialIcon} />;
    case 'DestinationRule':
      return <CustomIcon Icon={SearchIcon} />;
    case 'EnvoyFilter':
      return <CustomIcon Icon={FilterIcon} />;
    case 'Gateway':
      return <CustomIcon Icon={DocumentIcon} />;
    case 'PeerAuthentication':
      return <CustomIcon Icon={CopyIcon} />;
    case 'Sidecar':
      return <CustomIcon Icon={DeploymentSelectorIcon} />;
    case 'VirtualService':
      return <CustomIcon Icon={OrgIcon} />;
    case 'WorkloadEntry':
      return <CustomIcon Icon={ComponentIcon} />;
    // default
    default:
      return <CustomIcon Icon={AddIcon} />;
  }
}

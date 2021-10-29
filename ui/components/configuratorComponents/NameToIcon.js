import {
  IconButton, Tooltip
} from "@material-ui/core";
import { AddCircle, BuildRounded, DirectionsCar, Filter, SimCard, SupervisedUserCircle, TouchApp } from "@material-ui/icons";
import ExploreIcon from '@material-ui/icons/Explore';
import FileCopyIcon from '@material-ui/icons/FileCopy';
import ListAltIcon from '@material-ui/icons/ListAlt';
import LockIcon from '@material-ui/icons/Lock';
import React from "react";

/**
 * nameToIcon returns icons for the rjsf form titles
 *
 * @param {string} name Tooltip name
 * @param {Function} action onCLick handler for IconButton
 * @param {string} color default value is primary
 * @returns CustomIconButton
 */
export default function NameToIcon(props) {
  const { name, action, color = "#607D8B", key } = props

  function CustomIcon({ Icon }) {
    return (
      <Tooltip
        key={key}
        title={name}>
        <IconButton onClick={action}>
          <Icon style={{ color }} />
        </IconButton>
      </Tooltip>
    )
  }

  switch (name) {
    //core
    case "Application": return <CustomIcon Icon={TouchApp} />
    case "Kubernetes Service": return <CustomIcon Icon={BuildRounded} />
    // Istio
    case "AuthorizationPolicy": return <CustomIcon Icon={LockIcon} />
    case "DestinationRule": return <CustomIcon Icon={ExploreIcon} />
    case "EnvoyFilter": return <CustomIcon Icon={Filter} />
    case "Gateway": return <CustomIcon Icon={ListAltIcon} />
    case "PeerAuthentication": return <CustomIcon Icon={FileCopyIcon} />
    case "Sidecar": return <CustomIcon Icon={DirectionsCar} />
    case "VirtualService": return <CustomIcon Icon={SupervisedUserCircle} />
    case "WorkloadEntry": return <CustomIcon Icon={SimCard} />
    // default
    default: return <CustomIcon Icon={AddCircle} />
  }
}
import React, { useEffect, useState } from 'react';
import { Button, Grid } from '@material-ui/core';
import { ArrowForward, Delete, Edit } from '@material-ui/icons';

import OrgIcon from '../../assets/icons/OrgIcon';
import {
  AllocationWorkspace,
  BulkSelectCheckbox,
  CardTitle,
  CardWrapper,
  DateLabel,
  DescriptionLabel,
  EmptyDescription,
  OrganizationName,
  StyledIconButton,
  TabCount,
  TabTitle,
} from './styles';
import theme from '../../themes/app';
import FlipCard from '../Environments/flip-card';
import SyncAltIcon from '@mui/icons-material/SyncAlt';
import { useGetDesignsOfWorkspaceQuery, useGetEnvironmentsOfWorkspaceQuery } from '../../rtk-query/workspace';

export const formattoLongDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

export const TransferButton = ({ title, count, onAssign, disabled, classes }) => {
  return (
    <Button onClick={onAssign} disabled={disabled} color="primary" className={classes.popupButton}>
      <Grid>
        <TabCount>{count}</TabCount>
        <TabTitle>{title}</TabTitle>
        <SyncAltIcon />
      </Grid>
    </Button>
  );
};

export const RedirectButton = ({ title, count, disabled=true, classes }) => {
  return (
    <Button disabled={disabled} color="primary" className={classes.popupButton}>
      <Grid>
        <TabCount>{count}</TabCount>
        <TabTitle>{title}</TabTitle>
        <ArrowForward />
      </Grid>
    </Button>
  );
};

/**
 * Renders a Workspace card component.
 *
 * @param {Object} props - The component props.
 * @param {Object} props.workspaceDetails - The details of the workspace.
 * @param {string} props.workspaceDetails.name - The name of the workspace.
 * @param {string} props.workspaceDetails.description - The description of the workspace.
 * @param {Function} props.onDelete - Function to delete the workspace.
 * @param {Function} props.onEdit - Function to edit the workspace.
 * @param {Function} props.onSelect - Function to select workspace for bulk actions.
 * @param {Array} props.selectedWorkspaces - Selected workspace list for delete.
 * @param {Function} props.onAssignTeam - Function to open team assignment modal open.
 * @param {Function} props.onAssignDesign - Function to open design assignment modal open.
 *
 */

const WorkspaceCard = ({
  workspaceDetails,
  // onDelete,
  onEdit,
  onSelect,
  selectedWorkspaces,
  onAssignTeam,
  onAssignEnvironment,
  onAssignDesign,
  classes,
}) => {
  const [skip, setSkip] = useState(true);

  const deleted = workspaceDetails.deleted_at.Valid;

  const { data: environmentsOfWorkspace } = useGetEnvironmentsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
    },
    {
      skip,
    },
  );

  const { data: designsOfWorkspace } = useGetDesignsOfWorkspaceQuery({
      workspaceId: workspaceDetails.id
    },
    {
      skip
    }
  );

  useEffect(() => {
    if (!deleted) {
      setSkip(false);
    } else {
      setSkip(true);
    }
  }, [workspaceDetails, deleted]);

  const teamsOfWorkspaceCount = 0;

  const environmentsOfWorkspaceCount = environmentsOfWorkspace?.total_count
    ? environmentsOfWorkspace.total_count
    : 0;

  const designsOfWorkspaceCount = designsOfWorkspace?.total_count
    ? designsOfWorkspace.total_count
    : 0;

  return (
    <>
      <FlipCard
        disableFlip={
          selectedWorkspaces?.filter((id) => id == workspaceDetails.id).length === 1 ? true : false
        }
        frontComponents={
          <CardWrapper
            elevation={2}
            // minHeight={{ xs: "520px", sm: "390px" }}
            style={{
              minHeight: '390px',
            }}
          >
            <Grid style={{ display: 'flex', flexDirection: 'row', paddingBottom: '10px' }}>
              <CardTitle variant="body2" onClick={(e) => e.stopPropagation()}>
                {workspaceDetails?.name}
              </CardTitle>
            </Grid>
            <Grid style={{ display: 'flex', alignItems: 'center', marginTop: '10px', marginBottom: "10px" }}>
              <StyledIconButton onClick={(e) => e.stopPropagation()}>
                <OrgIcon width="24" height="24" />
              </StyledIconButton>
              <OrganizationName variant="span" onClick={(e) => e.stopPropagation()}>
                {workspaceDetails?.owner}
              </OrganizationName>
            </Grid>
            <Grid
              pt={{ xs: 1.5, md: 3 }}
              style={{
                display: 'flex',
                gap: '10px',
                flexDirection: { xs: 'column', sm: 'row' },
              }}
            >
              <Grid xs={12} sm={4}>
                <AllocationWorkspace onClick={(e) => e.stopPropagation()}>
                  <RedirectButton title="Connections" count={0} classes={classes} />
                  <TransferButton
                    title="Environments"
                    count={environmentsOfWorkspaceCount}
                    onAssign={onAssignEnvironment}
                    classes={classes}
                  />
                </AllocationWorkspace>
              </Grid>
              <Grid xs={12} sm={4} style={{ display: 'flex', justifyContent: 'center' }}>
                <AllocationWorkspace onClick={(e) => e.stopPropagation()}>
                  <RedirectButton disabled={true} title="Users" count={0} link="/identity/users" classes={classes} />
                  <TransferButton
                    title="Teams"
                    count={teamsOfWorkspaceCount}
                    onAssign={onAssignTeam}
                    disabled={true}
                    classes={classes}
                  />
                </AllocationWorkspace>
              </Grid>
              <Grid xs={12} sm={4} style={{ display: 'flex', justifyContent: 'flex-end' }}>
                <AllocationWorkspace onClick={(e) => e.stopPropagation()}>
                  <TransferButton
                    title="Designs"
                    count={designsOfWorkspaceCount}
                    onAssign={onAssignDesign}
                    disabled={true}
                    classes={classes}
                  />
                  <RedirectButton title="Deployments" count={0} classes={classes} />
                </AllocationWorkspace>
              </Grid>
            </Grid>
          </CardWrapper>
        }
        backComponents={
          <CardWrapper
            elevation={2}
            // minHeight={{ xs: "520px", sm: "390px" }}
            style={{
              background: 'linear-gradient(180deg, #007366 0%, #000 100%)',
              minHeight: '390px',
            }}
          >
            <Grid xs={12}>
              <Grid xs={12} style={{ display: 'flex', flexDirection: 'row' }}>
                <Grid xs={6} style={{ display: 'flex', alignItems: 'flex-start' }}>
                  <BulkSelectCheckbox onClick={(e) => e.stopPropagation()} onChange={onSelect} />
                  <CardTitle
                    style={{ color: theme.palette.secondary.white }}
                    variant="body2"
                    onClick={(e) => e.stopPropagation()}
                  >
                    {workspaceDetails?.name}
                  </CardTitle>
                </Grid>
                <Grid
                  xs={6}
                  style={{
                    display: 'flex',
                    alignItems: 'flex-start',
                    justifyContent: 'flex-end',
                  }}
                >
                  <Button
                    style={{
                      minWidth: 'fit-content',
                      '&.MuiButtonBase-root:hover': {
                        background: 'transparent',
                      },
                      padding: 0,
                    }}
                    onClick={onEdit}
                    disabled={
                      selectedWorkspaces?.filter((id) => id == workspaceDetails.id).length === 1
                        ? true
                        : false
                    }
                  >
                    <Edit style={{ color: 'white', margin: '0 2px' }} />
                  </Button>
                  <Button
                    style={{
                      minWidth: 'fit-content',
                      '&.MuiButtonBase-root:hover': {
                        background: 'transparent',
                      },
                      padding: 0,
                    }}
                    onClick={onEdit}
                    disabled={
                      selectedWorkspaces?.filter((id) => id == workspaceDetails.id).length === 1
                        ? true
                        : false
                    }
                  >
                    <Delete style={{ color: 'white', margin: '0 2px' }} />
                  </Button>
                </Grid>
              </Grid>
              <Grid>
                {workspaceDetails.description ? (
                  <DescriptionLabel
                    onClick={(e) => e.stopPropagation()}
                    style={{ padding: '10px 0', color: theme.palette.secondary.white, maxHeight: '105px' }}
                  >
                    {workspaceDetails.description}
                  </DescriptionLabel>
                ) : (
                  <EmptyDescription
                    onClick={(e) => e.stopPropagation()}
                    style={{ padding: '10px 0', color: `${theme.palette.secondary.white}90` }}
                  >
                    No description
                  </EmptyDescription>
                )}
              </Grid>
            </Grid>
            <Grid
              style={{
                display: 'flex',
                flexDirection: 'row',
                position: 'absolute',
                bottom: '20px',
                width: '100%',
                color: `${theme.palette.secondary.white}99`,
              }}
            >
              <Grid xs={6} style={{ textAlign: 'left' }}>
                <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
                  Updated At: {formattoLongDate(workspaceDetails?.updated_at)}
                </DateLabel>
              </Grid>
              <Grid xs={6} style={{ textAlign: 'left' }}>
                <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
                  Created At: {formattoLongDate(workspaceDetails?.created_at)}
                </DateLabel>
              </Grid>
            </Grid>
          </CardWrapper>
        }
      />
    </>
  );
};

export default WorkspaceCard;

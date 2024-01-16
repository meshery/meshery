import React, { useEffect, useState } from 'react';
import { Box, Button, Grid } from '@material-ui/core';
import { ArrowForward, Edit } from '@material-ui/icons';
import { DeleteIcon } from '@layer5/sistent-svg';
import {
  BulkSelectCheckbox,
  CardTitle,
  CardWrapper,
  DateLabel,
  DescriptionLabel,
  EmptyDescription,
  TabCount,
  TabTitle,
} from './styles';
import theme from '../../../themes/app';
import { FlipCard } from '../General';
import SyncAltIcon from '@mui/icons-material/SyncAlt';
import {
  useGetDesignsOfWorkspaceQuery,
  useGetEnvironmentsOfWorkspaceQuery,
} from '../../../rtk-query/workspace';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';

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
        <SyncAltIcon style={{ position: 'absolute', top: '10px', right: '10px' }} />
      </Grid>
    </Button>
  );
};

export const RedirectButton = ({ title, count, disabled = true, classes }) => {
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
 * @param {Function} props.onAssignEnvironment - Function to open environment assignment modal open.
 * @param {Function} props.onAssignDesign - Function to open design assignment modal open.
 *
 */

const WorkspaceCard = ({
  workspaceDetails,
  onDelete,
  onEdit,
  onSelect,
  selectedWorkspaces,
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

  const { data: designsOfWorkspace } = useGetDesignsOfWorkspaceQuery(
    {
      workspaceId: workspaceDetails.id,
    },
    {
      skip,
    },
  );

  useEffect(() => {
    if (!deleted) {
      setSkip(false);
    } else {
      setSkip(true);
    }
  }, [workspaceDetails, deleted]);

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
          <CardFront
            classes={classes}
            name={workspaceDetails?.name}
            description={workspaceDetails?.description}
            environmentsCount={environmentsOfWorkspaceCount}
            onAssignEnvironment={onAssignEnvironment}
            designsCount={designsOfWorkspaceCount}
            onAssignDesign={onAssignDesign}
          />
        }
        backComponents={
          <CardBack
            classes={classes}
            name={workspaceDetails?.name}
            workspaceId={workspaceDetails?.id}
            onEdit={onEdit}
            onDelete={onDelete}
            selectedWorkspaces={selectedWorkspaces}
            updatedDate={workspaceDetails?.updated_at}
            createdDate={workspaceDetails?.created_at}
            onSelect={onSelect}
          />
        }
      />
    </>
  );
};

export default WorkspaceCard;

const CardFront = ({
  classes,
  name,
  description,
  environmentsCount,
  onAssignEnvironment,
  designsCount,
  onAssignDesign,
}) => {
  return (
    <CardWrapper
      elevation={2}
      style={{
        minHeight: '300px',
      }}
    >
      <Grid style={{ display: 'flex', flexDirection: 'row', paddingBottom: '5px' }}>
        <CardTitle variant="body2" onClick={(e) => e.stopPropagation()}>
          {name}
        </CardTitle>
      </Grid>
      <Grid>
        {description ? (
          <DescriptionLabel
            onClick={(e) => e.stopPropagation()}
            style={{
              padding: '10px 0',
              maxHeight: '105px',
            }}
          >
            {description}
          </DescriptionLabel>
        ) : (
          <EmptyDescription onClick={(e) => e.stopPropagation()} style={{ padding: '10px 0' }}>
            No description
          </EmptyDescription>
        )}
      </Grid>
      <Grid
        pt={{ xs: 1.5, md: 3 }}
        style={{
          display: 'flex',
          gap: '10px',
          flexDirection: { xs: 'column', sm: 'row' },
        }}
      >
        <Grid
          xs={12}
          style={{
            paddingTop: '15px',
            display: 'flex',
            alignItems: 'flex-end',
            justifyContent: 'flex-end',
            gap: '10px',
          }}
        >
          <Box className={classes.allocationButton} onClick={(e) => e.stopPropagation()}>
            <TransferButton
              title="Environments"
              count={environmentsCount}
              onAssign={onAssignEnvironment}
              classes={classes}
            />
          </Box>
          <Box className={classes.allocationButton} onClick={(e) => e.stopPropagation()}>
            <TransferButton
              title="Designs"
              count={designsCount}
              onAssign={onAssignDesign}
              classes={classes}
            />
          </Box>
        </Grid>
      </Grid>
    </CardWrapper>
  );
};

const CardBack = ({
  onSelect,
  name,
  onEdit,
  onDelete,
  selectedWorkspaces,
  workspaceId,
  // description,
  updatedDate,
  createdDate,
}) => {
  return (
    <CardWrapper
      elevation={2}
      style={{
        background: 'linear-gradient(180deg, #007366 0%, #000 100%)',
        minHeight: '300px',
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
              {name}
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
                selectedWorkspaces?.filter((id) => id == workspaceId).length === 1
                  ? true
                  : !CAN(keys.EDIT_WORKSPACE.action, keys.EDIT_WORKSPACE.subject)
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
              onClick={onDelete}
              disabled={
                selectedWorkspaces?.filter((id) => id == workspaceId).length === 1
                  ? true
                  : !CAN(keys.DELETE_WORKSPACE.action, keys.DELETE_WORKSPACE.subject)
              }
            >
              <DeleteIcon fill={theme.palette.secondary.whiteIcon} />
            </Button>
          </Grid>
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
            Updated At: {formattoLongDate(updatedDate)}
          </DateLabel>
        </Grid>
        <Grid xs={6} style={{ textAlign: 'left' }}>
          <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
            Created At: {formattoLongDate(createdDate)}
          </DateLabel>
        </Grid>
      </Grid>
    </CardWrapper>
  );
};

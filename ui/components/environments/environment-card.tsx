import React from 'react';
import { FlipCard } from '../lifecycle/general';
import { useGetEnvironmentConnectionsQuery } from '../../rtk-query/environments';

import { Keys } from '@meshery/schemas/permissions';
import {
  DeleteIcon,
  EditIcon,
  Grid2,
  SyncAltIcon,
  useTheme,
  IconButton,
  CustomTooltip,
  Chip,
} from '@sistent/sistent';
import { iconMedium } from '../../css/icons.styles';

import {
  Name,
  CardWrapper,
  DateLabel,
  DescriptionLabel,
  EmptyDescription,
  TabCount,
  TabTitle,
  PopupButton,
  AllocationButton,
  BulkSelectCheckbox,
  CardTitle,
} from './styles';

export const formattoLongDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

export const TransferButton = ({ title, count, onAssign, permissionKey }) => {
  const theme = useTheme();
  return (
    <PopupButton permissionKey={permissionKey} onClick={onAssign}>
      <Grid2>
        <TabCount>{count}</TabCount>
        <TabTitle>{title}</TabTitle>
        <SyncAltIcon
          style={{
            position: 'absolute',
            top: '10px',
            right: '10px',
          }}
          fill={theme.palette.background?.neutral?.default}
        />
      </Grid2>
    </PopupButton>
  );
};

/**
 * Renders a environment card component.
 *
 * @param {Object} props - The component props.
 * @param {Object} props.environmentDetails - The details of the environment.
 * @param {string} props.environmentDetails.name - The name of the environment.
 * @param {string} props.environmentDetails.description - The description of the environment.
 * @param {string} [props.environmentDetails.purpose] - The purpose of the environment.
 *   When equal to "administrative", the environment is platform-provisioned and shown
 *   with an "Administrative" badge. Edit and delete affordances are disabled for such
 *   environments because actions on them are likely to be refused by the server.
 * @param {Function} props.onDelete - Function to delete the environment.
 * @param {Function} props.onEdit - Function to edit the environment.
 * @param {Function} props.onSelect - Function to select environment for bulk actions.
 * @param {Function} props.onAssignConnection - Function to open connection assignment modal open.
 * @param {Array} props.selectedEnvironments - Selected environments list for delete.
 * @param {String} props.classes - Styles property names for classes.
 *
 */

const EnvironmentCard = ({
  environmentDetails,
  selectedEnvironments,
  onDelete,
  onEdit,
  onSelect,
  onAssignConnection,
}) => {
  const theme = useTheme();
  const { data: environmentConnections } = useGetEnvironmentConnectionsQuery(
    {
      environmentId: environmentDetails.id,
    },
    { skip: !environmentDetails.id },
  );
  const environmentConnectionsCount = environmentConnections?.totalCount || 0;

  // Compare against the literal "administrative".
  // Do NOT test for "not user": the property is optional and absent for every
  // environment that predates it, so a negative test would render those as administrative.
  const isAdministrative = environmentDetails?.purpose === 'administrative';

  // this allows to handle both cases when deleted at is:
  // - timestamp or null
  // - object in format {Time: timestamp, Valid: boolean}
  // --
  // TODO:
  // - switch remote provider to have format of deletedAt as timestamp or null
  // - or update serialisation for deletedAt field of Environment to return object in format {Time: timestamp, Valid: boolean}
  const deleted =
    environmentDetails.deletedAt === null
      ? false
      : typeof environmentDetails.deletedAt === 'object' &&
          environmentDetails.deletedAt !== null &&
          'Valid' in environmentDetails.deletedAt
        ? !!environmentDetails.deletedAt.Valid
        : true;

  return (
    <>
      <FlipCard
        disableFlip={
          selectedEnvironments?.filter((id) => id == environmentDetails.id).length === 1
            ? true
            : false
        }
        frontComponents={
          <CardWrapper
            sx={{
              minHeight: '320px',
              height: '320px',
              borderRadius: 2,
            }}
          >
            <Grid2 sx={{ display: 'flex', flexDirection: 'row', pb: 1, alignItems: 'center' }}>
              <Name variant="body2" onClick={(e) => e.stopPropagation()}>
                {environmentDetails?.name}
              </Name>
              {isAdministrative && (
                <Chip
                  label="Administrative"
                  size="small"
                  sx={{
                    ml: 1,
                    backgroundColor: 'rgba(0, 179, 159, 0.12)',
                    color: '#00B39F',
                    border: '1px solid #00B39F',
                    fontWeight: 600,
                    fontSize: '0.7rem',
                    height: '22px',
                  }}
                />
              )}
            </Grid2>
            <Grid2
              sx={{
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'flex-start',
              }}
            >
              <Grid2
                sx={{ display: 'flex', justifyContent: 'flex-start' }}
                size={{
                  xs: 12,
                  sm: 9,
                  md: 12,
                }}
              >
                {environmentDetails.description ? (
                  <DescriptionLabel
                    onClick={(e) => e.stopPropagation()}
                    sx={{
                      marginBottom: { xs: 2, sm: 0 },
                      paddingRight: { sm: 2, lg: 0 },
                      marginTop: '0px',
                    }}
                  >
                    {environmentDetails.description}
                  </DescriptionLabel>
                ) : (
                  <EmptyDescription
                    onClick={(e) => e.stopPropagation()}
                    sx={{ color: 'rgba(122,132,142,1)' }}
                  >
                    No description
                  </EmptyDescription>
                )}
              </Grid2>
              <Grid2
                size={{
                  xs: 12,
                }}
                sx={{
                  paddingTop: '15px',
                  display: 'flex',
                  alignItems: 'flex-end',
                  justifyContent: 'flex-end',
                  gap: '10px',
                }}
              >
                <AllocationButton onClick={(e) => e.stopPropagation()}>
                  <TransferButton
                    title="Assigned Connections"
                    count={environmentConnectionsCount}
                    onAssign={onAssignConnection}
                    permissionKey={Keys.WorkspaceManagementViewConnections}
                  />
                </AllocationButton>
                {/* temporary disable workspace allocation button  */}
                {/* {false && (
                  <AllocationButton onClick={(e) => e.stopPropagation()}>
                    <TransferButton
                      title="Assigned Workspaces"
                      count={
                        environmentDetails.workspaces ? environmentDetails.workspaces?.length : 0
                      }
                      onAssign={onAssignConnection}
                      disabled={false} // TODO: re-enable with permissionKey={Keys.WorkspaceManagementViewWorkspace}
                    />
                  </AllocationButton>
                )} */}
              </Grid2>
            </Grid2>
          </CardWrapper>
        }
        backComponents={
          <CardWrapper
            elevation={2}
            sx={{
              minHeight: '320px',
              background: 'linear-gradient(180deg, #007366 0%, #000 100%)',
            }}
          >
            <Grid2 sx={{ display: 'flex', flexDirection: 'row' }} size={{ xs: 12 }}>
              <Grid2 sx={{ display: 'flex', alignItems: 'center', gap: 1 }} size={{ xs: 6 }}>
                <BulkSelectCheckbox
                  onClick={(e) => e.stopPropagation()}
                  onChange={onSelect}
                  disabled={deleted ? true : false}
                />
                <CardTitle
                  sx={{ color: 'white' }}
                  variant="body2"
                  onClick={(e) => e.stopPropagation()}
                >
                  {environmentDetails?.name}
                </CardTitle>
                {isAdministrative && (
                  <Chip
                    label="Administrative"
                    size="small"
                    sx={{
                      backgroundColor: 'rgba(255, 255, 255, 0.15)',
                      color: 'white',
                      border: '1px solid rgba(255, 255, 255, 0.5)',
                      fontWeight: 600,
                      fontSize: '0.7rem',
                      height: '22px',
                    }}
                  />
                )}
              </Grid2>
              <Grid2
                size={{ xs: 6 }}
                sx={{
                  display: 'flex',
                  alignItems: 'flex-start',
                  justifyContent: 'flex-end',
                }}
              >
                <CustomTooltip
                  title={isAdministrative ? 'Administrative environments cannot be edited' : 'Edit'}
                >
                  {/* span ensures the tooltip is visible even when the button is disabled */}
                  <span>
                    <IconButton
                      onClick={(ev) => {
                        ev.stopPropagation();
                        onEdit(ev);
                      }}
                      sx={{ color: 'white' }}
                      disabled={
                        isAdministrative ||
                        selectedEnvironments?.filter((id) => id == environmentDetails.id)
                          .length === 1
                      }
                      permissionKey={Keys.WorkspaceManagementEditEnvironment}
                    >
                      <EditIcon
                        style={{ ...iconMedium, margin: '0 2px' }}
                        fill={theme?.palette?.icon?.default}
                      />
                    </IconButton>
                  </span>
                </CustomTooltip>
                <CustomTooltip
                  title={
                    isAdministrative ? 'Administrative environments cannot be deleted' : 'Delete'
                  }
                >
                  {/* span ensures the tooltip is visible even when the button is disabled */}
                  <span>
                    <IconButton
                      onClick={(ev) => {
                        ev.stopPropagation();
                        onDelete(ev);
                      }}
                      sx={{ color: 'white' }}
                      disabled={
                        isAdministrative ||
                        selectedEnvironments?.filter((id) => id == environmentDetails.id)
                          .length === 1
                      }
                      permissionKey={Keys.WorkspaceManagementDeleteEnvironment}
                    >
                      <DeleteIcon
                        style={{ ...iconMedium, margin: '0 2px' }}
                        fill={theme?.palette?.icon?.default}
                      />
                    </IconButton>
                  </span>
                </CustomTooltip>
              </Grid2>
            </Grid2>
            <Grid2 sx={{ display: 'flex', flexDirection: 'row', color: 'white' }}>
              <Grid2 size={{ xs: 6 }} sx={{ textAlign: 'left' }}>
                <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
                  Updated At: {formattoLongDate(environmentDetails?.updatedAt)}
                </DateLabel>
              </Grid2>
              <Grid2 size={{ xs: 6 }} sx={{ textAlign: 'left' }}>
                <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
                  Created At: {formattoLongDate(environmentDetails?.createdAt)}
                </DateLabel>
              </Grid2>
            </Grid2>
          </CardWrapper>
        }
      />
    </>
  );
};

export default EnvironmentCard;

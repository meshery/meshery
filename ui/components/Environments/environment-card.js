import React from 'react';
import { Button, Grid } from '@mui/material';
import {
  AllocationButton,
  BulkSelectCheckbox,
  CardTitle,
  CardWrapper,
  DateLabel,
  DescriptionLabel,
  EmptyDescription,
  PopupButton,
  TabCount,
  TabTitle,
} from './styles';
import SyncAltIcon from '@mui/icons-material/SyncAlt';
import FlipCard from './flip-card';
import { Delete, Edit } from '@material-ui/icons';
import { useGetEnvironmentConnectionsQuery } from '../../rtk-query/environments';

export const formattoLongDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

export const TransferButton = ({ title, count, onAssign }) => {
  return (
    <PopupButton onClick={onAssign}>
      <TabCount>{count}</TabCount>
      <TabTitle>{title}</TabTitle>
      <SyncAltIcon />
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
 * @param {Function} props.onDelete - Function to delete the environment.
 * @param {Function} props.onEdit - Function to edit the environment.
 * @param {Function} props.onSelect - Function to select environment for bulk actions.
 * @param {Function} props.onAssignConnection - Function to open connection assignment modal open.
 * @param {Array} props.selectedEnvironments - Selected environments list for delete.
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
  // const [environmentConnectionsCount, setEenvironmentConnectionsCount] = useState(0);

  const { data: environmentConnections } = useGetEnvironmentConnectionsQuery(
    {
      environmentId: environmentDetails.id,
    },
    { skip: !environmentDetails.id },
  );
  const environmentConnectionsCount = environmentConnections?.connections?.length || 0;

  const deleted = environmentDetails.deleted_at.Valid;

  return (
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
            background: deleted
              ? 'rgba(255, 255, 255, .6)'
              : `linear-gradient(180deg, $#FBFBFB 0%, #F5F5F5 100%)`,
          }}
          elevation={2}
        >
          <Grid sx={{ display: 'flex', flexDirection: 'row', pb: 1 }}>
            <CardTitle variant="body2" onClick={(e) => e.stopPropagation()}>
              {environmentDetails?.name}
            </CardTitle>
          </Grid>
          <Grid
            sx={{
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'flex-start',
            }}
          >
            <Grid xs={12} sm={9} md={12} sx={{ display: 'flex', justifyContent: 'flex-start' }}>
              {environmentDetails.description ? (
                <DescriptionLabel
                  onClick={(e) => e.stopPropagation()}
                  sx={{ marginBottom: { xs: 2, sm: 0 }, paddingRight: { sm: 2, lg: 0 } }}
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
            </Grid>
            <Grid
              xs={12}
              sm={3}
              md={12}
              sx={{
                pt: 2,
                display: 'flex',
                alignItems: 'flex-end',
                justifyContent: 'flex-end',
                gap: 1,
              }}
            >
              <AllocationButton onClick={(e) => e.stopPropagation()}>
                <TransferButton
                  title="Assigned Connections"
                  count={environmentConnectionsCount}
                  onAssign={onAssignConnection}
                />
              </AllocationButton>
              {/* temporary disable workspace allocation button  */}
              {false && (
                <AllocationButton onClick={(e) => e.stopPropagation()}>
                  <TransferButton
                    title="Assigned Workspaces"
                    count={
                      environmentDetails.workspaces ? environmentDetails.workspaces?.length : 0
                    }
                    onAssign={() => {}}
                  />
                </AllocationButton>
              )}
            </Grid>
          </Grid>
        </CardWrapper>
      }
      backComponents={
        <CardWrapper
          elevation={2}
          sx={{
            minHeight: '320px',
            background: deleted
              ? 'rgba(255, 255, 255, .6)'
              : 'linear-gradient(180deg, #007366 0%, #000 100%)',
          }}
        >
          <Grid xs={12} sx={{ display: 'flex', flexDirection: 'row', height: '40px' }}>
            <Grid xs={6} sx={{ display: 'flex', alignItems: 'flex-start' }}>
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
            </Grid>
            <Grid
              xs={6}
              sx={{
                display: 'flex',
                alignItems: 'flex-start',
                justifyContent: 'flex-end',
              }}
            >
              <Button
                sx={{
                  minWidth: 'fit-content',
                  '&.MuiButtonBase-root:hover': {
                    bgcolor: 'transparent',
                  },
                  p: 0,
                }}
                onClick={onEdit}
                disabled={
                  selectedEnvironments?.filter((id) => id == environmentDetails.id).length === 1
                    ? true
                    : false
                }
              >
                <Edit style={{ color: 'white', margin: '0 2px' }} />
              </Button>
              <Button
                sx={{
                  minWidth: 'fit-content',
                  '&.MuiButtonBase-root:hover': {
                    bgcolor: 'transparent',
                  },
                  p: 0,
                }}
                onClick={onDelete}
                disabled={
                  selectedEnvironments?.filter((id) => id == environmentDetails.id).length === 1
                    ? true
                    : false
                }
              >
                <Delete style={{ color: 'white', margin: '0 2px' }} />
              </Button>
            </Grid>
          </Grid>
          <Grid sx={{ display: 'flex', flexDirection: 'row', color: 'white' }}>
            <Grid xs={6} sx={{ textAlign: 'left' }}>
              <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
                Updated At: {formattoLongDate(environmentDetails?.updated_at)}
              </DateLabel>
            </Grid>
            <Grid xs={6} sx={{ textAlign: 'left' }}>
              <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
                Created At: {formattoLongDate(environmentDetails?.created_at)}
              </DateLabel>
            </Grid>
          </Grid>
        </CardWrapper>
      }
    />
  );
};

export default EnvironmentCard;

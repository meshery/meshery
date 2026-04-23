import React from 'react';
import { SyncAlt as SyncAltIcon, Delete, Edit } from '../icons';
import { useGetEnvironmentConnectionsQuery } from '@/rtk-query/environments';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { Grid2, useTheme, Box, IconButton } from '@sistent/sistent';
import ResourceCard from '@/components/shared/Card/ResourceCard';

import {
  DateLabel,
  DescriptionLabel,
  EmptyDescription,
  TabCount,
  TabTitle,
  PopupButton,
  AllocationButton,
  BulkSelectCheckbox,
} from './styles';

export const formattoLongDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
};

export const TransferButton = ({ title, count, onAssign, disabled }) => {
  const theme = useTheme();
  return (
    <PopupButton disabled={disabled} onClick={onAssign}>
      <Grid2>
        <TabCount>{count}</TabCount>
        <TabTitle>{title}</TabTitle>
        <SyncAltIcon
          sx={{
            position: 'absolute',
            top: '10px',
            right: '10px',
            color: theme.palette.background?.neutral?.default,
          }}
        />
      </Grid2>
    </PopupButton>
  );
};

const EnvironmentCard = ({
  environmentDetails,
  selectedEnvironments,
  onDelete,
  onEdit,
  onSelect,
  onAssignConnection,
}) => {
  const { data: environmentConnections } = useGetEnvironmentConnectionsQuery(
    {
      environmentId: environmentDetails.id,
    },
    { skip: !environmentDetails.id },
  );
  const environmentConnectionsCount = environmentConnections?.total_count || 0;

  const deleted =
    environmentDetails.deleted_at === null
      ? false
      : typeof environmentDetails.deleted_at === 'object' &&
          environmentDetails.deleted_at !== null &&
          'Valid' in environmentDetails.deleted_at
        ? !!environmentDetails.deleted_at.Valid
        : true;

  const isFlippingDisabled =
    selectedEnvironments?.filter((id) => id == environmentDetails.id).length === 1;

  return (
    <ResourceCard
      title={environmentDetails?.name}
      isFlippable={!isFlippingDisabled}
      footer={
        <Grid2
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
              disabled={!CAN(keys.VIEW_CONNECTIONS.action, keys.VIEW_CONNECTIONS.subject)}
            />
          </AllocationButton>
        </Grid2>
      }
      backHeaderRight={
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <IconButton
            onClick={onEdit}
            disabled={
              isFlippingDisabled
                ? true
                : !CAN(keys.EDIT_ENVIRONMENT.action, keys.EDIT_ENVIRONMENT.subject)
            }
          >
            <Edit sx={{ color: 'white' }} />
          </IconButton>
          <IconButton
            onClick={onDelete}
            disabled={
              isFlippingDisabled
                ? true
                : !CAN(keys.DELETE_ENVIRONMENT.action, keys.DELETE_ENVIRONMENT.subject)
            }
          >
            <Delete sx={{ color: 'white' }} />
          </IconButton>
        </div>
      }
      backHeaderLeft={
        <BulkSelectCheckbox
          onClick={(e) => e.stopPropagation()}
          onChange={onSelect}
          disabled={deleted}
          sx={{ color: 'white', '&.Mui-checked': { color: 'white' } }}
        />
      }
      backFooter={
        <Grid2 container spacing={1} sx={{ color: 'white', width: '100%', mt: 'auto' }}>
          <Grid2 size={{ xs: 12 }}>
            <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
              Updated At: {formattoLongDate(environmentDetails?.updated_at)}
            </DateLabel>
          </Grid2>
          <Grid2 size={{ xs: 12 }}>
            <DateLabel variant="span" onClick={(e) => e.stopPropagation()}>
              Created At: {formattoLongDate(environmentDetails?.created_at)}
            </DateLabel>
          </Grid2>
        </Grid2>
      }
      backHeaderStyle={{
        background: 'linear-gradient(180deg, #007366 0%, #000 100%)',
        color: 'white',
      }}
      backContentStyle={{
        background: 'linear-gradient(180deg, #007366 0%, #000 100%)',
        color: 'white',
      }}
      backFooterStyle={{
        background: 'linear-gradient(180deg, #007366 0%, #000 100%)',
        color: 'white',
      }}
    >
      <Box sx={{ minHeight: '100px' }}>
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
      </Box>
    </ResourceCard>
  );
};

export default EnvironmentCard;

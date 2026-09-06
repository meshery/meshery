import React from 'react';
import {
  Autocomplete,
  TextField,
  Chip,
  Box,
  useTheme,
  ArrowDropDownIcon,
  alpha,
} from '@sistent/sistent';

export interface EnvironmentOption {
  label: string;
  value: string;
}

export interface WorkspaceEnvironmentSelectionProps {
  workspaceId: string;
  useAssignEnvironmentToWorkspaceMutation: any;
  useGetEnvironmentsOfWorkspaceQuery: any;
  useUnassignEnvironmentFromWorkspaceMutation: any;
  useNotificationHandlers: any;
  isAssignedEnvironmentAllowed?: boolean;
}

const WorkspaceEnvironmentSelection: React.FC<WorkspaceEnvironmentSelectionProps> = ({
  workspaceId,
  useAssignEnvironmentToWorkspaceMutation,
  useGetEnvironmentsOfWorkspaceQuery,
  useUnassignEnvironmentFromWorkspaceMutation,
  useNotificationHandlers,
  isAssignedEnvironmentAllowed = true,
}) => {
  const theme = useTheme();
  const { handleSuccess, handleError } = useNotificationHandlers();

  const iconColor = theme.palette.icon?.secondary || theme.palette.text.secondary;
  const inputBg = theme.palette.background?.card || theme.palette.background.default;
  const borderColor = theme.palette.border?.default || theme.palette.divider;
  const hoverBg = theme.palette.action?.hover || alpha(theme.palette.text.primary, 0.08);

  const { data: availableData, isLoading: isAvailableLoading } = useGetEnvironmentsOfWorkspaceQuery(
    {
      workspaceId,
      page: 0,
      pagesize: 'all',
      filter: '{"assigned":false}',
    },
  );

  const { data: assignedData, isLoading: isAssignedLoading } = useGetEnvironmentsOfWorkspaceQuery({
    workspaceId,
    page: 0,
    pagesize: 'all',
  });

  const availableOptions: EnvironmentOption[] =
    availableData?.environments?.map((env: { name: string; id: string }) => ({
      label: env.name,
      value: env.id,
    })) || [];

  const assignedOptions: EnvironmentOption[] =
    assignedData?.environments?.map((env: { name: string; id: string }) => ({
      label: env.name,
      value: env.id,
    })) || [];

  const [assignEnvironment] = useAssignEnvironmentToWorkspaceMutation();
  const [unassignEnvironment] = useUnassignEnvironmentFromWorkspaceMutation();

  const handleAssignment = (added: EnvironmentOption[], removed: EnvironmentOption[]) => {
    const addedIds = added.map((item) => item.value);
    const removedIds = removed.map((item) => item.value);

    if (removedIds.length > 0) {
      removedIds.forEach((envId) => {
        const envName =
          assignedData?.environments?.find((e: { id: string; name: string }) => e.id === envId)
            ?.name || 'Unknown';
        unassignEnvironment({ workspaceId, environmentId: envId })
          .unwrap()
          .then(() => handleSuccess(`Environment "${envName}" unassigned`))
          .catch((err: any) =>
            handleError(`Environment "${envName}" Unassign Error: ${err?.data || err?.message}`),
          );
      });
      return;
    }

    if (addedIds.length > 0) {
      addedIds.forEach((envId) => {
        if (assignedOptions.find((e) => e.value === envId)) return;
        const envName =
          availableData?.environments?.find((e: { id: string; name: string }) => e.id === envId)
            ?.name || 'Unknown';
        assignEnvironment({ workspaceId, environmentId: envId })
          .unwrap()
          .then(() => handleSuccess(`Environment "${envName}" assigned`))
          .catch((err: any) =>
            handleError(`Environment "${envName}" Assign Error: ${err?.data || err?.message}`),
          );
      });
      return;
    }
  };

  return (
    <Box
      onClick={(e: React.MouseEvent) => e.stopPropagation()}
      sx={{
        marginBlock: '0.25rem',
        minWidth: '14rem',
        maxWidth: '18rem',
        width: '100%',
      }}
    >
      <Autocomplete
        multiple
        options={availableOptions}
        value={assignedOptions}
        loading={isAvailableLoading || isAssignedLoading}
        onChange={(_event, value) => {
          const current = assignedOptions || [];
          const next = (value as EnvironmentOption[]) || [];
          const added = next.filter((item) => !current.some((c) => c.value === item.value));
          const removed = current.filter((item) => !next.some((n) => n.value === item.value));
          handleAssignment(added, removed);
        }}
        size="small"
        disableCloseOnSelect
        isOptionEqualToValue={(option, val) => option.value === val.value}
        getOptionLabel={(option) => (typeof option === 'string' ? option : option.label)}
        popupIcon={<ArrowDropDownIcon fontSize="small" sx={{ color: 'inherit' }} />}
        sx={{
          '& .MuiOutlinedInput-root': {
            backgroundColor: inputBg,
            borderRadius: '4px',
            paddingRight: '36px !important',
            minHeight: '36px',
            display: 'flex',
            alignItems: 'center',
            '& fieldset': {
              borderColor: borderColor,
            },
            '&:hover fieldset': {
              borderColor: theme.palette.primary.main,
            },
            '&.Mui-focused fieldset': {
              borderColor: theme.palette.primary.main,
              borderWidth: '1px',
            },
          },
          '& .MuiInputBase-input': {
            fontSize: '0.8125rem',
            padding: '4px 6px !important',
            color: theme.palette.text.primary,
            '&::placeholder': {
              color: theme.palette.text.secondary,
              opacity: 0.85,
            },
          },
          '& .MuiAutocomplete-endAdornment': {
            position: 'absolute',
            top: '50%',
            right: '8px',
            transform: 'translateY(-50%)',
            display: 'flex',
            alignItems: 'center',
            gap: '2px',
          },
          '& .MuiAutocomplete-popupIndicator': {
            color: iconColor,
            padding: '4px',
            margin: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            transition: 'color 0.2s ease-in-out, background-color 0.2s ease-in-out',
            '&:hover': {
              color: theme.palette.text.primary,
              backgroundColor: hoverBg,
            },
          },
          '& .MuiAutocomplete-clearIndicator': {
            color: iconColor,
            padding: '4px',
            margin: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            transition: 'color 0.2s ease-in-out, background-color 0.2s ease-in-out',
            '&:hover': {
              color: theme.palette.text.primary,
              backgroundColor: hoverBg,
            },
          },
        }}
        renderValue={(selected, getItemProps) =>
          selected.map((option, index) => (
            <Chip
              label={option.label}
              size="small"
              style={{
                margin: '0.15rem',
                borderRadius: '0.25rem',
                fontSize: '0.75rem',
                height: '24px',
              }}
              {...getItemProps({ index })}
              key={option.value}
            />
          ))
        }
        renderInput={(params) => (
          <TextField
            {...params}
            placeholder={assignedOptions.length === 0 ? 'Assigned Environment' : ''}
          />
        )}
        disabled={!isAssignedEnvironmentAllowed}
      />
    </Box>
  );
};

export default WorkspaceEnvironmentSelection;

import {
  Autocomplete,
  Avatar,
  Box,
  Button,
  Checkbox,
  CircularProgress,
  FormControl,
  InputLabel,
  ListItemText,
  MenuItem,
  Select,
  TextField,
  Typography,
  Grid,
  useTheme,
  PersonIcon,
  OutlinedInput,
} from '@layer5/sistent';
import React, { useState } from 'react';
import { capitalize } from 'lodash/fp';
import { getAllUsers } from '@/rtk-query/user';
import { FileUpload } from '@mui/icons-material';
import { ImportDesignModal } from '../MesheryPatterns';
import { useNotification } from '@/utils/hooks/useNotification';
import { updateProgress } from 'lib/store';
import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { EVENT_TYPES } from 'lib/event-types';
import { useImportPatternMutation } from '@/rtk-query/design';
import SettingsIcon from '@mui/icons-material/Settings';

export const UserSearchAutoComplete = ({ handleAuthorChange }) => {
  const [open, setOpen] = React.useState(false);
  const [options, setOptions] = React.useState([]);
  const [inputValue, setInputValue] = React.useState('');
  const [loading, setLoading] = React.useState(false);

  const fetchUsers = async (search) => {
    setLoading(true);
    try {
      const { data } = await getAllUsers({ page: 0, pagesize: 10, search });
      setOptions(data.data || []);
    } catch (err) {
      console.error('Failed to fetch users', err);
      setOptions([]);
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (event, value) => {
    setInputValue(value);
    fetchUsers(value);
  };

  const handleOpen = () => {
    fetchUsers(inputValue);
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
    setOptions([]);
  };

  return (
    <Autocomplete
      open={open}
      onOpen={handleOpen}
      onClose={handleClose}
      onInputChange={handleInputChange}
      onChange={(_, value) => {
        handleAuthorChange(value?.user_id || null);
      }}
      inputValue={inputValue}
      options={options}
      loading={loading}
      isOptionEqualToValue={(option, value) => option.id === value.id}
      getOptionLabel={(option) => option.email || ''}
      renderOption={(props, option) => (
        <Box component="li" sx={{ '& > img': { mr: 2, flexShrink: 0 } }} {...props}>
          <Grid container alignItems="center">
            <Grid item>
              <Box sx={{ color: 'text.secondary', mr: 2 }}>
                <Avatar alt={option.first_name} src={option.avatar_url}>
                  {option.avatar_url ? '' : <PersonIcon />}
                </Avatar>
              </Box>
            </Grid>
            <Grid item xs>
              {option.deleted_at?.Valid ? (
                <Typography variant="body2" color="text.secondary">
                  {option.email} (deleted)
                </Typography>
              ) : (
                <>
                  <Typography variant="body2">
                    {option.first_name} {option.last_name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {option.email}
                  </Typography>
                </>
              )}
            </Grid>
          </Grid>
        </Box>
      )}
      renderInput={(params) => (
        <TextField
          {...params}
          label="Author"
          InputProps={{
            ...params.InputProps,
            endAdornment: (
              <>
                {loading ? <CircularProgress color="inherit" size={20} /> : null}
                {params.InputProps.endAdornment}
              </>
            ),
          }}
        />
      )}
    />
  );
};

export const VisibilitySelect = ({ visibility, handleVisibilityChange, visibilityItems }) => {
  return (
    <FormControl fullWidth>
      <InputLabel>Visibility</InputLabel>
      <Select
        sx={{
          '& .MuiSelect-select': {
            paddingBlock: '0.85rem',
          },
        }}
        value={visibility}
        label="Visibility"
        multiple
        input={<OutlinedInput label="Tag" />}
        renderValue={(selected) => {
          return selected.map((value) => capitalize(value)).join(', ');
        }}
        onChange={handleVisibilityChange}
      >
        {visibilityItems.map((name) => (
          <MenuItem key={name} value={name}>
            <Checkbox checked={visibility.includes(name)} />
            <ListItemText primary={capitalize(name)} />
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

export const SortBySelect = ({ sortBy, handleSortByChange }) => {
  const sortOptions = [
    { value: 'name asc', label: 'Name (A-Z)' },
    { value: 'name desc', label: 'Name (Z-A)' },
    { value: 'updated_at desc', label: 'Last Updated (Newest First)' },
    { value: 'updated_at asc', label: 'Last Updated (Oldest First)' },
    { value: 'created_at desc', label: 'Created At (Newest First)' },
    { value: 'created_at asc', label: 'Created At (Oldest First)' },
  ];

  return (
    <FormControl fullWidth>
      <InputLabel>Sort By</InputLabel>
      <Select
        sx={{
          '& .MuiSelect-select': {
            paddingBlock: '0.85rem',
          },
        }}
        value={sortBy}
        label="Sort By"
        onChange={handleSortByChange}
        input={<OutlinedInput label="Sort By" />}
      >
        {sortOptions.map((option) => (
          <MenuItem key={option.value} value={option.value}>
            {option.label}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

export const TableListHeader = () => {
  return (
    <Box display={'flex'} width={'100%'} paddingInline={'1rem'} gap={'1rem'}>
      <Typography style={{ width: '40%' }} variant="body1">
        Name
      </Typography>
      <Typography style={{ width: '30%' }} variant="body1">
        Author
      </Typography>
      <Typography style={{ width: '10%' }} variant="body1">
        Visibility
      </Typography>
      <Typography style={{ width: '20%' }} variant="body1">
        Actions
      </Typography>
    </Box>
  );
};

export const ImportButton = () => {
  const [importModal, setImportModal] = useState(false);
  const handleImportModalOpen = () => {
    setImportModal(true);
  };
  const handleImportModalClose = () => {
    setImportModal(false);
  };
  const [importPattern] = useImportPatternMutation();
  const { notify } = useNotification();
  const theme = useTheme();
  function handleImportDesign(data) {
    updateProgress({ showProgress: true });
    const { uploadType, name, url, file } = data;

    let requestBody = null;
    switch (uploadType) {
      case 'File Upload': {
        const fileElement = document.getElementById('root_file');
        const fileName = fileElement.files[0].name;
        requestBody = JSON.stringify({
          name,
          file_name: fileName,
          file: getUnit8ArrayDecodedFile(file),
        });
        break;
      }
      case 'URL Import':
        requestBody = JSON.stringify({
          url,
          name,
        });
        break;
    }

    importPattern({
      importBody: requestBody,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `"${name}" design uploaded`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        // getPattern();
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        notify({
          message: 'Error uploading design',
          event_type: EVENT_TYPES.ERROR,
        });
      });
  }
  return (
    <>
      {importModal && (
        <ImportDesignModal
          handleClose={handleImportModalClose}
          handleImportDesign={handleImportDesign}
        />
      )}
      <Button
        color="primary"
        variant="contained"
        onClick={handleImportModalOpen}
        sx={{
          minWidth: 'fit-content',
          padding: '0.85rem',
        }}
        startIcon={<FileUpload color={theme.palette.common.white} />}
      >
        Import
      </Button>
    </>
  );
};

export const AssignDesignViewButton = ({ type, handleAssign, disabled }) => {
  return (
    <Button
      disabled={disabled}
      color="primary"
      variant="contained"
      onClick={handleAssign}
      sx={{
        minWidth: 'fit-content',
        padding: '0.85rem',
      }}
      startIcon={<SettingsIcon />}
    >
      {type === 'design' ? 'Manage Designs' : 'Manage Views'}
    </Button>
  );
};

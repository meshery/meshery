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
  Grid2,
  useTheme,
  PersonIcon,
  OutlinedInput,
  FormControlLabel,
  FormGroup,
  DeleteIcon,
  ShareIcon,
  ExportIcon,
  IconButton,
  CloseIcon,
} from '@sistent/sistent';
import React, { useContext, useState } from 'react';
import { capitalize } from 'lodash/fp';
import { getAllUsers } from '@/rtk-query/user';
import { FileUpload } from '@mui/icons-material';
import { ImportDesignModal } from '../MesheryPatterns/MesheryPatterns';
import { useNotification } from '@/utils/hooks/useNotification';
import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { EVENT_TYPES } from 'lib/event-types';
import { useImportPatternMutation } from '@/rtk-query/design';
import SettingsIcon from '@mui/icons-material/Settings';
import { updateProgress } from '@/store/slices/mesheryUi';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { useAssignDesignToWorkspaceMutation } from '@/rtk-query/workspace';
import { RESOURCE_TYPE } from '@/utils/Enum';
import { iconMedium } from 'css/icons.styles';
import MoveFileIcon from '@/assets/icons/MoveFileIcon';
import { StyledMuiDoubleCheckbox, StyledResponsiveButton } from './styles';

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
          <Grid2 container alignItems="center" size="grow">
            <Grid2>
              <Box sx={{ color: 'text.secondary', mr: 2 }}>
                <Avatar alt={option.first_name} src={option.avatar_url}>
                  {option.avatar_url ? '' : <PersonIcon />}
                </Avatar>
              </Box>
            </Grid2>
            <Grid2 size={{ sx: 12 }}>
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
            </Grid2>
          </Grid2>
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
    { value: 'name asc', label: 'A to Z' },
    { value: 'name desc', label: 'Z to A' },
    { value: 'updated_at desc', label: 'Latest Update' },
    { value: 'updated_at asc', label: 'Oldest Update' },
    { value: 'created_at desc', label: 'Newest' },
    { value: 'created_at asc', label: 'Oldest' },
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

export const TableListHeader = ({ content = [], isMultiSelectMode = false }) => {
  const { setMultiSelectedContent, multiSelectedContent } = useContext(WorkspaceModalContext);
  const theme = useTheme();
  return (
    <Grid2
      container
      width="100%"
      size="grow"
      paddingInline="1rem"
      spacing={2}
      alignItems="center"
      wrap="nowrap"
    >
      {isMultiSelectMode && (
        <Grid2 size={{ xs: 1, md: 0.5, lg: 0.25 }}>
          <FormGroup>
            <FormControlLabel
              control={
                <StyledMuiDoubleCheckbox
                  checked={
                    multiSelectedContent.length != 0 &&
                    multiSelectedContent.length === content.length
                  }
                  onChange={(e) => {
                    if (e.target.checked) {
                      setMultiSelectedContent(content);
                    } else {
                      setMultiSelectedContent([]);
                    }
                  }}
                />
              }
            />
          </FormGroup>
        </Grid2>
      )}
      <Grid2 size={{ xs: 6, md: 3.5, lg: 3 }}>
        <Typography variant="body1" noWrap>
          Name
        </Typography>
      </Grid2>
      <Grid2 size={{ xs: 4, md: 4, lg: isMultiSelectMode ? 2.75 : 3 }}>
        <Typography sx={{
          [theme.breakpoints.down('sm')]: {
            display: "flex",
            justifyContent: "center",
          }
        }} variant="body1" noWrap>
          Author
        </Typography>
      </Grid2>
      <Grid2 size={{ md: 2, lg: 1.5 }} sx={{ display: { xs: 'none', md: 'block' } }}>
        <Typography>Organization</Typography>
      </Grid2>
      <Grid2 size={{ lg: 1.5 }} sx={{ display: { xs: 'none', lg: 'block' } }}>
        <Typography>Workspace</Typography>
      </Grid2>

      <Grid2 size={{ md: 1, lg: 1 }} sx={{ display: { xs: 'none', md: 'block' } }}>
        <Typography variant="body1" noWrap>
          Visibility
        </Typography>
      </Grid2>
      <Grid2
        size={{ xs: 2, md: 1, lg: 2 }}
        sx={{
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
        }}
      >
        <Typography variant="body1" noWrap>
          Actions
        </Typography>
      </Grid2>
    </Grid2>
  );
};

export const ImportButton = ({ workspaceId, disabled = false, refetch }) => {
  const [importModal, setImportModal] = useState(false);
  const handleImportModalOpen = () => {
    setImportModal(true);
  };
  const [assignDesignToWorkspace] = useAssignDesignToWorkspaceMutation();
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
      .then((data) => {
        updateProgress({ showProgress: false });
        notify({
          message: `"${name}" design uploaded`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        if (workspaceId) {
          assignDesignToWorkspace({
            workspaceId: workspaceId,
            designId: data[0].id,
          });
        }
        handleImportModalClose();
        if (refetch) {
          refetch();
        }
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
      <StyledResponsiveButton
        color="primary"
        variant="contained"
        onClick={handleImportModalOpen}
        disabled={disabled}
        sx={{
          minWidth: 'fit-content',
          padding: '0.85rem !important',
          width: '100%',
        }}
        startIcon={<FileUpload color={theme.palette.common.white} />}
      >
        <Box sx={{ display: { xs: 'none', sm: 'block' } }}>Import</Box>
      </StyledResponsiveButton>
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
      {type === RESOURCE_TYPE.DESIGN ? 'Manage Designs' : 'Manage Views'}
    </Button>
  );
};

export const MultiContentSelectToolbar = ({
  type,
  handleContentMove,
  handleDownload,
  handleViewDownload,
  handleDelete,
  handleShare,
  refetch,
}) => {
  const theme = useTheme();
  const { multiSelectedContent, setMultiSelectedContent } = useContext(WorkspaceModalContext);
  return (
    <>
      {multiSelectedContent.length > 0 && (
        <Box
          width={'100%'}
          sx={{ backgroundColor: theme.palette.background.default }}
          height={'4rem'}
          borderRadius={'0.5rem'}
          display={'flex'}
          justifyContent={'space-between'}
          alignItems={'center'}
          paddingInline={'1rem'}
        >
          <Box display={'flex'} alignItems={'center'} gap={'0.5rem'}>
            <IconButton onClick={() => setMultiSelectedContent([])}>
              <CloseIcon />
            </IconButton>
            <Typography>
              {multiSelectedContent.length} {type} selected
            </Typography>
          </Box>
          <Box style={{ display: 'flex', gap: '0.5rem' }}>
            {handleContentMove && (
              <StyledResponsiveButton
                variant="contained"
                startIcon={<MoveFileIcon style={iconMedium} fill={theme.palette.common.white} />}
                onClick={() => {
                  handleContentMove(true);
                }}
                disabled={!multiSelectedContent.length}
              >
                <Box sx={{ display: { xs: 'none', sm: 'block' } }}>Move</Box>
              </StyledResponsiveButton>
            )}
            <StyledResponsiveButton
              variant="contained"
              startIcon={<ExportIcon style={iconMedium} fill={theme.palette.common.white} />}
              onClick={() => {
                type === RESOURCE_TYPE.DESIGN
                  ? handleDownload(multiSelectedContent)
                  : handleViewDownload(multiSelectedContent);
                setMultiSelectedContent([]);
              }}
              disabled={!multiSelectedContent.length}
            >
              <Box sx={{ display: { xs: 'none', sm: 'block' } }}>Download</Box>
            </StyledResponsiveButton>{' '}
            {handleShare && (
              <StyledResponsiveButton
                variant="contained"
                startIcon={<ShareIcon style={iconMedium} fill={theme.palette.common.white} />}
                onClick={() => {
                  handleShare(multiSelectedContent);
                  setMultiSelectedContent([]);
                }}
                disabled={!multiSelectedContent.length}
              >
                <Box sx={{ display: { xs: 'none', sm: 'block' } }}>Share</Box>
              </StyledResponsiveButton>
            )}
            <StyledResponsiveButton
              color="error"
              variant="contained"
              onClick={() => {
                handleDelete(
                  multiSelectedContent,
                  type === RESOURCE_TYPE.DESIGN ? RESOURCE_TYPE.DESIGN : RESOURCE_TYPE.VIEW,
                  refetch,
                );
                setMultiSelectedContent([]);
              }}
              sx={{
                backgroundColor: `${theme.palette.error.dark} !important`,
              }}
              startIcon={<DeleteIcon style={iconMedium} fill={theme.palette.common.white} />}
            >
              <Box sx={{ display: { xs: 'none', sm: 'block' } }}>Delete</Box>
            </StyledResponsiveButton>
          </Box>
        </Box>
      )}
    </>
  );
};

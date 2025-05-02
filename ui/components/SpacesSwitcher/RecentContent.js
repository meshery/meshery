//@ts-check
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import {
  Autocomplete,
  Avatar,
  Box,
  Checkbox,
  CircularProgress,
  DesignIcon,
  FormControl,
  InputLabel,
  ListItemText,
  MenuItem,
  Select,
  Tab,
  Tabs,
  TextField,
  Typography,
  useTheme,
  ViewIcon,
} from '@layer5/sistent';
import { iconMedium } from 'css/icons.styles';
import React, { useState } from 'react';
import DesignsContent from './DesignsContent';
import ViewsContent from './ViewsContent';
import { StyledSearchBar } from '@layer5/sistent';
import MainDesignsContent from './MainDesignsContent';
import { useGetUserDesignsQuery } from '@/rtk-query/design';
import MainViewsContent from './MainViewsContent';
import { useFetchViewsQuery } from '@/rtk-query/view';
import { VISIBILITY } from '@/utils/Enum';
import { OutlinedInput } from '@layer5/sistent';
import { capitalize } from 'lodash/fp';
import { getAllUsers } from '@/rtk-query/user';
import { Grid } from '@layer5/sistent';
import { PersonIcon } from '@layer5/sistent';

const RecentContent = () => {
  const isViewVisible = CAN(keys.VIEW_VIEWS.action, keys.VIEW_VIEWS.subject);
  const isDesignsVisible = CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject);

  const [searchQuery, setSearchQuery] = useState('');
  const visibilityItems = [VISIBILITY.PUBLIC, VISIBILITY.PRIVATE];

  const [type, setType] = React.useState('design');
  const [author, setAuthor] = React.useState('');
  const [modified, setModified] = useState('updated_at desc');
  const [visibility, setVisibility] = useState(visibilityItems);
  const handleTypeChange = (event) => {
    setType(event.target.value);
    setDesignsPage(0);
    setViewsPage(0);
  };
  const handleAuthorChange = (user_id) => {
    setDesignsPage(0);
    setViewsPage(0);
    setAuthor(user_id);
  };
  const handleModifiedChange = (event) => {
    setDesignsPage(0);
    setViewsPage(0);
    setModified(event.target.value);
  };
  const handleVisibilityChange = (event) => {
    const value = event.target.value;
    setVisibility(typeof value === 'string' ? value.split(',') : value);
    setDesignsPage(0);
    setViewsPage(0);
  };
  const onSearchChange = (e) => {
    setDesignsPage(0);
    setViewsPage(0);
    setSearchQuery(e.target.value);
  };

  const menuItems = [
    { value: '', label: 'Sort by', disabled: true },
    { value: 'name asc', label: 'Alphabetically (A-Z)' },
    { value: 'name desc', label: 'Alphabetically (Z-A)' },
    { value: 'updated_at desc', label: 'Most Recently Updated' },
    { value: 'updated_at asc', label: 'Least Recently Updated' },
  ];
  const [designsPage, setDesignsPage] = useState(0);
  const [viewsPage, setViewsPage] = useState(0);
  const {
    data: designsData,
    isLoading,
    isFetching,
  } = useGetUserDesignsQuery(
    {
      expandUser: true,
      page: designsPage,
      pagesize: 10,
      order: modified,
      metrics: true,
      search: searchQuery,
      visibility: visibility,
      user_id: author,
    },
    {
      skip: type !== 'design',
    },
  );

  const {
    data: viewsData,
    isLoading: isViewLoading,
    isFetching: isViewFetching,
  } = useFetchViewsQuery(
    {
      page: viewsPage,
      pagesize: 10,
      order: modified,
      user_id: author,
      visibility: visibility,
      search: searchQuery,
    },
    {
      skip: type !== 'view',
    },
  );

  const theme = useTheme();
  return (
    <>
      <Box style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
        <StyledSearchBar
          sx={{
            backgroundColor: 'transparent',
          }}
          width="auto"
          placeholder={type === 'design' ? 'Search Designs' : 'Search Views'}
          value={searchQuery}
          onChange={onSearchChange}
          endAdornment={
            type === 'design' ? (
              <p style={{ color: theme.palette.text.default }}>
                Total Designs: {designsData?.total_count ?? 0}
              </p>
            ) : (
              <p style={{ color: theme.palette.text.default }}>
                Total Views: {viewsData?.total_count ?? 0}
              </p>
            )
          }
        />
        <Box display={'flex'} alignItems="center" marginBottom="1rem" gap={'1rem'}>
          <Box sx={{ minWidth: 120 }}>
            <FormControl fullWidth>
              <InputLabel>Type</InputLabel>
              <Select value={type} label="Type" onChange={handleTypeChange}>
                <MenuItem value={'design'}>Design</MenuItem>
                <MenuItem value={'view'}>View</MenuItem>
              </Select>
            </FormControl>
          </Box>

          <Box sx={{ minWidth: 120 }}>
            <FormControl fullWidth>
              <InputLabel>Sort By</InputLabel>
              <Select value={modified} label="Sort By" onChange={handleModifiedChange}>
                {menuItems.map((item) => (
                  <MenuItem key={item.value} value={item.value} disabled={item.disabled}>
                    {item.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
          <Box sx={{ minWidth: 300 }}>
            <FormControl fullWidth>
              <UserAutoComplete handleAuthorChange={handleAuthorChange} />
            </FormControl>
          </Box>
          <Box sx={{ minWidth: 120 }}>
            <FormControl fullWidth>
              <InputLabel>Visibility</InputLabel>
              <Select
                value={visibility}
                label="Visibility"
                multiple
                input={<OutlinedInput label="Tag" />}
                //make the first letter capital
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
          </Box>
        </Box>
        <Box display={'flex'} width={'100%'} paddingInline={'1rem'} gap={'1rem'}>
          <Typography style={{ width: '50%' }} variant="body1">
            Name
          </Typography>
          <Typography style={{ width: '30%' }} variant="body1">
            Author
          </Typography>
          <Typography style={{ width: '10%' }} variant="body1">
            Visibility
          </Typography>
          <Typography style={{ width: '10%' }} variant="body1">
            Actions
          </Typography>
        </Box>

        {type == 'design' && (
          <MainDesignsContent
            designsData={designsData}
            setPage={setDesignsPage}
            isLoading={isLoading}
            isFetching={isFetching}
            designs={designsData?.patterns}
          />
        )}
        {type == 'view' && (
          <MainViewsContent
            viewsData={viewsData}
            setPage={setViewsPage}
            isLoading={isViewLoading}
            isFetching={isViewFetching}
            views={viewsData?.views}
          />
        )}
      </Box>
    </>
  );
};

export default RecentContent;

const UserAutoComplete = ({ handleAuthorChange }) => {
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

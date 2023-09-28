import * as React from 'react';
import TextField from '@mui/material/TextField';
import Autocomplete, { createFilterOptions } from '@mui/material/Autocomplete';
import AddIcon from '@mui/icons-material/Add';
import { withStyles } from '@material-ui/core/styles';

const filter = createFilterOptions();

const styles = (theme) => ({
  heading: {
    color: theme.palette.secondary.iconMain,
    '& .MuiInputLabel-root': {
      color: theme.palette.secondary.iconMain,
    },
    '& .MuiOutlinedInput-root': {
      color: theme.palette.secondary.iconMain,
    },
    '& .MuiSvgIcon-root': {
      color: theme.palette.secondary.iconMain,
    },
  },
  wrapper: {
    background:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    padding: '8px',
  },
  autocompleteOptions: {
    '& .MuiAutocomplete-option': {
      background: theme.palette.secondary.mainBackground,
    },
  },
  dialogTitle: {
    textAlign: 'center',
    minWidth: 400,
    padding: '10px',
    color: '#fff',
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.headerColor
        : theme.palette.secondary.mainBackground,
  },
  subtitle: {
    minWidth: 400,
    overflowWrap: 'anywhere',
    textAlign: 'center',
    padding: '5px',
  },
  button: {
    backgroundColor: theme.palette.secondary.focused,
    '&:hover': {
      backgroundColor: theme.palette.secondary.focused,
    },
    color: '#fff',
  },
  icon: {
    color: theme.palette.secondary.iconMain,
  },
});

function MySelectComponent({ classes }) {
  const [value, setValue] = React.useState(null);
  const [labelVisible, setLabelVisible] = React.useState(true);

  React.useEffect(() => {
    if (value) {
      setLabelVisible(false);
    }
  }, [value]);

  return (
    <React.Fragment>
      <Autocomplete
        className={classes.autocompleteOptions}
        value={value}
        onChange={(event, newValue) => {
          if (typeof newValue === 'string') {
            setValue({
              title: newValue,
            });
          } else if (newValue && newValue.inputValue) {
            // Create a new value from the user input
            setValue({
              title: newValue.inputValue,
            });
          } else {
            setValue(newValue);
          }
        }}
        filterOptions={(options, params) => {
          const filtered = filter(options, params);

          if (params.inputValue !== '') {
            filtered.push({
              inputValue: params.inputValue,
              title: `Create:  "${params.inputValue}"`,
            });
          }

          return filtered;
        }}
        id="free-solo-dialog-demo"
        options={environmentOptions}
        getOptionLabel={(option) => {
          // e.g. value selected with enter, right from the input
          if (typeof option === 'string') {
            return option;
          }
          if (option.inputValue) {
            return option.inputValue;
          }
          return option.title;
        }}
        selectOnFocus
        clearOnBlur
        handleHomeEndKeys
        renderOption={(props, option) => (
          <div className={classes.wrapper}>
            <li
              className={classes.wrapper}
              {...props}
              style={{ display: 'flex', justifyContent: 'space-between' }}
            >
              <div className={classes.heading}>
                {option.inputValue ? <>{option.title}</> : option.title}
              </div>
              <div>{option.inputValue && <AddIcon className={classes.icon} />}</div>
            </li>
          </div>
        )}
        sx={{ '& fieldset': { border: 'none' }, width: 250 }}
        renderInput={(params) => (
          <TextField
            className={classes.heading}
            {...params}
            label={labelVisible ? 'Choose an environment' : ''}
            placeholder="Environment"
          />
        )}
      />
    </React.Fragment>
  );
}

const environmentOptions = [
  { title: 'Environment 1' },
  { title: 'Environment 2' },
  { title: 'Environment 3' },
  { title: 'Environment 4' },
  { title: 'Environment 5' },
];

export default withStyles(styles)(MySelectComponent);

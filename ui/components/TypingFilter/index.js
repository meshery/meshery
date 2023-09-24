import {
  ClickAwayListener,
  Divider,
  Fade,
  IconButton,
  InputAdornment,
  List,
  Popper,
  TextField,
  Typography,
  useTheme,
} from '@material-ui/core';
import ContentFilterIcon from '../../assets/icons/ContentFilterIcon';
import { useEffect, useReducer, useRef, useState } from 'react';
import CrossCircleIcon from '../../assets/icons/CrossCircleIcon';
import clsx from 'clsx';
import { useStyles, useFilterStyles } from './style';
import { FILTERING_STATE, FILTER_EVENTS, filterReducer } from './state';
import { getFilters, getCurrentFilterAndValue } from './utils';

const Filters = ({ filterStateMachine, dispatchFilterMachine, filterSchema }) => {
  const classes = useFilterStyles();
  const selectFilter = (filter) => {
    dispatchFilterMachine({
      type: FILTER_EVENTS.SELECT,
      payload: {
        value: filter,
      },
    });
  };

  const { filter: currentFilter } = getCurrentFilterAndValue(filterStateMachine);
  const matchingFilters = currentFilter
    ? Object.values(filterSchema).filter((filter) => filter.value.startsWith(currentFilter))
    : Object.values(filterSchema);
  return (
    <List>
      {matchingFilters.length == 0 && (
        <div className={classes.item}>
          <Typography variant="body1" className={classes.label}>
            Sorry we dont currently support this filter
          </Typography>
        </div>
      )}
      {matchingFilters.map((filter) => {
        return (
          <>
            <div
              key={filter.value}
              className={classes.item}
              disableGutters
              onClick={() => selectFilter(filter.value)}
            >
              <Typography variant="body1" className={classes.label}>
                {filter.value}:
              </Typography>
              <Typography variant="body1" className={classes.description}>
                {filter.description}
              </Typography>
            </div>
            <Divider light />
          </>
        );
      })}
    </List>
  );
};

const FilterValueSuggestions = ({ filterStateMachine, dispatchFilterMachine, filterSchema }) => {
  const classes = useFilterStyles();

  const selectValue = (value) => {
    dispatchFilterMachine({
      type: FILTER_EVENTS.SELECT,
      payload: {
        value,
      },
    });
  };
  const { filter, value } = getCurrentFilterAndValue(filterStateMachine);
  const currentFilter = Object.values(filterSchema).find((f) => f.value == filter);
  const suggestions = currentFilter?.values?.filter((v) => v.startsWith(value)) || [];

  return (
    <List>
      {suggestions.length == 0 && (
        <div className={classes.item}>
          <Typography variant="body1" className={classes.label}>
            No results available
          </Typography>
        </div>
      )}
      {suggestions.map((value) => {
        return (
          <>
            <div
              key={value.value}
              className={classes.item}
              disableGutters
              onClick={() => selectValue(value)}
            >
              <Typography variant="body1" className={classes.description}>
                {value}
              </Typography>
            </div>
            <Divider light />
          </>
        );
      })}
    </List>
  );
};

/**
 * Filter Schema Object
 *
 * The `filterSchema` object defines available filter options for the TypingFilter component.
 * It provides information about different filter categories, their descriptions, and possible values.
 *
 * @typedef {object} FilterSchema
 * @property {object} [CATEGORY_NAME] - An object representing a filter category.
 * @property {string} [CATEGORY_NAME.value] - The filter name used for filtering within the category.
 * @property {string} [CATEGORY_NAME.description] - Description of the filter category.
 * @property {string} [CATEGORY_NAME.type] - The data type of the filter (optional).
 * @property {string[]} [CATEGORY_NAME.values] - Possible values for the filter (optional).
 *
 * @example
 * // Example filter schema with multiple filter categories
 * const filterSchema = {
 *   SEVERITY: {
 *     value: "severity",
 *     description: "Filter by severity",
 *     values: ["Low", "Medium", "High"],
 *   },
 *   STATUS: {
 *     value: "status",
 *     description: "Filter by status",
 *     type: "string",
 *     values: ["Open", "Closed", "In Progress"],
 *   },
 *   CUSTOM_FILTER: {
 *     value: "custom",
 *     description: "Custom filter description",
 *     type: "number",
 *   },
 *   // Add more filter categories as needed
 * };
 */

/**
 * TypingFilter Component
 *
 * A component for real-time filtering and selection with typing. It provides a user-friendly
 * interface for filtering data based on user input.
 *
 * @component
 * @param {object} props - Component props.
 * @param {FilterSchema} filterSchema - The schema defining available filter options.
 * @param {function} handleFilter - A callback function to handle filter changes.
 * @param {boolean} autoFilter - A boolean to indicate if the filter should be applied automatically (on user input) .
 * @returns {JSX.Element} - A React JSX element representing the TypingFilter component.
 */
const TypingFilter = ({ filterSchema, handleFilter , autoFilter=false }) => {
  const theme = useTheme();
  // console.log("initialFilter", initialFilter)
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = useState(null);
  const isPopperOpen = Boolean(anchorEl);
  const inputFieldRef = useRef(null);
  const [filteringState, dispatch] = useReducer(filterReducer, {
    context: {
      value: '',
      prevValue: [''],
    },
    state: FILTERING_STATE.IDLE,
  });


  const handleFilterChange = (e) => {
    if (!anchorEl) {
      setAnchorEl(e.currentTarget);
    }

    if (e.target.value === '') {
      return dispatch({
        type: FILTER_EVENTS.CLEAR,
      });
    }

    dispatch({
      type : FILTER_EVENTS.INPUT_CHANGE,
      payload : {
        value : e.target.value,
      },
    });

  };

  const handleClear = () => {
    dispatch({
      type: FILTER_EVENTS.EXIT,
    });

    handleFilter({});
  };

  const handleFocus = (e) => {
    setAnchorEl(e.currentTarget);
    dispatch({ type: 'START' });
  };

  const handleClickAway = (e) => {
    if (inputFieldRef.current.contains(e.target)) {
      return;
    }

    setAnchorEl(null);
  };

  //add enter event listener to the input fieldse
  //add esc event listener to the input fields
  useEffect(() => {
    if (!inputFieldRef.current) {
      return;
    }

    const handleKeyDown = (e) => {
      if (e.key == 'Enter') {
        handleFilter(getFilters(e.target.value, filterSchema));
        setAnchorEl(null);
      }
    };
    inputFieldRef?.current?.addEventListener('keydown', handleKeyDown);
    return () => {
      inputFieldRef?.current?.removeEventListener("keydown", handleKeyDown)
    }
  }, [inputFieldRef.current])


  useEffect(() => {
    if (autoFilter && filteringState.state == FILTERING_STATE.SELECTING_FILTER) {
      handleFilter(getFilters(filteringState.context.value, filterSchema))
    }
  }, [filteringState.state])


  return (
    <div className={clsx(classes.root, 'mui-fixed')}>
      <TextField
        ref={inputFieldRef}
        variant="outlined"
        placeholder="Filter Notifications"
        fullWidth
        size="small"
        className={classes.input}
        value={filteringState.context.value}
        onChange={handleFilterChange}
        onFocus={handleFocus}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              {' '}
              <ContentFilterIcon fill={theme.palette.secondary.iconMain} />{' '}
            </InputAdornment>
          ),
          endAdornment: (
            <InputAdornment position="end">
              <IconButton onClick={handleClear}>
                {filteringState.state !== FILTERING_STATE.IDLE && (
                  <CrossCircleIcon fill={theme.palette.secondary.iconMain} />
                )}
              </IconButton>
            </InputAdornment>
          ),
        }}
      />

      <Popper
        open={filteringState.state != FILTERING_STATE.IDLE && isPopperOpen}
        anchorEl={inputFieldRef.current}
        placement="bottom-start"
        style={{ zIndex: 2000 }}
        transition
        className="mui-fixed"
      >
        {({ TransitionProps }) => {
          return (
            <Fade {...TransitionProps} timeout={100}>
              <ClickAwayListener onKeydown onClickAway={handleClickAway}>
                <div
                  className={classes.dropDown}
                  style={{
                    width: inputFieldRef.current ? inputFieldRef.current.clientWidth : 0,
                  }}
                >
                  {filteringState.state == FILTERING_STATE.SELECTING_FILTER && (
                    <Filters
                      filterStateMachine={filteringState}
                      dispatchFilterMachine={dispatch}
                      filterSchema={filterSchema}
                    />
                  )}
                  {filteringState.state == FILTERING_STATE.SELECTING_VALUE && (
                    <FilterValueSuggestions
                      filterStateMachine={filteringState}
                      dispatchFilterMachine={dispatch}
                      filterSchema={filterSchema}
                    />
                  )}
                </div>
              </ClickAwayListener>
            </Fade>
          );
        }}
      </Popper>
    </div>
  );
};

export default TypingFilter;

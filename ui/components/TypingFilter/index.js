import {
  ClickAwayListener,
  Divider,
  IconButton,
  InputAdornment,
  List,
  Popper,
  useTheme,
  Fade,
} from '@layer5/sistent';
import { Description, DropDown, InputField, Item, Label, Root } from './style';
import { UsesSistent } from '../SistentWrapper';
import ContentFilterIcon from '../../assets/icons/ContentFilterIcon';
import { useEffect, useReducer, useRef, useState } from 'react';
import CrossCircleIcon from '../../assets/icons/CrossCircleIcon';

import { FILTERING_STATE, FILTER_EVENTS, filterReducer } from './state';
import { getFilters, getCurrentFilterAndValue } from './utils';
import { ExpandLess, ExpandMore } from '@mui/icons-material';

const Filters = ({ filterStateMachine, dispatchFilterMachine, filterSchema }) => {
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
        <Item>
          <Label variant="body1">Sorry we dont currently support this filter</Label>
        </Item>
      )}
      {matchingFilters.map((filter) => {
        return (
          <>
            <Item key={filter.value} disableGutters onClick={() => selectFilter(filter.value)}>
              <Label variant="body1">{filter.value}:</Label>
              <Description variant="body1">{filter.description}</Description>
            </Item>
            <Divider light />
          </>
        );
      })}
    </List>
  );
};

const FilterValueSuggestions = ({ filterStateMachine, dispatchFilterMachine, filterSchema }) => {
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
        <Item>
          <Label variant="body1">No results available</Label>
        </Item>
      )}
      {suggestions.map((value) => {
        return (
          <>
            <Item key={value.value} disableGutters onClick={() => selectValue(value)}>
              <Description variant="body1">{value}</Description>
            </Item>
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
const TypingFilter = ({ filterSchema, handleFilter, autoFilter = false, placeholder }) => {
  const theme = useTheme();
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
      type: FILTER_EVENTS.INPUT_CHANGE,
      payload: {
        value: e.target.value,
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
      if (e.key == 'Escape') {
        setAnchorEl(null);
        //remove focus from the input field
        setTimeout(() => {
          document.activeElement.blur();
        });
      }
    };

    inputFieldRef?.current?.addEventListener('keydown', handleKeyDown);
    return () => {
      inputFieldRef?.current?.removeEventListener('keydown', handleKeyDown);
    };
  }, [inputFieldRef.current]);

  const toggleSuggestions = () => {
    if (isPopperOpen) {
      setAnchorEl(null);
    } else {
      setAnchorEl(inputFieldRef.current);
    }
  };

  useEffect(() => {
    if (autoFilter && filteringState.state == FILTERING_STATE.SELECTING_FILTER) {
      handleFilter(getFilters(filteringState.context.value, filterSchema));
    }
  }, [filteringState.state]);

  return (
    <UsesSistent>
      <Root className="mui-fixed">
        <InputField
          ref={inputFieldRef}
          variant="outlined"
          placeholder={placeholder}
          fullWidth
          size="small"
          value={filteringState.context.value}
          onChange={handleFilterChange}
          onFocus={handleFocus}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                {' '}
                <ContentFilterIcon fill={theme.palette.icon.default} />{' '}
              </InputAdornment>
            ),
            endAdornment: (
              <InputAdornment position="end">
                <IconButton onClick={handleClear}>
                  {filteringState.state !== FILTERING_STATE.IDLE && (
                    <CrossCircleIcon fill={theme.palette.icon.default} />
                  )}
                </IconButton>
                <IconButton onClick={toggleSuggestions}>
                  {filteringState.state !== FILTERING_STATE.IDLE && !anchorEl && (
                    <ExpandMore fill={theme.palette.icon.default} />
                  )}
                  {filteringState.state !== FILTERING_STATE.IDLE && anchorEl && (
                    <ExpandLess fill={theme.palette.icon.default} />
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
                <div>
                  <ClickAwayListener onKeydown onClickAway={handleClickAway}>
                    <DropDown
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
                    </DropDown>
                  </ClickAwayListener>
                </div>
              </Fade>
            );
          }}
        </Popper>
      </Root>
    </UsesSistent>
  );
};

export default TypingFilter;

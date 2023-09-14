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
  alpha,
  makeStyles,
  useTheme,
} from "@material-ui/core";
import ContentFilterIcon from "../../assets/icons/ContentFilterIcon";
import { useEffect, useReducer, useRef, useState } from "react";
import CrossCircleIcon from "../../assets/icons/CrossCircleIcon";
import clsx from "clsx";
import { SEVERITY, STATUS } from "./constants";

const useStyles = makeStyles((theme) => ({
  root : {
    position : "relative",
    backgroundColor : theme.palette.secondary.elevatedComponents,
  },
  input : {
    width : "100%",
    "& .MuiOutlinedInput-root" : {
      borderRadius : "6px",
      backgroundColor : theme.palette.secondary.searchBackground,
      "& fieldset" : {
        borderRadius : "6px",
        border : `2px solid ${theme.palette.secondary.searchBorder}`,
      },
    },
  },

  dropDown : {
    backgroundColor : theme.palette.secondary.searchBackground,
    borderRadius : "6px",
    boxShadow :
      "0px 2px 4px 0px rgba(0, 0, 0, 0.20), 0px 1px 10px 0px rgba(0, 0, 0, 0.12), 0px 4px 5px 0px rgba(0, 0, 0, 0.14)",
    border : `2px solid ${theme.palette.secondary.searchBorder}`,
    marginTop : "0.2rem",
  },
}));

const useFilterStyles = makeStyles((theme) => ({
  item : {
    fontFamily : "Qanelas Soft, sans-serif",
    display : "flex",
    gap : "0.3rem",
    margin : "0.3rem",
    padding : "0.3rem",
    paddingInline : "1rem",
    borderRadius : "6px",
    cursor : "pointer",
    "&:hover" : {
      backgroundColor : alpha(theme.palette.secondary.link2, 0.25),
    },
  },

  label : {
    fontWeight : 500,
    color : theme.palette.secondary.icon,
  },
  description : {
    fontWeight : 400,
    color : theme.palette.secondary.number,
  },
}));

const FILTERS = {
  SEVERITY : {
    value : "severity",
    label : "Severity",
    description : "Filter by severity",
    values : Object.values(SEVERITY),
  },

  STATUS : {
    value : "status",
    label : "Status",
    description : "Filter by status",
    values : Object.values(STATUS),
    type : "string"
  },

  TYPE : {
    value : "type",
    label : "Type",
    description : "Filter by type",
  },

  AUTHOR : {
    value : "author",
    label : "Author",
    description : "Filter by any user or system",
  },

  CATEGORY : {
    value : "category",
    label : "Category",
    description : "Filter by category",
    values : ["pattern", "connection"],
  },
};

const FILTERING_STATE = {
  IDLE : "idle",
  SELECTING_FILTER : "selecting_filter",
  SELECTING_VALUE : "selecting_value",
};

const FILTER_EVENTS = {
  START : "start",
  SELECT : "select_filter",
  SELECT_FILTER : "select_filter",
  INPUT_CHANGE : "input_change",
  SELECT_FILTER_VALUE : "select_filter_value",
  CLEAR : "clear",
  EXIT : "exit",
};

const Delimiter = {
  FILTER : " ",
  FILTER_VALUE : ":",
};

/**
 * Parses a filter string and returns a filter object.
 *
 * @param {string} filterString - The input filter string of the form "type:value type2:value2 type:value2".
 * @returns {Object} - The filter object with types as keys and arrays of values as values.
 */
const getFilters = (filterString) => {
  const filters = {};
  const filterValuePairs = filterString.split(Delimiter.FILTER);
  filterValuePairs.forEach((filterValuePair) => {
    const [filter, value] = filterValuePair.split(Delimiter.FILTER_VALUE);

    if (filter == FILTERS.STATUS.value) {
      filters[filter] = value;
      return
    }

    if (filter && value) {
      filters[filter] = filters[filter] || [];
      if (!filters[filter].includes(value)) {
        filters[filter].push(value)
      }
    }
  });


  return filters;
};

// return a filter string of form "type:value type2:value2 type:value2"
// from a filter object of form { type : {values} , type2 : {values}  }
export const getFilterString = (filters) => {
  return Object.entries(filters).reduce((filterString, [filter, values]) => {
    return filterString + [...values].map((value) => `${filter}${Delimiter.FILTER_VALUE}${value}`).join(" ");
  }, "");
};



const commonReducer = (stateMachine, action) => {
  const { context } = stateMachine;
  switch (action.type) {
    case FILTER_EVENTS.CLEAR:
      return {
        state : FILTERING_STATE.SELECTING_FILTER,
        context : {
          ...context,
          value : "",
          prevValue : [""],
        },
      };

    case FILTER_EVENTS.EXIT:
      return {
        state : FILTERING_STATE.IDLE,
        context : {
          ...context,
          value : "",
          prevValue : [""],
        },
      };

    default:
      return stateMachine;
  }
};

const filterSelectionReducer = (stateMachine, action, nextState, nextValue) => {
  const { state, context } = stateMachine;
  const nextDelimiter = nextState == FILTERING_STATE.SELECTING_FILTER ? Delimiter.FILTER : Delimiter.FILTER_VALUE;
  const prevDelimiter = nextDelimiter == Delimiter.FILTER_VALUE ? Delimiter.FILTER : Delimiter.FILTER_VALUE;
  const prevState = nextState; // same beccuase the prevState is the same as the nextState ( as we have only two states)
  switch (action.type) {
    // Select a filter and move to start entring its value
    case FILTER_EVENTS.SELECT: {
      const newValue = nextValue(context.prevValue.at(-1), action.payload.value); // ":" is used to separate the filter and its value)
      return {
        state : nextState,
        context : {
          ...context,
          value : newValue + nextDelimiter,
          prevValue : [...context.prevValue, newValue],
        },
      };
    }
    //" " is used to separate multiple filters
    case FILTER_EVENTS.INPUT_CHANGE:
      // prevent transition when the the filter/value is empty
      if (action.payload.value.endsWith(nextDelimiter) && context.value.endsWith(prevDelimiter)) {
        return stateMachine;
      }

      // prevent adding multiple delimeters together
      if (action.payload.value.endsWith(prevDelimiter) && context.value.endsWith(prevDelimiter)) {
        return stateMachine;
      }

      if (action.payload.value == context.prevValue.at(-1)) {
        return {
          state : prevState,
          context : {
            ...context,
            prevValue : context.prevValue.slice(0, -1),
            value : action.payload.value,
          },
        };
      }

      if (action.payload.value.endsWith(nextDelimiter)) {
        const newValue = action.payload.value;
        return {
          state : nextState,
          context : {
            ...context,
            value : action.payload.value,
            prevValue : [...context.prevValue, newValue.slice(0, -1)],
          },
        };
      }

      return {
        state, // stay in the same state
        context : {
          ...context,
          value : action.payload.value,
        },
      };
    default:
      return commonReducer(stateMachine, action);
  }
};

const filterReducer = (stateMachine, action) => {
  const { state } = stateMachine;
  switch (state) {
    // Initial State
    case FILTERING_STATE.IDLE:
      switch (action.type) {
        // Start the filter process
        case "START":
          return {
            ...stateMachine,
            state : FILTERING_STATE.SELECTING_FILTER,
          };
        default:
          return stateMachine;
      }

    case FILTERING_STATE.SELECTING_FILTER:
      // return filterSelectionReducer(stateMachine, action);
      return filterSelectionReducer(
        stateMachine,
        action,
        FILTERING_STATE.SELECTING_VALUE,
        (prevValue, value) => prevValue + Delimiter.FILTER + value
      );

    case FILTERING_STATE.SELECTING_VALUE:
      return filterSelectionReducer(
        stateMachine,
        action,
        FILTERING_STATE.SELECTING_FILTER,
        (prevValue, value) => prevValue + Delimiter.FILTER_VALUE + value
      );

    // runs for all states
    default:
      return stateMachine;
  }
};

const getCurrentFilterAndValue = (filteringState) => {
  const { context } = filteringState;
  const currentFilterValue = context.value.split(Delimiter.FILTER).at(-1);
  const currentFilter = currentFilterValue.split(Delimiter.FILTER_VALUE)?.[0] || "";
  const currentValue = currentFilterValue.split(Delimiter.FILTER_VALUE)?.[1] || "";
  return {
    filter : currentFilter,
    value : currentValue,
  };
};

const Filters = ({ filterStateMachine, dispatchFilterMachine }) => {
  const classes = useFilterStyles();
  const selectFilter = (filter) => {
    dispatchFilterMachine({
      type : FILTER_EVENTS.SELECT,
      payload : {
        value : filter,
      },
    });
  };

  const { filter : currentFilter } = getCurrentFilterAndValue(filterStateMachine);
  const matchingFilters = currentFilter
    ? Object.values(FILTERS).filter((filter) => filter.value.startsWith(currentFilter))
    : Object.values(FILTERS);
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
            <div key={filter.value} className={classes.item} disableGutters onClick={() => selectFilter(filter.value)}>
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

const FilterValueSuggestions = ({ filterStateMachine, dispatchFilterMachine }) => {
  const classes = useFilterStyles();

  const selectValue = (value) => {
    dispatchFilterMachine({
      type : FILTER_EVENTS.SELECT,
      payload : {
        value,
      },
    });
  };
  const { filter, value } = getCurrentFilterAndValue(filterStateMachine);
  const currentFilter = Object.values(FILTERS).find((f) => f.value == filter);
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
            <div key={value.value} className={classes.item} disableGutters onClick={() => selectValue(value)}>
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

const Filter = ({ handleFilter }) => {
  const theme = useTheme();
  // console.log("initialFilter", initialFilter)
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = useState(null);
  const isPopperOpen = Boolean(anchorEl);
  const inputFieldRef = useRef(null);
  const [filteringState, dispatch] = useReducer(filterReducer, {
    context : {
      value : "",
      prevValue : [""],
    },
    state : FILTERING_STATE.IDLE,
  });

  const handleFilterChange = (e) => {

    if (!anchorEl) {
      setAnchorEl(e.currentTarget)
    }

    if (e.target.value === "") {
      return dispatch({
        type : FILTER_EVENTS.CLEAR,
      });
    }

    return dispatch({
      type : FILTER_EVENTS.INPUT_CHANGE,
      payload : {
        value : e.target.value,
      },
    });
  };

  const handleClear = () => {
    dispatch({
      type : FILTER_EVENTS.EXIT,
    });

    handleFilter({})
  };

  const handleFocus = (e) => {
    setAnchorEl(e.currentTarget);
    dispatch({ type : "START" });
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
      if (e.key == "Enter") {
        handleFilter(getFilters(e.target.value))
        setAnchorEl(null);
      }
    }
    inputFieldRef?.current?.addEventListener("keydown", handleKeyDown)
    return () => {
      inputFieldRef?.current?.removeEventListener("keydown", handleKeyDown)
    }
  }, [inputFieldRef.current])



  return (
    <div className={clsx(classes.root, "mui-fixed")} >
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
          startAdornment : (
            <InputAdornment position="start">
              {" "}
              <ContentFilterIcon fill={theme.palette.secondary.iconMain} />{" "}
            </InputAdornment>
          ),
          endAdornment : (
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
        style={{ zIndex : 2000 }}
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
                    width : inputFieldRef.current ? inputFieldRef.current.clientWidth : 0,
                  }}
                >
                  {filteringState.state == FILTERING_STATE.SELECTING_FILTER && (
                    <Filters filterStateMachine={filteringState} dispatchFilterMachine={dispatch} />
                  )}
                  {filteringState.state == FILTERING_STATE.SELECTING_VALUE && (
                    <FilterValueSuggestions filterStateMachine={filteringState} dispatchFilterMachine={dispatch} />
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

export default Filter;
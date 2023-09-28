export const FILTERING_STATE = {
  IDLE: 'idle',
  SELECTING_FILTER: 'selecting_filter',
  SELECTING_VALUE: 'selecting_value',
};

export const FILTER_EVENTS = {
  START: 'start',
  SELECT: 'select_filter',
  SELECT_FILTER: 'select_filter',
  INPUT_CHANGE: 'input_change',
  SELECT_FILTER_VALUE: 'select_filter_value',
  CLEAR: 'clear',
  EXIT: 'exit',
};

export const Delimiter = {
  FILTER: ' ',
  FILTER_VALUE: ':',
};

const commonReducer = (stateMachine, action) => {
  const { context } = stateMachine;
  switch (action.type) {
    case FILTER_EVENTS.CLEAR:
      return {
        state: FILTERING_STATE.SELECTING_FILTER,
        context: {
          ...context,
          value: '',
          prevValue: [''],
        },
      };

    case FILTER_EVENTS.EXIT:
      return {
        state: FILTERING_STATE.IDLE,
        context: {
          ...context,
          value: '',
          prevValue: [''],
        },
      };

    default:
      return stateMachine;
  }
};

const filterSelectionReducer = (stateMachine, action, nextState, nextValue) => {
  const { state, context } = stateMachine;
  const nextDelimiter =
    nextState == FILTERING_STATE.SELECTING_FILTER ? Delimiter.FILTER : Delimiter.FILTER_VALUE;
  const prevDelimiter =
    nextDelimiter == Delimiter.FILTER_VALUE ? Delimiter.FILTER : Delimiter.FILTER_VALUE;
  const prevState = nextState; // same beccuase the prevState is the same as the nextState ( as we have only two states)
  switch (action.type) {
    // Select a filter and move to start entring its value
    case FILTER_EVENTS.SELECT: {
      const newValue = nextValue(context.prevValue.at(-1), action.payload.value); // ":" is used to separate the filter and its value)
      return {
        state: nextState,
        context: {
          ...context,
          value: newValue + nextDelimiter,
          prevValue: [...context.prevValue, newValue],
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
          state: prevState,
          context: {
            ...context,
            prevValue: context.prevValue.slice(0, -1),
            value: action.payload.value,
          },
        };
      }

      if (action.payload.value.endsWith(nextDelimiter)) {
        const newValue = action.payload.value;
        return {
          state: nextState,
          context: {
            ...context,
            value: action.payload.value,
            prevValue: [...context.prevValue, newValue.slice(0, -1)],
          },
        };
      }

      return {
        state, // stay in the same state
        context: {
          ...context,
          value: action.payload.value,
        },
      };
    default:
      return commonReducer(stateMachine, action);
  }
};

export const filterReducer = (stateMachine, action) => {
  const { state } = stateMachine;
  switch (state) {
    // Initial State
    case FILTERING_STATE.IDLE:
      switch (action.type) {
        // Start the filter process
        case 'START':
          return {
            ...stateMachine,
            state: FILTERING_STATE.SELECTING_FILTER,
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
        (prevValue, value) => prevValue + Delimiter.FILTER + value,
      );

    case FILTERING_STATE.SELECTING_VALUE:
      return filterSelectionReducer(
        stateMachine,
        action,
        FILTERING_STATE.SELECTING_FILTER,
        (prevValue, value) => prevValue + Delimiter.FILTER_VALUE + value,
      );

    // runs for all states
    default:
      return stateMachine;
  }
};

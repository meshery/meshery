import { assign, createMachine } from 'xstate';
import { reply } from '../utils';

export const DATA_VALIDATOR_COMMANDS = {
  VALIDATE_DATA: 'VALIDATE_DATA',
};

export const DATA_VALIDATOR_EVENTS = {
  // when the design validation is done successfully (doesnt mean the design is valid)
  DESIGN_VALIDATION_DONE: 'DESIGN_VALIDATION_DONE',
  // when the design is able to be validated ( either the design is corrupt or the there is some issue with the system)
  DESIGN_VALIDATION_FAILED: 'DESIGN_VALIDATION_FAILED',
};

export const dataValidatorCommands = {
  validateData: ({ validationPayload, returnAddress }) => ({
    type: DATA_VALIDATOR_COMMANDS.VALIDATE_DATA,
    returnAddress,
    data: { validationPayload },
  }),
};

export const dataValidatorEvents = {
  dataValidated: ({ validationPayload, validationResults }) => ({
    type: DATA_VALIDATOR_EVENTS.DESIGN_VALIDATION_DONE,
    data: { validationPayload, validationResults },
  }),

  dataValidationFailed: ({ validationPayload, systemErrors }) => ({
    type: DATA_VALIDATOR_EVENTS.DESIGN_VALIDATION_FAILED,
    data: { validationPayload, systemErrors },
  }),
};

export const dataValidatorMachine = createMachine(
  {
    id: 'validationMachine',

    initial: 'idle',

    context: {
      validationResults: null,
      validationPayload: {}, //data to be used in the validation process like the design to validate
      returnAddress: null, // the address to send the validation result to
    },

    states: {
      idle: {
        description:
          'when the machine is idle , i.e no process is going on and the machine the ready to accept new request',
        initial: 'waiting',
        on: {
          [DATA_VALIDATOR_COMMANDS.VALIDATE_DATA]: {
            target: '.debouncing',
            actions: ['setValidationPayload', 'setReturnAddress'],
          },
        },
        states: {
          waiting: {},
          debouncing: {
            after: {
              debounceTimeout: '#validationMachine.validatingData',
            },
          },
        },
      },

      validatingData: {
        invoke: {
          src: 'ValidateActor',
          input: ({ context }) => ({
            validationPayload: context.validationPayload,
            prevValidationResults: context.validationResults,
          }),
          onDone: {
            target: 'idle',
            actions: [
              reply(({ context, event }) =>
                dataValidatorEvents.dataValidated({
                  validationPayload: context.validationPayload,
                  validationResults: event.output.validationResults,
                }),
              ),
              'setValidationResults',
            ],
          },
          onError: {
            target: 'idle',
            actions: [
              reply(({ context, event }) =>
                dataValidatorEvents.dataValidationFailed({
                  validationPayload: context.validationPayload,
                  systemErrors: event.error,
                }),
              ),
              ({ event }) => console.error('Failed to validate data', event),
              assign({
                validationResults: ({ event }) => `Failed to validate data: ${event.error || ''}`,
              }),
            ],
          },
        },
      },
    },
  },

  {
    delays: {
      debounceTimeout: 300,
    },
    actors: {},
    actions: {
      setReturnAddress: assign({
        returnAddress: ({ event }) => event.returnAddress,
      }),
      setValidationPayload: assign({
        validationPayload: ({ event }) => event.data.validationPayload,
      }),

      resetValidationPayload: assign({
        validationPayload: null,
      }),

      resetValidationResults: assign({
        validationResults: null,
      }),
      setValidationResults: assign({
        validationResults: ({ event }) => event.output.validationResults,
      }),
    },
  },
);

export const selectValidationResults = (state) => state.context.validationResults;

export const selectIsValidating = (state) => state.matches('validatingData');

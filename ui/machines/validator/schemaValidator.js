import { createMachine, fromPromise, sendTo, assign, emit } from 'xstate';
import Ajv from 'ajv';
import _ from 'lodash';

import { processDesign } from '@/utils/utils';
import { dataValidatorMachine } from '@layer5/sistent';

const ajv = new Ajv({
  allErrors: true,
  strict: false, // allow additional properties like x-kubernetes-attributes ( this is safe the schema is sourced from the component definition and is not ours)
});

// dynamically add schemas to ajv to avoid recompiling the same schema and cache it
const validateSchema = (schema, data, id) => {
  let validate = ajv.getSchema(id);
  if (!validate) {
    ajv.addSchema(schema, id);
    validate = ajv.getSchema(id);
  }
  const valid = validate(data);

  return {
    isValid: valid,
    errors: validate.errors,
  };
};

// const AjvValidatorActor = createMachine({
//   id: 'AjvValidatorActor',
//   initial: 'idle',
//   context: {
//     ajv,
//   },
//   states: {
//     idle: {
//       entry: DeferEvents.recall,
//       on: {
//         VALIDATE: 'validating',
//       },
//     },
//     validating: {
//       entry: [
//         assign({
//           validationResults: ({ event, context }) => {
//             const schema = JSON.parse(event.data.schema);
//             const results = validateSchema(context.ajv, schema, event.data.data, event.data.id);
//             return results;
//           },
//         }),
//         sendTo(
//           ({ event }) => event.returnAddress,
//           ({ context }) => ({
//             type: 'VALIDATION_DONE',
//             data: context.validationResults,
//           }),
//         ),
//       ],

//       on: {
//         VALIDATE: DeferEvents.defer,
//       },
//     },
//   },
// });

const validateComponent = (component, validateAnnotations = false, componentDef) => {
  // const componentDef = await getComponentDefinition(component.type, component.model, {
  //   apiVersion: component.apiVersion,
  //   annotations: 'include',
  // });
  console.log('validateComponent', component, componentDef);

  if (!componentDef || (componentDef?.metadata?.isAnnotation && !validateAnnotations)) {
    // skip validation for annotations
    return {
      errors: [],
      componentDefinition: componentDef,
      component,
    };
  }
  const schema = JSON.parse(componentDef.component.schema);
  const results = validateSchema(schema, component.settings || {}, componentDef.id);

  const validationResults = {
    ...results,
    componentDefinition: componentDef,
    component,
  };

  return validationResults;
};

// const componentSchemaValidationMachine = createMachine({
//   id: 'componentSchemaValidationMachine',
//   initial: 'validating',
//   context: ({ input }) => ({
//     component: input.component,
//     validateAnnotations: input.validateAnnotations,
//   }),

//   states: {
//     getDefinition: {
//       entry: sendToActors([ACTOR_SYSTEM.RTK_QUERY], ({ context }) =>
//         rtkQueryActorCommands.initiateQuery({
//           endpointName: 'getComponentDefinition',
//           params: {
//             type: context.component.type,
//             model: context.component.model,
//             params: {
//               apiVersion: context.component.apiVersion,
//               annotations: 'include',
//             },
//           },
//         }),
//       ),

//       on: {
//         [RTK_EVENTS.QUERY_RESULT]: {
//           actions: ({ event }) => console.log('Query result wooo', event),
//           target: 'done',
//         },
//       },
//     },
//     done: {
//       type: 'final',
//       output: () => {
//         console.log('done');
//         return {
//           errors: [],
//         };
//       },
//     },
//   },
// });

export const componentKey = ({ type, model, apiVersion }) => `${type}-${model}-${apiVersion}`;

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

const validateDesign = (design, componentDefsStore) => {
  const { configurableComponents } = processDesign(design);

  const validationResults = {};

  for (const configurableComponent of configurableComponents) {
    try {
      const componentDef = componentDefsStore?.[componentKey(configurableComponent)];
      const componentValidationResults = validateComponent(
        configurableComponent,
        false,
        componentDef,
      );
      validationResults[configurableComponent.name] = componentValidationResults;
    } catch (error) {
      console.error('Error validating component', error);
    }
  }
  return validationResults;
};

const SchemaValidateDesignActor = fromPromise(async ({ input }) => {
  const { validationPayload, prevValidationResults } = input;
  const { validationPayloadType } = validationPayload;

  if (validationPayloadType === 'design') {
    console.log('Validating design', validationPayload);
    const { design, componentDefs } = validationPayload;
    const validationResults = validateDesign(design, componentDefs);
    console.log('Validation results', validationResults);
    return {
      validationResults,
    };
  }

  if (validationPayloadType === 'component') {
    console.log('Validating component', validationPayload);
    const { component, componentDef } = validationPayload;
    const validationResults = validateComponent(
      component,
      validationPayload.validateAnnotations || false,
      componentDef,
    );

    return {
      validationResults: _.set(prevValidationResults || {}, component.name, validationResults),
    };
  }

  throw new Error('Invalid validation payload type', validationPayloadType);
});

export const schemaValidatorMachine = dataValidatorMachine.provide({
  actors: {
    ValidateActor: SchemaValidateDesignActor,
  },
});

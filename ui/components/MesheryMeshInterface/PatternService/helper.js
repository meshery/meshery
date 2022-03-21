
import { formatString } from "../helpers";

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title : "" };
}

function deleteDescriptionFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, description : "" }
}

/**
 * remove top-level title, top-level description and
 * replace internal description with "help" key for
 * tooltip description
 *
 * @param {Object.<String, Object>} jsonSchema
 * @returns
 */
export function getRefinedJsonSchema(jsonSchema, hideTitle = true, handleError) {
  let refinedSchema;
  try {
    refinedSchema = hideTitle ? deleteTitleFromJSONSchema(jsonSchema) : jsonSchema
    refinedSchema = deleteDescriptionFromJSONSchema(refinedSchema)
    refinedSchema = addTitleToPropertiesJSONSchema(refinedSchema)
  } catch (e) {
    console.trace(e)
    handleError(e, "schema parsing problem")
  }
  return refinedSchema
}

/**
 * Check if the exsxceptional type fields are present in the schema
 * returns false if any of the exceptional-fields are not present else
 * returns the type of the exceptional-field
 *
 * @see https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#specifying-a-structural-schema
 */
function handleExceptionalFields(schema) {
  const xKubernetesIntOrString = 'x-kubernetes-int-or-string';
  const xKubernetesPreserveUnknownFields = 'x-kubernetes-preserve-unknown-fields';
  const additionalProperties = 'additionalProperties';

  const exceptionalFieldToTypeMap = {
    [xKubernetesIntOrString] : ['string'], // string can hold integers too
    [xKubernetesPreserveUnknownFields] : [schema?.type || 'object'],
    [additionalProperties] : ['string']
  }

  let returnedType;

  Object.keys(exceptionalFieldToTypeMap).some(field => {
    if ( Object.prototype.hasOwnProperty.call(schema, field) && schema[field] === true) {
      returnedType = exceptionalFieldToTypeMap[field];
      return;
    }
  })

  return returnedType || false;
}

function getDefaults(types) {
  if (types.includes('null')) {
    return null
  }
  if (types.includes('string')) {
    return ''
  }
  if (types.includes('integer')) {
    return 0
  } if (types.includes('array')) {
    return []
  }

  return ''
}

function addTitleToPropertiesJSONSchema(jsonSchema) {
  const newProperties = jsonSchema?.properties


  if (newProperties && typeof newProperties === 'object') {
    Object.keys(newProperties).map(key => {
      if (Object.prototype.hasOwnProperty.call(newProperties, key)) {
        let types = []

        const eField = handleExceptionalFields(newProperties[key]); // false if it is not a exceptional field
        if (eField) {
          types = eField;
          newProperties[key]['type'] = types
        } else {
          if (!Array.isArray(newProperties[key].type) && Object.prototype.hasOwnProperty.call(newProperties[key], 'type')) {
            types.push(newProperties[key].type)
          } else if (newProperties[key]?.type) {
            types.push(...newProperties[key].type)
          } else {
            console.log("an error occurred with schema", newProperties[key])
            throw new Error("There is something wrong with the json schema")
          }
        }

        const defaultValue = getDefaults(types)

        newProperties[key] = {
          ...newProperties[key],
          title : formatString(key),
          default : defaultValue,
        }
        // if (typeof newProperties[key] === 'object' && Object.prototype.hasOwnProperty.call(newProperties[key], 'properties')){
        //   newProperties[key] = {
        //     ...newProperties[key],
        //     properties : addTitleToPropertiesJSONSchema(newProperties[key])
        //   }
        // }
      }

    })

    return { ...jsonSchema, properties : newProperties };

  }
  return undefined
}

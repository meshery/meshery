

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

function addTitleToPropertiesJSONSchema(jsonSchema) {
  const newProperties = jsonSchema?.properties


  if (newProperties && typeof newProperties === 'object') {
    Object.keys(newProperties).map(key => {
      if (Object.prototype.hasOwnProperty.call(newProperties, key)) {
        const eField = handleExceptionalFields(newProperties[key]); // false if it is not a exceptional field
        if (eField) {
          newProperties[key]['type'] = eField
        }

        newProperties[key] = {
          ...newProperties[key],
        }
      }

    })

    return { ...jsonSchema, properties : newProperties };

  }
  return undefined
}

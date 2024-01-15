export const selectCompSchema = (enums, description, title, name, multiSelect = false) => {
  const intermediary =
    multiSelect === true ? getMultiselectSchema(enums) : getSingleSelectSchema(enums);

  return {
    properties: {
      [name]: {
        description: description,
        ...intermediary,
        title: title,
        uniqueItems: true,
        'x-rjsf-grid-area': 12,
      },
    },
    required: [name],
    type: 'object',
  };
};

const getSingleSelectSchema = (enums) => {
  return {
    enum: enums,
  };
};

const getMultiselectSchema = (enums) => {
  return {
    items: {
      enum: enums,
      type: 'string',
    },
    minItems: 1,
    type: 'array',
  };
};

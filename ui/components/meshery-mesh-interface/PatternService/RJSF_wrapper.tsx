import React from 'react';
import HandleError from '../../ErrorHandling';
import { buildUiSchema } from '../helpers';
import { getRefinedJsonSchema, hideRootObjectTitle } from './helper';
// import MesheryArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate";
// import MesheryCustomObjFieldTemplate from "./RJSFCustomComponents/ObjectFieldTemplate";
import _ from 'lodash';
import RJSFForm from './RJSF';

function RJSFWrapper(props) {
  const {
    formData,
    jsonSchema,
    onChange,
    // When set (default true), the root object's title is suppressed via the
    // UI schema in the effect below — not by mutating the JSON schema.
    hideTitle,
    uiSchema = {},
    formRef = null,
    liveValidate = true,
    RJSFWrapperComponent = React.Fragment,
    ...restProps
    //.. temporarily ignoring till handler is attached successfully
  } = props;
  const errorHandler = HandleError();

  const [data, setData] = React.useState((prev) => ({ ...formData, ...prev }));
  const [schema, setSchema] = React.useState({ rjsfSchema: {}, uiSchema: {} });
  const [isLoading, setIsLoading] = React.useState(true);

  React.useEffect(() => {
    // Apply debouncing mechanism for the state propagation
    const timer = setTimeout(() => {
      // callback fired, that triggers save operations, and other related side-effects
      onChange?.(data);
    }, 400);

    return () => clearTimeout(timer);
  }, [data]);

  React.useEffect(() => {
    const rjsfSchema = getRefinedJsonSchema(jsonSchema, errorHandler);
    // UI schema builds responsible for customizations in the RJSF fields shown to user
    const baseUiSchema = buildUiSchema(rjsfSchema);
    // Hide the duplicative top-level object title via the UI schema (the shared
    // Sistent standard: `ui:options.label = false`) rather than by mutating the
    // JSON schema. Defaults to hidden (`hideTitle` unset) to preserve behavior.
    const rootUiSchema = (hideTitle ?? true) ? hideRootObjectTitle(baseUiSchema) : baseUiSchema;
    setSchema({ rjsfSchema, uiSchema: rootUiSchema });
    // Recompute when the schema or its title visibility changes. (errorHandler
    // is only used on the parse-error path and is intentionally excluded.)
  }, [jsonSchema, hideTitle]);

  React.useEffect(() => {
    if (!_.isEqual(schema, { rjsfSchema: {}, uiSchema: {} })) {
      setTimeout(() => {
        setIsLoading(false);
      }, 300); // for showing circular progress
    }
  }, [schema]);

  return (
    <RJSFWrapperComponent {...props}>
      <RJSFForm
        isLoading={isLoading}
        schema={schema}
        formRef={formRef}
        uiSchema={uiSchema}
        data={data}
        liveValidate={liveValidate}
        onChange={(e) => {
          setData(e.formData);
        }}
        jsonSchema={jsonSchema}
        {...restProps}
      />
    </RJSFWrapperComponent>
  );
}

export default RJSFWrapper;

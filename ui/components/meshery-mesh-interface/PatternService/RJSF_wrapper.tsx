import React from 'react';
import HandleError from '../../ErrorHandling';
import { buildUiSchema } from '../helpers';
import { getRefinedJsonSchema } from './helper';
// import MesheryArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate";
// import MesheryCustomObjFieldTemplate from "./RJSFCustomComponents/ObjectFieldTemplate";
import _ from 'lodash';
import RJSFForm from './RJSF';

function RJSFWrapper(props) {
  const {
    formData,
    jsonSchema,
    onChange,
    hideTitle,
    uiSchema = {},
    formRef = null,
    liveValidate = true,
    RJSFWrapperComponent = React.Fragment,
    ...restProps
    //.. temporarily ignoring till handler is attached successfully
  } = props;
  const errorHandler = HandleError();

  const [data, setData] = React.useState({ ...formData });
  const [schema, setSchema] = React.useState({ rjsfSchema: {}, uiSchema: {} });
  const [isLoading, setIsLoading] = React.useState(true);

  // Track the last formData reference that we synced from, so we can tell
  // when the parent sends a genuinely new dataset (e.g. after async meshModels
  // resolves and rewrites the compatibility/Technology values) vs. a re-render
  // with the same prop value.
  const prevFormDataRef = React.useRef(formData);

  React.useEffect(() => {
    // Sync internal data when the parent provides a meaningfully different
    // formData. This covers the async race where meshModels resolves after
    // the first render and the parent calls setFormData() with the resolved
    // displayName values for the Technology/compatibility field.
    //
    // Guards:
    //  • Skip if formData is empty/undefined — no point resetting to nothing.
    //  • Skip if the reference hasn't changed — avoids redundant merges when
    //    the parent just re-renders without touching formData.
    if (!formData || Object.keys(formData).length === 0) return;
    if (formData === prevFormDataRef.current) return;
    prevFormDataRef.current = formData;
    setData((prev) => {
      // For each field arriving in the new formData, use the incoming value
      // only if the current value is empty/absent (i.e. was never filled in due
      // to the async race), so we don't clobber fields the user has already edited.
      const merged = { ...prev };
      let changed = false;
      Object.keys(formData).forEach((key) => {
        const current = prev[key];
        const isEmpty =
          current === undefined ||
          current === null ||
          current === '' ||
          (Array.isArray(current) && current.length === 0);
        if (isEmpty) {
          merged[key] = formData[key];
          changed = true;
        }
      });
      return changed ? merged : prev;
    });
  }, [formData]);

  React.useEffect(() => {
    // Apply debouncing mechanism for the state propagation
    const timer = setTimeout(() => {
      // callback fired, that triggers save operations, and other related side-effects
      onChange?.(data);
    }, 400);

    return () => clearTimeout(timer);
  }, [data]);

  React.useEffect(() => {
    const rjsfSchema = getRefinedJsonSchema(jsonSchema, hideTitle, errorHandler);
    // UI schema builds responsible for customizations in the RJSF fields shown to user
    const uiSchema = buildUiSchema(rjsfSchema);
    setSchema({ rjsfSchema, uiSchema });
  }, [jsonSchema]); // to reduce heavy lifting on every react render

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

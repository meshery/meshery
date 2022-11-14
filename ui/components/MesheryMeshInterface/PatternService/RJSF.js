import { MuiThemeProvider } from '@material-ui/core/styles';
import Form, { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme, } from "@rjsf/material-ui";
import React, { useEffect } from "react";
import JS4 from "../../../assets/jsonschema/schema-04.json";
import { rjsfTheme } from "../../../themes";
import { recursiveCleanObject } from "../helpers";
import MesheryArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate";
import MesheryCustomObjFieldTemplate from "./RJSFCustomComponents/ObjectFieldTemplate";
import MesheryWrapIfAdditionalTemplate from './RJSFCustomComponents/WrapIfAdditionalTemplate';
import { customizeValidator } from "@rjsf/validator-ajv6";
import _ from "lodash"
import CustomTextWidget from './RJSFCustomComponents/CustomTextWidget';
import CustomDateTimeWidget from './RJSFCustomComponents/CustomDateTimeWidget';

class RJSFOverridenComponent extends Form {
  constructor(props){
    super(props)
    let oldValidate = this.validate;
    this.validate = (
      formData,
      schema,
    ) => {
      let fixedFormData = recursiveCleanObject(_.cloneDeep(formData));
      return oldValidate.call(this, fixedFormData, schema);
    }
  }
}

// This is Patched change to include customised Forms
const MuiRJSFForm = withTheme(MaterialUITheme, RJSFOverridenComponent);
const validator = customizeValidator({ additionalMetaSchemas : [JS4] });

/**
 * The Custom RJSF Form that accepts custom fields from the extension
 * or seed it's own default
 * Adding a new custom component:
 * 1. Pass the new prop from the Meshery Extension
 * 2. Extract from props in the RJSFForm Component
 * @param {*} props
 * @returns
 */
function RJSFForm(props) {
  const {
    schema,
    jsonSchema,
    data,
    onChange,
    isLoading,
    ArrayFieldTemplate = MesheryArrayFieldTemplate,
    ObjectFieldTemplate = MesheryCustomObjFieldTemplate,
    WrapIfAdditionalTemplate = MesheryWrapIfAdditionalTemplate,
    LoadingComponent,
    ErrorList,
    // prop should be present in order for the cloned element to override this property
    transformErrors
  } = props;
  const templates={
    ArrayFieldTemplate,
    ObjectFieldTemplate,
    WrapIfAdditionalTemplate,
  }

  useEffect(() => {
    const extensionTooltipPortal = document.getElementById("extension-tooltip-portal");
    if (extensionTooltipPortal) {
      rjsfTheme.props.MuiMenu.container = extensionTooltipPortal;
    }
    rjsfTheme.zIndex.modal = 99999;
  }, [])

  if (isLoading && LoadingComponent) {
    return <LoadingComponent />
  }

  return (
    <MuiThemeProvider theme={rjsfTheme}>
      <MuiRJSFForm
        schema={schema.rjsfSchema}
        idPrefix={jsonSchema?.title}
        onChange={onChange}
        formData={data}
        validator={validator}
        templates={templates}
        uiSchema={schema.uiSchema}
        widgets={{
          TextWidget : CustomTextWidget,
          // Custom components to be added here
          DateTimeWidget : CustomDateTimeWidget,
          // SelectWidget: CustomSelectWidget,
          // CheckboxWidget: CustomBooleanWidget,
        }}
        liveValidate
        showErrorList={false}
        noHtml5Validate
        ErrorList={ErrorList}
        transformErrors={transformErrors}
      >
        {/* {hideSubmit ? true : <RJSFButton handler={onSubmit} text="Submit" {...restparams} />}
{hideSubmit ? true : <RJSFButton handler={onDelete} text="Delete" />} */}
        {/* <RJSFFormChildComponent /> */}
        <></> {/* temporary change for functionality */}
      </MuiRJSFForm>

    </MuiThemeProvider>
  )
}

export default RJSFForm;

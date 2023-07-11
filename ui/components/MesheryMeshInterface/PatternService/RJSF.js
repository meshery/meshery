import { MuiThemeProvider, useTheme } from '@material-ui/core/styles';
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import ajv8validator from "@rjsf/validator-ajv8";
import React, { useEffect } from "react";
import { rjsfTheme } from "../../../themes";
import darkRjsfTheme from '../../../themes/rjsf';
import { CustomTextTooltip } from './CustomTextTooltip';
import MesheryArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate";
import CustomDateTimeWidget from './RJSFCustomComponents/CustomDateTimeWidget';
import CustomTextWidget from './RJSFCustomComponents/CustomTextWidget';
import { CustomFieldTemplate } from './RJSFCustomComponents/FieldTemplate';
import MesheryCustomObjFieldTemplate from "./RJSFCustomComponents/ObjectFieldTemplate";
import MesheryWrapIfAdditionalTemplate from './RJSFCustomComponents/WrapIfAdditionalTemplate';
import _ from "lodash"
import { CustomCheckboxWidget } from './RJSFCustomComponents/CustomCheckboxWidget';
import CustomSelectWidget from './RJSFCustomComponents/CustomSelectWidget';
import CustomTextAreaWidget from './RJSFCustomComponents/CustomTextAreaWidget';

const MuiRJSFForm = withTheme(MaterialUITheme);

/**
 * The Custom RJSF Form that accepts custom fields from the extension
 * or seed it's own default
 * Adding a new custom component:
 * 1. Pass the new prop from the Meshery Extension
 * 2. Extract from props in the RJSFForm Component
 * @param {*} props
 * @returns
 */
function RJSFForm({
  schema,
  jsonSchema,
  data,
  onChange,
  isLoading,
  ArrayFieldTemplate = MesheryArrayFieldTemplate,
  ObjectFieldTemplate = MesheryCustomObjFieldTemplate,
  WrapIfAdditionalTemplate = MesheryWrapIfAdditionalTemplate,
  LoadingComponent,
  // prop should be present in order for the cloned element to override this property
  transformErrors,
  override,
  uiSchema={}
}) {
  const globalTheme = useTheme();
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
    <MuiThemeProvider
      theme={globalTheme.palette.type === "dark" ? darkRjsfTheme : rjsfTheme}>
      <MuiRJSFForm
        schema={schema.rjsfSchema}
        idPrefix={jsonSchema?.title}
        onChange={onChange}
        formData={data}
        validator={ajv8validator}
        templates={{
          ArrayFieldTemplate,
          ObjectFieldTemplate,
          WrapIfAdditionalTemplate,
          FieldTemplate : CustomFieldTemplate, // applying field template universally to every field type.
        }}
        formContext={{ overrideFlag : override, CustomTextTooltip : CustomTextTooltip }}
        uiSchema={_.merge(schema.uiSchema, uiSchema)}
        widgets={{
          // Custom components to be added here
          TextWidget : CustomTextWidget,
          DateTimeWidget : CustomDateTimeWidget,
          SelectWidget : CustomSelectWidget,
          CheckboxWidget : CustomCheckboxWidget,
          TextareaWidget : CustomTextAreaWidget,
        }}
        liveValidate
        showErrorList={false}
        noHtml5Validate
        transformErrors={transformErrors}
      >
        <div></div>
      </MuiRJSFForm>
    </MuiThemeProvider>
  )
}

export default RJSFForm;
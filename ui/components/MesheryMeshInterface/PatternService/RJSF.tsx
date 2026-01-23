import { withTheme } from '@rjsf/core';
import { Theme as MaterialUITheme } from '@rjsf/mui';
import customValidator from '../../../utils/rjsfValidator';
import React, { useEffect, useMemo, useState } from 'react';
import { rjsfTheme } from '../../../themes';
import darkRjsfTheme from '../../../themes/rjsf';
import { createTheme, useTheme, ThemeProvider } from '@sistent/sistent';
import { CustomTextTooltip } from './CustomTextTooltip';
import MesheryArrayFieldTemplate from './RJSFCustomComponents/ArrayFieldTemlate';
import CustomDateTimeWidget from './RJSFCustomComponents/CustomDateTimeWidget';
import CustomTextWidget from './RJSFCustomComponents/CustomTextWidget';
import { CustomFieldTemplate } from './RJSFCustomComponents/FieldTemplate';
import MesheryCustomObjFieldTemplate from './RJSFCustomComponents/ObjectFieldTemplate';
import MesheryWrapIfAdditionalTemplate from './RJSFCustomComponents/WrapIfAdditionalTemplate';
import _ from 'lodash';
import { CustomCheckboxWidget } from './RJSFCustomComponents/CustomCheckboxWidget';
import MesheryCustomSelectWidget from './RJSFCustomComponents/CustomSelectWidget';
import CustomTextAreaWidget from './RJSFCustomComponents/CustomTextAreaWidget';
import CustomFileWidget from './RJSFCustomComponents/CustomFileWidget';
import CustomURLWidget from './RJSFCustomComponents/CustomURLWidget';
import CustomColorWidget from './RJSFCustomComponents/CustomColorWidget';
// @ts-expect-error
import { ErrorBoundary } from '@sistent/sistent';
import CustomErrorFallback from '@/components/General/ErrorBoundary';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';

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
function RJSFForm_({
  schema,
  jsonSchema,
  data,
  onChange,
  isLoading,
  children,
  ArrayFieldTemplate = MesheryArrayFieldTemplate,
  ObjectFieldTemplate = MesheryCustomObjFieldTemplate,
  BaseInputTemplate,
  WrapIfAdditionalTemplate = MesheryWrapIfAdditionalTemplate,
  SelectWidget = MesheryCustomSelectWidget,
  LoadingComponent,
  liveValidate,
  // prop should be present in order for the cloned element to override this property
  transformErrors,
  override,
  formRef = null,
  uiSchema = {},
  validator,
  fieldTemplates = {},
  customFields = {},
  widgets = {},
  extraErrors,
  isExtensionTooltipPortal = true,
  ...restProps
}) {
  const globalTheme = useTheme();
  const [menuContainer, setMenuContainer] = useState<HTMLElement | null>(null);
  const baseTheme = globalTheme.palette.mode === 'dark' ? darkRjsfTheme : rjsfTheme;

  useEffect(() => {
    if (!isExtensionTooltipPortal) {
      setMenuContainer(null);
      return;
    }

    const extensionTooltipPortal = document.getElementById('extension-tooltip-portal');
    setMenuContainer(extensionTooltipPortal || null);
  }, [isExtensionTooltipPortal]);

  const resolvedRjsfTheme = useMemo(() => {
    if (!menuContainer) {
      return baseTheme;
    }

    return createTheme(baseTheme, {
      components: {
        MuiMenu: {
          defaultProps: {
            ...baseTheme.components?.MuiMenu?.defaultProps,
            container: menuContainer,
          },
        },
      },
    });
  }, [baseTheme, menuContainer]);

  if (isLoading && LoadingComponent) {
    return <LoadingComponent />;
  }

  return (
    <ErrorBoundary customFallback={CustomErrorFallback}>
      {/* Putting RJSF into error boundary, so that error can be catched.. */}{' '}
      <ThemeProvider theme={resolvedRjsfTheme}>
        <MuiRJSFForm
          schema={schema.rjsfSchema}
          idPrefix={jsonSchema?.title}
          ref={formRef}
          onChange={onChange}
          formData={data}
          extraErrors={extraErrors}
          validator={validator || customValidator}
          templates={{
            ArrayFieldTemplate,
            ObjectFieldTemplate: ObjectFieldTemplate as any,
            WrapIfAdditionalTemplate: WrapIfAdditionalTemplate as any,
            BaseInputTemplate,
            FieldTemplate: CustomFieldTemplate, // applying field template universally to every field type.
            ...fieldTemplates,
          }}
          formContext={{ overrideFlag: override, CustomTextTooltip: CustomTextTooltip }}
          uiSchema={_.merge(schema.uiSchema, uiSchema)}
          widgets={{
            // Custom components to be added here
            TextWidget: CustomTextWidget,
            DateTimeWidget: CustomDateTimeWidget,
            SelectWidget: SelectWidget as any,
            ColorWidget: CustomColorWidget,
            CheckboxWidget: CustomCheckboxWidget,
            TextareaWidget: CustomTextAreaWidget,
            FileWidget: CustomFileWidget,
            URLWidget: CustomURLWidget,
            ...widgets,
          }}
          fields={customFields}
          liveValidate={liveValidate}
          showErrorList={false}
          noHtml5Validate
          transformErrors={transformErrors}
          {...restProps}
        >
          {children}
          <div></div>
        </MuiRJSFForm>
      </ThemeProvider>
    </ErrorBoundary>
  );
}

const RJSFForm = (props) => {
  return <ProviderStoreWrapper>{<RJSFForm_ {...props} />}</ProviderStoreWrapper>;
};
export default RJSFForm;

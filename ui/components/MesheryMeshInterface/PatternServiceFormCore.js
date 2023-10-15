// import React from 'react';
// @ts-check
import useStateCB from '../../utils/hooks/useStateCB';
import PatternService from './PatternService';
import { getPatternAttributeName, createPatternFromConfig } from './helpers';
import React, { useEffect, useState } from 'react';
import { scrollToTop } from '../../utils/utils';

/**
 * usePatternServiceForm seperates the form logic from its UI representation
 * @param {{
 *  schemaSet: { workload: any, type: string };
 *  onSubmit: Function;
 *  onDelete: Function;
 *  namespace: string;
 *  onChange?: Function
 *  onSettingsChange?: Function;
 *  formData?: Record<String, unknown>
 *  reference?: Record<any, any>;
 * 	children?: Function;
 *  scroll?: Boolean; // If the window should be scrolled to zero after re-rendering
 * }} param0 props for the component
 */
function PatternServiceFormCore({
  formData,
  schemaSet,
  onSubmit,
  onDelete,
  reference,
  namespace,
  onSettingsChange,
  children,
  scroll = false,
}) {
  const [settings, setSettings, getSettingsRefValue] = useStateCB(
    formData && !!formData.settings ? formData.settings : {},
    onSettingsChange,
  );
  const [update, forceUpdate] = useState(0);
  // eslint-disable-next-line no-unused-vars
  const [rjsfReferenceKey, _] = useState((Math.random() + 1).toString(32).substring(2));

  useEffect(() => {
    if (schemaSet.type !== 'addon') {
      child.current = children(...propagatedChildren());
      forceUpdate(update + 1); // updating the state for simulating re-rendering of changed children
      scroll && scrollToTop();
    }
  }, [schemaSet]);

  const child = React.useRef(null);

  const submitHandler = (val) => {
    try {
      onSubmit?.(
        createPatternFromConfig({ [getPatternAttributeName(schemaSet.workload)]: val }, namespace),
      );
    } catch (e) {
      console.log('error while submitting form-data', e);
    }
  };

  const deleteHandler = (val) => {
    onDelete?.(
      createPatternFromConfig({ [getPatternAttributeName(schemaSet.workload)]: val }, namespace),
      true,
    );
  };

  const propagatedChildren = () => [
    function (props = {}) {
      return (
        <PatternService
          type="workload"
          formData={settings}
          jsonSchema={schemaSet.workload}
          onChange={setSettings}
          onSubmit={() => submitHandler({ settings: getSettingsRefValue() })}
          onDelete={() => deleteHandler({ settings: getSettingsRefValue() })}
          {...props}
        />
      );
    },
  ];

  if (reference) {
    if (!reference.current) reference.current = {};

    reference.current.submit = (cb) => submitHandler(cb?.(getSettingsRefValue()));
    reference.current.delete = (cb) => deleteHandler(cb?.(getSettingsRefValue()));
    reference.current.getSettings = () => getSettingsRefValue();
    reference.current.referKey = rjsfReferenceKey;
  }

  // Return cached child -- Prevents rerenders
  if (child.current) return child.current;

  if (schemaSet.type === 'addon') {
    return (child.current = children(
      function (props = {}) {
        return (
          <PatternService
            formData={settings}
            type="workload"
            jsonSchema={schemaSet.workload}
            onChange={setSettings}
            onSubmit={() => submitHandler({ settings: getSettingsRefValue() })}
            onDelete={() => deleteHandler({ settings: getSettingsRefValue() })}
            {...props}
          />
        );
      },
      () => null,
    ));
  }

  return child.current;
}

export default PatternServiceFormCore;

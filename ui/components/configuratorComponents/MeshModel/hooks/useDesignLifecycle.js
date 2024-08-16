import { useEffect, useState } from 'react';
import jsYaml from 'js-yaml';
// eslint-disable-next-line no-unused-vars
import * as Types from './types';
import { promisifiedDataFetch } from '../../../../lib/data-fetch';
import { useNotification } from '../../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../../lib/event-types';
import { getUnit8ArrayForDesign } from '@/utils/utils';

export default function useDesignLifecycle() {
  const [designId, setDesignId] = useState();
  const [designJson, setDesignJson] = useState({
    name: 'Untitled Design',
    components: [],
    schemaVersion: 'designs.meshery.io/v1beta1',
  });
  const [designYaml, setDesignyaml] = useState('');
  const { notify } = useNotification();

  useEffect(
    function updateDesignYamlFromJson() {
      setDesignyaml(jsYaml.dump(designJson));
    },
    [designJson],
  );

  /**
   *
   * @param {Types.ComponentDefinition} componentDefinition
   */
  function onSettingsChange(componentDefinition, formReference) {
    const {
      component,
      schemaVersion,
      version,
      model: {
        name: modelName,
        registrant: modelRegistrant,
        version: modelVersion,
        category: modelCategory,
      },
    } = componentDefinition;

    /**
     * Handles actual design-json change in response to form-data change
     * @param {*} formData
     */
    return function handledesignJsonChange(formData) {
      const referKey = formReference.current.referKey;
      const { name, namespace, labels, annotations, ...configuration } = formData;
      const newInput = {
        id: referKey,
        schemaVersion,
        version,
        component,
        displayName: name,
        model: {
          name: modelName,
          version: modelVersion,
          category: modelCategory,
          registrant: modelRegistrant,
        },
        configuration: {
          metadata: {
            labels,
            annotations,
            namespace,
          },
          ...configuration,
        },
      };
      setDesignJson((prev) => {
        let newestKey = false;
        const currentJson =
          prev.components?.map((val) => {
            if (val.id == referKey) {
              newestKey = true;
              return newInput;
            }
            return val;
          }) || [];
        if (!newestKey) {
          currentJson.push(newInput);
        }

        return { ...prev, components: [...currentJson] };
      });
    };
  }

  function onSubmit() {}

  function onDelete() {}

  function designSave() {
    promisifiedDataFetch('/api/pattern', {
      body: JSON.stringify({
        pattern_data: {
          name: designJson.name,
          pattern_file: getUnit8ArrayForDesign(designYaml),
        },
        save: true,
      }),
      method: 'POST',
    })
      .then((data) => {
        setDesignId(data[0].id);
        notify({
          message: `"${designJson.name}" saved successfully`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch((err) => {
        notify({
          message: `failed to save design file`,
          event_type: EVENT_TYPES.ERROR,
          details: err.toString(),
        });
      });
  }

  async function designUpdate() {
    try {
      await promisifiedDataFetch('/api/pattern', {
        body: JSON.stringify({
          pattern_data: {
            name: designJson.name,
            pattern_file: getUnit8ArrayForDesign(designYaml),
            id: designId,
          },
        }),
        method: 'POST',
      });
      notify({
        message: `"${designJson.name}" updated successfully`,
        event_type: EVENT_TYPES.SUCCESS,
      });
    } catch (err) {
      notify({
        message: `failed to update design file`,
        event_type: EVENT_TYPES.ERROR,
        details: err.toString(),
      });
    }
  }

  async function designDelete() {
    try {
      await promisifiedDataFetch('/api/pattern/' + designId, { method: 'DELETE' });
      notify({ message: `Design "${designJson.name}" Deleted`, event_type: EVENT_TYPES.SUCCESS });
      setDesignId(undefined);
      setDesignJson({
        name: 'Untitled Design',
        components: [],
        schemaVersion: 'designs.meshery.io/v1beta1',
      });
    } catch (err) {
      return notify({
        message: `failed to delete design file`,
        event_type: EVENT_TYPES.ERROR,
        details: err.toString(),
      });
    }
  }

  const updateDesignName = (name) => {
    setDesignJson((prev) => ({ ...prev, name: name }));
  };

  const loadDesign = async (design_id) => {
    try {
      const data = await promisifiedDataFetch('/api/pattern/' + design_id);
      setDesignId(design_id);
      setDesignJson(jsYaml.load(data.pattern_file));
    } catch (err) {
      notify({
        message: `failed to load design file`,
        event_type: EVENT_TYPES.ERROR,
        details: err.toString(),
      });
    }
  };

  const updateDesignData = ({ yamlData }) => {
    try {
      const designData = jsYaml.load(yamlData);
      setDesignJson(designData);
    } catch (err) {
      notify({
        message: `Invalid Yaml Data`,
        event_type: EVENT_TYPES.ERROR,
        details: err.toString(),
      });
    }
  };
  return {
    designJson,
    onSettingsChange,
    onSubmit,
    onDelete,
    designYaml,
    designSave,
    designUpdate,
    designId,
    designDelete,
    updateDesignName,
    loadDesign,
    updateDesignData,
  };
}

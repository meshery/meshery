import { useEffect, useState } from "react";
import jsYaml from "js-yaml";
// eslint-disable-next-line no-unused-vars
import * as Types from "./types";
import { promisifiedDataFetch } from "../../../../lib/data-fetch";
import { useNotification } from "../../../../utils/hooks/useNotification";
import { EVENT_TYPES } from "../../../../lib/event-types";

export default function useDesignLifecycle() {
  const [designName, setDesignName] = useState("Unitled Design")
  const [designId, setDesignId] = useState();
  const [designJson, setDesignJson] = useState({
    name : designName,
    services : {}
  })
  const [designYaml, setDesignyaml] = useState("");
  const { notify } = useNotification()


  useEffect(function updateDesignYamlFromJson() {
    setDesignyaml(jsYaml.dump(designJson))
  }, [designJson])

  /**
   *
   * @param {Types.ComponentDefinition} componentDefinition
   */
  function onSettingsChange(componentDefinition, formReference) {
    const { kind, apiVersion, model } = componentDefinition;
    const modelVersion = model.version;
    const modelName = model.name;

    /**
     * Handles actual design-json change in response to form-data change
     * @param {*} formData
     */
    return function handledesignJsonChange(formData) {
      const referKey = formReference.current.referKey;
      const { name, namespace, labels, annotations, ...settings } = formData;

      const currentJson = { ...designJson };
      currentJson.services[referKey] = {
        name,
        namespace,
        labels,
        annotations,
        type : kind,
        apiVersion,
        model : modelName,
        version : modelVersion,
        settings
      }
      setDesignJson(currentJson);
    }
  }

  function onSubmit() {
  }

  function onDelete() {

  }

  function designSave() {
    promisifiedDataFetch("/api/pattern", {
      body : JSON.stringify({
        pattern_data : {
          name : designName,
          pattern_file : designYaml
        },
        save : true
      }),
      method : "POST"
    }).then(data => {
      setDesignId(data[0].id);
      notify({ message : `"${designName}" saved successfully`, event_type : EVENT_TYPES.SUCCESS })
    }).catch((err) => {
      notify({ message : `failed to save design file`, event_type : EVENT_TYPES.ERROR, details : err.toString() })
    })
  }

  async function designUpdate() {
    try {
      await promisifiedDataFetch("/api/pattern", {
        body : JSON.stringify({
          pattern_data : {
            name : designName,
            pattern_file : designYaml,
            id : designId
          }
        }),
        method : "POST"
      });
      notify({ message : `"${designName}" updated successfully`, event_type : EVENT_TYPES.SUCCESS });
    } catch (err) {
      notify({ message : `failed to update design file`, event_type : EVENT_TYPES.ERROR, details : err.toString() });
    }
  }

  async function designDelete() {
    try {
      await promisifiedDataFetch("/api/pattern/" + designId, { method : "DELETE" });
      notify({ message : `Design "${designName}" Deleted`, event_type : EVENT_TYPES.SUCCESS });
      setDesignId(undefined);
      setDesignName("Unitled Design");
    } catch (err) {
      return notify({ message : `failed to delete design file`, event_type : EVENT_TYPES.ERROR, details : err.toString() });
    }
  }

  const updateDesignName = (name) => {
    setDesignName(name);
  }

  const loadDesign = async (design_id) => {
    try {
      const data = await promisifiedDataFetch("/api/pattern/" + design_id);
      console.log("loaded data file", data);
      setDesignId(design_id);
      setDesignName(data.name);
      setDesignJson(jsYaml.load(data.pattern_file));
    } catch (err) {
      notify({ message : `failed to load design file`, event_type : EVENT_TYPES.ERROR, details : err.toString() });
    }
  }

  const updateDesignData =({ yamlData }) => {
    try {
      const designData = jsYaml.load(yamlData);
      setDesignJson(designData);
    } catch (err) {
      notify({ message : `Invalid Yaml Data`, event_type : EVENT_TYPES.ERROR, details : err.toString() });
    }
  }



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
    designName,
    updateDesignName,
    loadDesign,
    updateDesignData
  }
}
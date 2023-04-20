import { useEffect, useState } from "react";
import jsYaml from "js-yaml";
// eslint-disable-next-line no-unused-vars
import * as Types from "./types";

export default function useDesignLifecycle() {
  const [designJson, setDesignJson] = useState({
    name : "Unitled Design",
    services : {}
  })
  const [designYaml, setDesignyaml] = useState("");


  useEffect(function updateDesignYamlFromJson() {
    console.log("....des", designJson)
    setDesignyaml(jsYaml.dump(designJson))
  }, [designJson])

  /**
   *
   * @param {Types.ComponentDefinition} componentDefinition
   */
  function onSettingsChange(componentDefinition, formReference) {
    const { kind, apiVersion, model, metadata } = componentDefinition;
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
        settings,
        traits : {
          "meshmodel-metadata" : { ...metadata },
          ...formReference.current.getTraits()
        }
      }
      setDesignJson(currentJson);
    }
  }

  function onSubmit(e, v, f) {
    console.log({ e, v, f })

  }

  function onDelete(formData) {
    console.log(formData)
  }

  return {
    designJson,
    onSettingsChange,
    onSubmit,
    onDelete,
    designYaml
  }
}
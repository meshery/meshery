import { useEffect, useState } from "react";
import jsYaml from "js-yaml";
// eslint-disable-next-line no-unused-vars
import * as Types from "./types";
import { promisifiedDataFetch } from "../../../../lib/data-fetch";
import { useSnackbar } from "notistack";

export default function useDesignLifecycle() {
  const [designName, setDesignName] = useState("Unitled Design")
  const [designId, setDesignId] = useState();
  const [designJson, setDesignJson] = useState({
    name : designName,
    services : {}
  })
  const [designYaml, setDesignyaml] = useState("");
  const { enqueueSnackbar } = useSnackbar();


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
      enqueueSnackbar(`"${designName}" Design saved`, { variant : "success" })
    }).catch(() => {
      enqueueSnackbar("failed to save design file", { variant : "error" })
    })
  }

  function designUpdate() {
    return promisifiedDataFetch("/api/pattern", {
      body : JSON.stringify({
        pattern_data : {
          name : designName,
          pattern_file : designYaml,
          id : designId
        }
      }),
      method : "POST"
    }).then(() => {
      enqueueSnackbar(`"${designName}" Design updated`, { variant : "success" })
    }).catch(() => {
      enqueueSnackbar(`couldn't update "${designName}"`, { variant : "error" })
    })
  }

  function designDelete() {
    return promisifiedDataFetch("/api/pattern/" + designId, { method : "DELETE" }).then(() => {
      enqueueSnackbar(`"${designName}" Design deleted`, { variant : "success" })
      setDesignJson({
        name : "Unitled Design",
        services : {}
      })
      setDesignId(undefined)
      setDesignName("Unitled Design")
    }).catch(() => enqueueSnackbar(`error deleting "${designName}" design`, { variant : "error" })
    )
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
    designDelete
  }
}
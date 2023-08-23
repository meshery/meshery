import { Button, Typography,FormGroup,TextField,InputAdornment } from '@material-ui/core'
import React from 'react'
import { useRef } from 'react';
import AddIconCircleBorder from '../assets/icons/AddIconCircleBorder'
import PromptComponent from './PromptComponent';
import CloudUploadIcon from "@material-ui/icons/CloudUpload";
import dataFetch, { promisifiedDataFetch } from "../lib/data-fetch";
import { updateProgress } from '../lib/store';
import { extractKubernetesCredentials } from './ConnectionWizard/helpers/kubernetesHelpers';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';

const MesherySettingsEnvButtons = () => {
  let k8sfileElementVal = "";
  let formData = new FormData();
  const ref = useRef(null)
  const { notify } = useNotification()

  const handleConfigSnackbars = ctxs => {
    updateProgress({ showProgress : false });
    for (let ctx of ctxs.inserted_contexts) {
      handleCredentialsPost(ctx);
      const msg = `Cluster ${ctx.name} at ${ctx.server} connected`
      notify({ message : msg, event_type : EVENT_TYPES.SUCCESS })
    }
    for (let ctx of ctxs.updated_contexts) {
      const msg = `Cluster ${ctx.name} at ${ctx.server} already exists`
      notify({ message : msg, event_type : EVENT_TYPES.INFO })
    }

    for (let ctx of ctxs.errored_contexts) {
      const msg = `Failed to add cluster ${ctx.name} at ${ctx.server}`
      notify({ message : msg, event_type : EVENT_TYPES.ERROR , details : ctx.error.toString() })
    }
  }

  function handleCredentialsPost(obj){
    // right now we just posting the credentials when we insert a new context
    const data = {
      name : obj.name,
      type : "kubernetes",
      secret : extractKubernetesCredentials(obj),
    }

    dataFetch(
      "/api/integrations/credentials",
      {
        credentials : "include",
        method : "POST",
        body : JSON.stringify(data),
      },
      () => {
        notify({ message : "Credentials saved successfully!", event_type : EVENT_TYPES.SUCCESS })
      }
    );
  }

  const handleError = (msg) => (error) => {
    updateProgress({ showProgress : false });
    notify({ message : `${msg}: ${error}`, event_type : EVENT_TYPES.ERROR, details : error.toString() })
  };
  const handleChange = () => {
    const field = document.getElementById("k8sfile");
    const textField = document.getElementById("k8sfileLabelText");
    if (field instanceof HTMLInputElement) {
      if (field.files.length < 1) return;
      const name = field.files[0].name;
      const formdata = new FormData();
      formdata.append("k8sfile", field.files[0])
      textField.value = name;
      formData = formdata;

    }
  }
  const uploadK8SConfig = async () => {
    return await promisifiedDataFetch(
      "/api/system/kubernetes",
      {
        method : "POST",
        body : formData,
      }
    )
  }
  const handleClick = async () => {
    const modal = ref.current;
    let response = await modal.show({
      title : "Add Kubernetes Cluster(s)",
      subtitle :
        <>
          <div
            style={{ "overflow" : "hidden" }}
          >
            <Typography variant="h6">
              Upload your kubeconfig
            </Typography>
            <Typography variant="body2">
              commonly found at ~/.kube/config
            </Typography>
            <FormGroup>
              <input
                id="k8sfile"
                type="file"
                value={k8sfileElementVal}
                onChange={handleChange}
                style={{ "display" : "none", }}
              />

              <TextField
                id="k8sfileLabelText"
                name="k8sfileLabelText"
                style={{ "cursor" : "pointer", }}
                placeholder="Upload kubeconfig"
                variant="outlined"
                fullWidth
                onClick={() => {
                  document.querySelector("#k8sfile")?.click();
                }}
                margin="normal"
                InputProps={{
                  readOnly : true,
                  endAdornment : (
                    <InputAdornment position="end">
                      <CloudUploadIcon />
                    </InputAdornment>
                  ),
                }}
              />
            </FormGroup>
          </div>
        </>,
      options : ["IMPORT","CANCEL"]
    })

    if (response === "IMPORT") {
      if (formData.get("k8sfile") === null) {
        handleError("No file selected")("Please select a valid kube config")
        return;
      }

      const inputFile = ( formData.get( "k8sfile" ).name );
      const invalidExtensions = /^.*\.(jpg|gif|jpeg|pdf|png|svg)$/i;

      if (invalidExtensions.test(inputFile)  ) {
        handleError("Invalid file selected")("Please select a valid kube config")
        return;
      }

      uploadK8SConfig().then((obj) => {
        handleConfigSnackbars(obj);
      }).
        catch(err => {
          handleError("failed to upload kubernetes config")(err)
        })
      formData.delete("k8sfile");
    }
  }

  return (
    <div>
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        onClick={handleClick}
        style={{
          "padding" : "8px",
          "borderRadius" : 5,
          "marginRight" : "2rem"
        }}
        data-cy="btnResetDatabase"
      >
        <AddIconCircleBorder style={{ width : "20px",height : "20px" }} />
        <Typography
          style={{
            "paddingLeft" : "4px" ,
            "marginRight" : "4px"
          }}
        > Add Cluster</Typography>
      </Button>
      <PromptComponent ref={ref} />
    </div>
  )
}

export default MesherySettingsEnvButtons
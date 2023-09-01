/* eslint-disable no-unused-vars, no-undef  */
import {
  FormGroup,
  TextField,
  InputAdornment,
  MenuItem,
  Grid,
} from "@material-ui/core/";
import BackupIcon from "@material-ui/icons/Backup";
import { useState } from "react";
import { fetchContexts, submitConfig } from "../helpers/kubernetesHelpers";
import { useNotification } from "../../../utils/hooks/useNotification";

const KubernetesConfig = ({
  updateK8SConfig, updateProgress
}) => {

  const [state, setState] = useState({
    contextNameForForm : "",
    contextsFromFile : [],
    k8sfile : "",
    k8sfileError : false,
    k8sfileElement : null,
    k8sfileElementVal : "",
    inClusterConfigForm : false,
  })
  const { notify } = useNotification();

  const handleChange = (name) => {
    return (event) => {
      let fileInput = document.querySelector("#k8sfile");
      let k8sfile = fileInput.files[0]
      if (name === "k8sfile") {
        if (event.target.value !== "") {
          setState({ ...state, k8sfileError : false });
        }
        fetchContexts(updateProgress, k8sfile )
          .then(res => {
            setState({
              ...state, contextsFromFile : res.result, k8sfile, contextNameForForm : res.currentContextName
            })
            if (res.result.length === 1)
              submitConfig(notify, updateProgress, updateK8SConfig, () => null, res.currentContextName, k8sfile);
          })
          .catch(err => alert(err))
      }
      if ( name === "contextNameChange"){
        submitConfig(notify, updateProgress, updateK8SConfig, action, event.target.value, k8sfile);
      }
    };
  };

  return (
    <>
      <Grid item xs={12} style={{ marginBottom : "0.6rem" }}>
        {/*<TextField
          id="k8sfileLabelText"
          name="k8sfileLabelText"
          label="Upload kubeconfig"
          variant="outlined"
          fullWidth
          margin="normal"
          InputProps={{
            readOnly: true,
            endAdornment: (
              <InputAdornment position="end">
                <BackupIcon style={{ marginRight: "1rem" }} />
              </InputAdornment>
            ),
          }}
          disabled
        />
      </Grid>
      <Grid item xs={12}>
        <TextField
          select
          id="contextName"
          name="contextName"
          label="Context Name"
          fullWidth
          margin="normal"
          variant="outlined"
          // disabled={inClusterConfigForm === true}
        >
          {false &&
                  contextsFromFile.map((ct) => (
                    <MenuItem key={`ct_---_${ct.contextName}`} value={ct.contextName}>
                      {ct.contextName}
                      {ct.currentContext ? " (default)" : ""}
                    </MenuItem>
                  ))}
        </TextField> */}

        <FormGroup>
          <input
            id="k8sfile"
            type="file"
            value={state.k8sfileElementVal}
            hidden
            onChange={handleChange("k8sfile")}
            // className={classes.fileInputStyle}
          />
          <TextField
            id="k8sfileLabelText"
            name="k8sfileLabelText"
            // className={classes.fileLabelText}
            label="Upload kubeconfig"
            variant="outlined"
            fullWidth
            value={""}
            onClick={() => document.querySelector("#k8sfile").click()}
            margin="normal"
            InputProps={{ readOnly : true,
              endAdornment : (
                <InputAdornment position="end">
                  <BackupIcon />
                </InputAdornment>
              ), }}
          />
        </FormGroup>
        <TextField
          select
          id="contextName"
          name="contextName"
          label="Context Name"
          fullWidth
          value={state.contextNameForForm}
          margin="normal"
          variant="outlined"
          // disabled={inClusterConfigForm === true}
          onChange={handleChange("contextNameChange")}
        >
          {state.contextsFromFile &&
                      state.contextsFromFile.map((ct) => (
                        <MenuItem key={`ct_---_${ct.contextName}`} value={ct.contextName}>
                          {ct.contextName}
                          {ct.currentContext
                            ? " (default)"
                            : ""}
                        </MenuItem>
                      ))}
        </TextField>

      </Grid>
    </>
  )
}


export default KubernetesConfig

/* eslint-disable no-unused-vars */
import {
  withStyles,
  FormGroup,
  TextField,
  InputAdornment,
  MenuItem,
  Grid,
} from "@material-ui/core/";
import BackupIcon from "@material-ui/icons/Backup";


const styles = () => ({
  contentBottomInputChecked: {
    background: "white",
    height: "6rem",
    width: "100%",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
  },
  contentBottomInputUnchecked: {
    display: "none",
  },
  contentBottomInput: {
    border: "1px solid lightgray",
    borderRadius: "5px",
    width: "9rem",
    height: "2rem",
    marginBottom: "0.15rem",
    fontSize: "0.75rem",
    padding: "0.50rem",
  },
  topInputIcon: {
    position: "absolute",
    fontSize: "1.25rem",
    color: "lightgray",
    bottom: "4.25rem",
    left: "9rem",
    cursor: "pointer",
    zIndex: "99999",
    "&:hover": {
      color: "grey",
    },
  },
  file: {
    display: "none",
  },
  // Inputs
  contentBottomChecked: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "flex-end",
  },
  contentBottomUpperInput: {
    width: "11rem",
    fontSize: "0.75rem",
    marginLeft: "2.4rem",
    marginTop: "-1rem",
    marginBottom: "0rem",
  },
  contentBottomLowerInput: {
    width: "11rem",
    marginBottom: "-1rem",
    fontSize: "0.75rem",
    marginTop: "0",
  },
})

const KubernetesConfig = ({classes}) => {
  return(
    <>
      <Grid item xs={12}>
        <TextField
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
        </TextField>
      </Grid>
    </>
  )
}


export default withStyles(styles)(KubernetesConfig)

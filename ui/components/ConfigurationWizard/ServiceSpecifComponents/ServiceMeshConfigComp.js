import {
  withStyles,
  TextField,
  MenuItem,
  Grid,
} from "@material-ui/core/";


const styles = () => ({
}) 

const ServiceMeshConfig = ({classes}) => {
  return(
    <>
       <Grid item xs={12}>
                <TextField
                  id="service-mesh-config-url-input"
                  name="serviceMeshConfigURL"
                  label="URL"
                  variant="outlined"
                  fullWidth
                  margin="normal"
                />
          </Grid>
          <Grid item xs={12}>
              <TextField
                select
                id="service-mesh-config-context-input"
                name="serviceMeshConfigContext"
                label="Context"
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


export default withStyles(styles)(ServiceMeshConfig)

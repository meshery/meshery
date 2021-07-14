/* eslint-disable no-unused-vars */
import {
  withStyles,
  TextField,
  MenuItem,
  Grid,
} from "@material-ui/core/";


const styles = () => ({
}) 

const ExternalsConfig = ({classes}) => {
  return(
    <>
      <Grid item xs={12}>
        <TextField
          id="external-config-url-input"
          name="externalConfigURL"
          label="URL"
          variant="outlined"
          fullWidth
          margin="normal"
        />
      </Grid>
      <Grid item xs={12}>
        <TextField
          select
          id="external-config-api-key-input"
          name="externalConfigapiKey"
          label="API key"
          fullWidth
          margin="normal"
          variant="outlined"
          // disabled={inClusterConfigForm === true}
        >
        </TextField>
      </Grid>
    </>
  )
}


export default withStyles(styles)(ExternalsConfig)

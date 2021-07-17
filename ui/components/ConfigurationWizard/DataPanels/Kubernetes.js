import {
  withStyles,
  Typography,
  Grid,
} from "@material-ui/core/";

const styles = theme => ({

    infoContainer: {
        width: "20rem",
        height: "15rem",
        padding: "1rem 1rem",
        boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.20)",
        borderRadius: '1rem',
      },
      infoTitle: {
        color: "#647881",
        width: "3rem",
        background: "#F1F3F4",
        padding: ".5rem 5rem .75rem 1.5rem",
        borderRadius: "0.25rem",
        fontSize: ".8rem",
      },
      infoLabel: {
        fontSize: ".9rem",
        color: theme.palette.text.primary,
        marginRight: "1rem",
      },
      infoData: {
        fontSize: ".9rem",
        color: theme.palette.text.secondary,
      },

})

const KubernetesDataPanel = ({clusterInformation, classes}) => (

       <Grid container className={classes.infoContainer} xs={12}>
         <Grid item xs={12} style={{marginBottom: "1rem"}}>
            <Typography className={classes.infoTitle}>Status</Typography>
         </Grid>
         <Grid item xs={12} container>
            <Grid item xs={12}>
              <Typography className={classes.infoLabel}>Current-Context:</Typography>
              <Typography className={classes.infoData}>
                {clusterInformation.inClusterConfig ? "Using In Cluster Config" : "Name"}
              </Typography>
            </Grid>
            <Grid item xs={12}>
              <Typography className={classes.infoLabel}>Cluster:</Typography>
              <Typography className={classes.infoData}>
                {clusterInformation.inClusterConfig ? "Using In Cluster Config" : "Using Out Of Cluster Config"}
              </Typography>
            </Grid>
          </Grid>
        </Grid>

)

export default withStyles(styles)(KubernetesDataPanel)

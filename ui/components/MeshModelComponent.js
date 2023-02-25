import {  withStyles } from '@material-ui/core'
import { withSnackbar } from 'notistack';
import React from 'react'
const meshmodelStyles = (theme) => ({
  wrapperClass : {
    padding : theme.spacing(2),backgroundColor : "white",height : "60vh",marginTop : theme.spacing(2),borderRadius : theme.spacing(1),
  },
  dashboardSection : {
    backgroundColor : "#fff",padding : theme.spacing(2),borderRadius : 4,height : "100%",overflowY : "scroll"
  },

})
const MeshModelComponent = ({ showMeshModelSummary,classes }) => {
  return (
    <div className={classes.wrapperClass}>
      <div className={classes.dashboardSection} data-test="workloads">
        {showMeshModelSummary()}
      </div>
    </div>
  )
}

export default withStyles(meshmodelStyles)((withSnackbar(MeshModelComponent)));

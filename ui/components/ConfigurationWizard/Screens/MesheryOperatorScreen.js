import MesheryOperatorIcon from "../icons/MesheryOperatorIcon.js"
import ServiceSwitch from "../ServiceSwitchCard.js"
import {Grid} from "@material-ui/core"


const MesheryOperatorScreen = () => {

  const mesheryOperatorinfo = {
    name: "Meshery Operator",
    logoComponent: MesheryOperatorIcon,
    configComp :  null
  }

  return (
    <Grid item container justify="center" alignItems="center" xs={12}>
      <ServiceSwitch serviceInfo={mesheryOperatorinfo} /> 
    </Grid>
  )
}


export default MesheryOperatorScreen

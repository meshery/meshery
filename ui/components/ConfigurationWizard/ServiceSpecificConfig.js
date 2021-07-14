import {Grid} from "@material-ui/core"


const ServiceSpecificConfig = ({components}) => {
  return (
    <Grid xs={12} container justify="center" alignItems="center" style={{padding: "1rem"}}>
    {components.map(Comp => (
        <Comp />
      ))} 
    </Grid>
  )
}


export default ServiceSpecificConfig

/* eslint-disable no-unused-vars */
import MesheryOperatorIcon from "../icons/MesheryOperatorIcon.js"
import fetchMesheryOperatorStatus from '../../graphql/queries/OperatorStatusQuery';
import subscribeOperatorStatusEvents from '../../graphql/subscriptions/OperatorStatusSubscription';
import ServiceCard from "../ServiceCard"
import {Grid} from "@material-ui/core"
import MesheryOperatorDataPanel from "../DataPanels/MesheryOperator"
import { useEffect, useState } from "react"
import { getOperatorStatusFromQueryResult, isMesheryOperatorConnected } from "../helpers/mesheryOperator.js";
 

const MesheryOperatorScreen = ({setStepStatus}) => {

  const [operatorInformation,setOperatorInformation] = useState({
    operatorInstalled: false,
    NATSInstalled: false,
    meshSyncInstalled: false,
    operatorVersion:"N/A",
    meshSyncVersion:"N/A",
    NATSVersion:"N/A",
  })
  const [isConnected,setIsConnected] = useState(false)

  const mesheryOperatorinfo = {
    name: "Meshery Operator",
    logoComponent: MesheryOperatorIcon,
    configComp :  <div />,
    operatorInformation
  }

  const showDataPanel = () => 
    isMesheryOperatorConnected(operatorInformation)

  useEffect(() => {
    setStepStatus(prev => ({...prev, operator: isConnected}))
  },[isConnected])

  useEffect(() => {

    subscribeOperatorStatusEvents(setOperatorState);
    fetchMesheryOperatorStatus().subscribe({
      next: (res) => {
        setOperatorState(res);
      },
      error: (err) => console.log("error at operator scan: " + err),
    });

  },[])

  useEffect(() => {
    setIsConnected(isMesheryOperatorConnected(operatorInformation))
  },[])

  const setOperatorState = (res) => {
    const [isReachable, operatorInformation] = getOperatorStatusFromQueryResult(res) 
    setOperatorInformation(operatorInformation)
  }
     
  return (
    <Grid item xs={12} container justify="center" alignItems="flex-start"> 
      <Grid item container justify="center" alignItems="center" xs={6} style={{paddingLeft: "1rem"}}>
        <ServiceCard serviceInfo={mesheryOperatorinfo} isConnected={isConnected} style={{paddingRight: "1rem"}}/> 
      </Grid>
      <Grid item xs={6} container justify="center">
        {
          showDataPanel() &&
        <MesheryOperatorDataPanel operatorInformation={operatorInformation}/>
        }
      </Grid>

    </Grid>
  )
}


export default MesheryOperatorScreen

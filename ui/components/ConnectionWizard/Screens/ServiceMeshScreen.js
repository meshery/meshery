/* eslint-disable no-unused-vars */
import ConsulIcon from "../icons/ConsulIcon.js"
import OpenServiceMeshIcon from "../icons/OpenServiceMeshIcon.js"
import ArrowForwardIosIcon from '@material-ui/icons/ArrowForwardIos';
import LinkerdIcon from "../icons/LinkerdIcon.js"
import { bindActionCreators } from "redux";
import { connect } from "react-redux";
import ServiceMeshConfig from "../ConfigComponents/ServiceMesh"
import ServiceCard from "../ServiceCard"
import {Grid} from "@material-ui/core"
import VerticalCarousel from "../../VerticalCarousel/VerticalCarousel"
import ServiceMeshDataPanel from "../DataPanels/ServiceMesh"
import {createRef, useEffect, useState} from "react"
import {fetchAvailableAdapters} from "../helpers/serviceMeshes"
import { updateProgress } from "../../../lib/store.js";
import {ScrollIndicator} from "../ScrollIndicator"


 
const ServiceMeshScreen = ({meshAdapters, meshAdaptersts, updateProgress}) => {



  const [availableAdapters, setAvailableAdapters] = useState([])
  const [activeAdapters, setActiveAdapters] = useState(meshAdapters)
  const [activeIndex, setActiveIndex] = useState(0);
  const sliderRef = createRef()
 

  const handleAfterSlideChange = (curSlide) => setActiveIndex(curSlide) 

  const serviceMeshComponents = availableAdapters.map(adapter => ({
    name: adapter.name ? adapter.name : adapter.label.split(":")[0],
    logoComponent:  OpenServiceMeshIcon,
    configComp :  <ServiceMeshConfig adapterLoc={adapter.value}/>,
    adapterInfo: adapter
  }))

  const isAdapterActive = (adapterLoc) => {
    let isActive = false
    activeAdapters.forEach(adapter => {
      if(adapter.adapter_location === adapterLoc) isActive = true
    }) 
    return isActive
  }


  const scrollItems = serviceMeshComponents.map(smesh => ({ activeIcon: "/static/img/meshery-logo.png", inactiveIcon:"/static/img/meshery-logo/meshery-white.png"}))


  const itemsToBeRendered = serviceMeshComponents.map(comp => {
    return(
      <ServiceCard serviceInfo={comp} isConnected={isAdapterActive(comp.adapterInfo.value)} /> 
    ) 
  })


  useEffect(() => {
    fetchAvailableAdapters()
      .then(res => {
        setAvailableAdapters( res)
      })
      .catch( err => alert(err))
  },[meshAdapters])

  useEffect(() => {
    setActiveAdapters(meshAdapters) 
  }, [meshAdapters])

  const handleIndicatorClick = (index) => () => {
     sliderRef?.current?.slickGoTo(index,false) 
  }



  return (
    <Grid xs={12} container justify="center" alignItems="flex-start">
      <Grid item xs={6} container justify="flex-start" >
        <div style={{height: "18rem", overflow: "scroll", marginTop: '1.3rem'}} className="hide-scrollbar"> 
          <ScrollIndicator items={scrollItems} handleClick={handleIndicatorClick} activeIndex={activeIndex} />
        </div>
        <VerticalCarousel slides={itemsToBeRendered} handleAfterSlideChange={handleAfterSlideChange} sliderRef={sliderRef}/>
      </Grid>
      <Grid item xs={6} container justify="center">
        <ServiceMeshDataPanel adapterInfo={serviceMeshComponents[activeIndex]?.adapterInfo} isActive={isAdapterActive(serviceMeshComponents[activeIndex]?.adapterInfo.value)}  />
      </Grid>
    </Grid>
  )
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  return {
    meshAdapters,
    meshAdaptersts,
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(ServiceMeshScreen)


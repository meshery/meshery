/* eslint-disable no-unused-vars */
/* eslint-disable react/display-name */
import LinkerdIcon from "../icons/LinkerdIcon.js"
import { bindActionCreators } from "redux";
import { connect } from "react-redux";
import ServiceMeshConfig from "../ConfigComponents/ServiceMesh"
import ServiceCard from "../ServiceCard"
import { Grid } from "@material-ui/core"
import VerticalCarousel from "../../VerticalCarousel/VerticalCarousel"
import ServiceMeshDataPanel from "../DataPanels/ServiceMesh"
import { createRef, useEffect, useState } from "react"
import { fetchAvailableAdapters } from "../helpers/serviceMeshes"
import { updateProgress } from "../../../lib/store.js";
import { ScrollIndicator } from "../ScrollIndicator"



const ServiceMeshScreen = ({ meshAdapters, meshAdaptersts, updateProgress }) => {



  const [availableAdapters, setAvailableAdapters] = useState([])
  const [activeAdapters, setActiveAdapters] = useState(meshAdapters)
  const [activeIndex, setActiveIndex] = useState(0);
  const sliderRef = createRef()


  const handleAfterSlideChange = (curSlide) => setActiveIndex(curSlide)

  const isAdapterActive = (adapterLoc) => {
    let isActive = false
    activeAdapters.forEach(adapter => {
      if (adapter.adapter_location === adapterLoc) isActive = true
    })
    return isActive
  }

  const AdapterIcon = (name) => ({ isActive }) => {

    let image = "/static/img/" + name?.toLowerCase() + ".svg";
    return  <img style={{ height : "4rem", width : "4rem" }} src={isActive ? image :"/static/img/meshery-logo/meshery-white.png"} />
  }

  const serviceMeshComponents = availableAdapters.map(adapter => ({
    name : adapter.name
      ? adapter.name
      : adapter.label.split(":")[0],
    logoComponent : AdapterIcon(adapter.name),
    configComp : <ServiceMeshConfig adapterLoc={adapter.value}/>,
    adapterInfo : adapter
  }))



  const scrollItems = serviceMeshComponents.map(smesh => ({ activeIcon : "/static/img/meshery-logo.png", inactiveIcon : "/static/img/meshery-logo/meshery-white.png" }))


  const itemsToBeRendered = serviceMeshComponents.map(comp => {
    return (
      <ServiceCard serviceInfo={comp} isConnected={isAdapterActive(comp.adapterInfo.value)} />
    )
  })


  useEffect(() => {
    fetchAvailableAdapters()
      .then(res => {
        setAvailableAdapters(res.sort((fe,se) => isAdapterActive(se.value)
          ? 1
          : -1))
      })
      .catch( err => alert(err))
  },[meshAdapters])

  useEffect(() => {
    setActiveAdapters(meshAdapters)
  }, [meshAdapters])

  const handleIndicatorClick = (index) => (e) => {
    e.preventDefault()
    sliderRef?.current?.slickGoTo(index,false)
    setActiveIndex(index)
  }


  return (
    <Grid xs={12} container justify="center" alignItems="flex-start">
      <Grid item lg={6} sm={12} md={12} container justify="center" alignItems="flex-start" style={{ paddingLeft : "1rem" }}>
        <div style={{ height : "18rem", overflow : "scroll", marginTop : "-1.2rem" }} className="hide-scrollbar">
          <ScrollIndicator items={scrollItems} handleClick={handleIndicatorClick} activeIndex={activeIndex} />
        </div>
        <VerticalCarousel slides={itemsToBeRendered} handleAfterSlideChange={handleAfterSlideChange} sliderRef={sliderRef}/>
      </Grid>
      <Grid item lg={6} sm={12} md={12} container justify="center" style={{ paddingRight : "1rem" }}>
        <ServiceMeshDataPanel adapterInfo={serviceMeshComponents[activeIndex]?.adapterInfo} isActive={isAdapterActive(serviceMeshComponents[activeIndex]?.adapterInfo.value)}  />
      </Grid>
    </Grid>
  )
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

const mapStateToProps = (state) => {
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  return { meshAdapters,
    meshAdaptersts, };
};

export default connect(mapStateToProps, mapDispatchToProps)(ServiceMeshScreen)


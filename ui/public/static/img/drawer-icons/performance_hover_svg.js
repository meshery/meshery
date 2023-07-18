import ArrowDropDownIcon from "@material-ui/icons/ArrowDropDown";
import PerformanceIcon from "./performance_svg";
const drawerIconsStyle = { height: "1.21rem", width: "1.21rem", fontSize: "1.0rem", transform: "scale(1.5)" };
const PerformanceHover = () => (
  <>
    <PerformanceIcon style={drawerIconsStyle} />
    <ArrowDropDownIcon />
  </>
);

export default PerformanceHover;

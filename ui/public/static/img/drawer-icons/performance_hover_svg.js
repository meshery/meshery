import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import PerformanceIcon from "./performance_svg";
const drawerIconsStyle = { height: "19.36px", width: "19.36px", fontSize: "1.0rem", transform: "scale(1.5)" };
const PerformanceHover = () => (
  <>
    <PerformanceIcon style={drawerIconsStyle} />
    <ArrowDropDownIcon />
  </>
);

export default PerformanceHover;

import { styled } from "@mui/material/styles";

export const StyledDiv = styled('div')(() => ({
  width: "16%",
  flexDirection: "column",
  display: "flex",
  padding: "0.3rem"
}));


export const AccountDiv = styled('div')(() => ({
  width: "50%",
  
  display: "flex",
  flexDirection: "column",
  justifyContent: 'center', alignItems: 'center' ,
}));

export const ServiceMeshAdapters = styled('div')(({ theme }) => ({
  display: "flex",
  justifyContent: 'center', alignItems: 'center',
  alignContent: 'space-between',
  [theme.breakpoints.down("md")]: {
    width: "100%"
  },
}));

export const AdapterDiv = styled("div")(({ theme, inactiveAdapter }) => ({
  filter: inactiveAdapter ? "grayscale(1) invert(0.35)" : ""
}))

// export const SpinnerDiv = styled("div")(({noHover}) => ({
//   "&:hover {   
//     transform: scale(2),
//     transitionDuration: "0.5s" "
// }))

export const ExtensionWrapper = styled('div')(({ theme }) => ({
  margin: theme.spacing(2),
  display: "flex",
  justifyContent: 'center', alignItems: 'center',
  height:"12rem",
  backgroundColor: "#393F49",
  borderRadius: "20px ",
  padding: "2rem",
  textAlign: "center",
  [theme.breakpoints.down("md")]: {
    width: "100%",
  }
}));

export const ComponentWrapper = styled('div')(({theme }) => ({
  textAlign: "center",
    backgroundColor: "#222C32",
    padding: "5rem",
    maxHeight: "100vh"
}))

export const SectionWrapper = styled('div')(({theme }) => ({
  padding: "2rem", display: "flex", flexDirection: 'row',
  justifyContent: 'center', alignItems: 'center' 
}))
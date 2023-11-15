
import { styled } from "@mui/material/styles";
import { ButtonBase, Typography } from "@mui/material";

export const StyledDiv = styled('div')(() => ({
  paddingLeft: "0.2rem",
  paddingRight: "0.2rem"
}));


export const AccountDiv = styled('div')(() => ({
  width: "50%",
  display: "flex",
  flexDirection: "column",
  justifyContent: 'center', alignItems: 'center',
}));

export const ServiceMeshAdapters = styled('div')(({ theme }) => ({
  display: "flex",
  alignItems: 'center',
  [theme.breakpoints.down("md")]: {
    display: "flex",
    flexWrap: "wrap",
  },
}));

export const AdapterDiv = styled("div")(({ theme, inactiveAdapter }) => ({
  filter: inactiveAdapter ? "grayscale(1) invert(0.35)" : ""
}))


export const ExtensionWrapper = styled('div')(({ theme }) => ({
  margin: theme.spacing(2),
  display: "flex",
  justifyContent: 'center',
  alignItems: 'center',
  height: "14rem",
  borderRadius: "20px ",
  padding: "2rem",
  textAlign: "center",
}));

export const ComponentWrapper = styled('div')(({ theme }) => ({
  textAlign: "center",
  padding: "3rem 5rem",
  maxHeight: "100vh",
}))

export const SectionWrapper = styled('div')(({ theme }) => ({
  padding: "2rem", display: "flex", flexWrap: "wrap",
  justifyContent: 'center', alignItems: 'center',
}))

export const LoadingDiv = styled('div')(({ theme }) => ({
  top: "40%",
  left: "43%",
  position: "absolute",

  zIndex: 10,
}))  

export const VersionDiv = styled('div')(({theme}) => ({
  position: "relative",
  bottom: 0,
  right: 0,
}))

export const VersionText = styled(Typography)(({ theme }) => ({
  padding: "0rem 0.5rem",
  paddingBottom: "0",
  width: "fit-content",
  marginLeft: "auto",
}));

export const LogoutButton = styled('div')(({ theme }) => ({
  transform:" translateX(39%)",
  width: "fit-content",
}));

export const StyledButton = styled(ButtonBase)(() => ({
  marginTop: '1rem',
  whiteSpace: 'nowrap',
  backgroundColor: "#00B39F",
  borderRadius: "5px",
  color: "white",
  padding: "10px 20px",
  "&:hover": {
    boxShadow: "0 1px 2px 0 rgba(0, 0, 0, 0.2), 0 1px 5px 0 rgba(0, 0, 0, 0.19)",
  },
}));

export const LinkButton = styled(ButtonBase)(() => ({
  marginTop: '1rem',
  whiteSpace: 'nowrap',
  padding: "0",
  backgroundColor: "#00B39F",
  borderRadius: "5px",
  color: "white",
  "&:hover": {
    boxShadow: "0 1px 2px 0 rgba(0, 0, 0, 0.2), 0 1px 5px 0 rgba(0, 0, 0, 0.19)",
  },
}));

export const StyledLink = styled('a')(() => ({
  padding: "10px 20px",
  borderRadius: "5px",
  color: "white",
}));
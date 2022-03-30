import { styled } from "@mui/material/styles";
//     inactiveAdapter: {
//       filter: "grayscale(1) invert(0.35)"
//     },
export const StyledDiv = styled('div')(() => ({
  width: "16%",
  float: "left",
  flexDirection: "row",
  padding: "0.3rem"
}));


export const AccountDiv = styled('div')(() => ({
  width: "50%",
  float: "left",
}));

export const ServiceMeshAdapters = styled('div')(({ theme }) => ({
  width: "50%",
  float: "right",
  [theme.breakpoints.down("xs")]: {
    width: "100%"
  },
}));

export const AdapterDiv = styled("div")(({ theme, inactiveAdapter }) => ({
  filter: inactiveAdapter ? "grayscale(1) invert(0.35)" : ""
}))

export const ExtensionWrapper = styled('div')(({ theme }) => ({
  margin: theme.spacing(5),
  backgroundColor: "#393F49",
  borderRadius: "20px ",
  padding: "1rem",
  height: "300px",
  [theme.breakpoints.down("xs")]: {
    height: "400px",
  }
}));

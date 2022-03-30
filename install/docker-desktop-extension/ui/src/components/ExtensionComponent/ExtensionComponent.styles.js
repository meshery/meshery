import makeStyles from '@mui/styles/makeStyles'
import { useTheme } from '@mui/material/styles';

export const useStyles = makeStyles((theme) => ({
    root: {
      textAlign: "center",
      backgroundColor: "#222C32",
      padding: "5rem",
      maxHeight: "100vh"
    },
    main: {
      margin: useTheme().spacing(5),
      backgroundColor: "#393F49",
      borderRadius: "20px ",
      padding: "1rem",
      height: "300px",
      [useTheme().breakpoints.down("xs")]: {
        height: "400px",
      },
    },
    paper: {
      padding: useTheme().spacing(1.5),
      textAlign: "center",
      color: "#ffffff",
      width: "240px",
      height: "45px"
    },
    OAuth: {
      padding: "2rem",
  
    },
    serviceMeshAdapters: {
      width: "50%",
      float: "right",
      [useTheme().breakpoints.down("xs")]: {
        width: "100%"
      }
    },
    account: {
      width: "50%",
      float: "left",
    },
    sm: {
      width: "16%",
      float: "left",
      flexDirection: "row",
      padding: "0.3rem"
    },
    subText: {
      color: "#AAAAAA",
      padding: "0.7rem"
    },
    headText: {
      maxWidth: "60%",
      margin: "auto",
      padding: "1rem"
    },
    button: {
      padding: "0.5rem"
    },
    inactiveAdapter: {
      filter: "grayscale(1) invert(0.35)"
    },
    mesheryConfig: {
      backgroundColor: "#7794AB",
      color: "#FFFFFF",
    },
    link: {
      textDecoration: "none"
    },
    Icon: {
      width: useTheme().spacing(2.5),
      paddingRight: useTheme().spacing(0.5),
    },
  }));

import React from "react";
import { useStyles } from "./Footer.styles";
import FavoriteIcon from "@mui/icons-material/Favorite";
import { Typography } from "@mui/material";
const FooterComponent = () => {
  const classes = useStyles();

  const handleComClick = () => {
    if (typeof window !== "undefined") {
      const L = window.open("https://layer5.io", "_blank");
      L.focus();
    }
  };
  return (
    <footer className={classes.footer}>
      <Typography variant="body2" align="center" color="textSecondary">
        <span onClick={handleComClick} className={classes.footerText}>
          Built with
          <FavoriteIcon className={classes.footericon} /> by the Layer5 Community
        </span>
      </Typography>
    </footer>
  );
};
export default FooterComponent;

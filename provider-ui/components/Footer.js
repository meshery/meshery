import React, { Fragment } from "react";
import { styled } from "@mui/material/styles";
import Paper from "@mui/material/Paper";
import NoSsr from "@mui/material/NoSsr";
import Typography from "@mui/material/Typography";
import FavoriteIcon from "@mui/icons-material/Favorite";

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor : "rgb(54, 54, 54)",
  padding : theme.spacing(2),
  color : "rgba(255, 255, 255, 0.7)",
  minWidth : "100%",
  border : "0",
}));

export default function Footer() {
  const handleL5CommunityClick = () => {
    if (typeof window !== "undefined") {
      const w = window.open("https://layer5.io", "_blank");
      w.focus();
    }
  };

  return (
    <Fragment>
      <NoSsr>
        <Item component="footer" square variant="outlined">
          <Typography
            variant="body2"
            align="center"
            color="textSecondary"
            component="p"
          >
            <span
              onClick={handleL5CommunityClick}
              style={{
                cursor : "pointer",
                display : "inline",
                verticalAlign : "middle",
                color : "rgba(255, 255, 255, 0.7)",
              }}
            >
              Built with{" "}
              <FavoriteIcon
                sx={{
                  display : "inline",
                  verticalAlign : "top",
                  color : "rgb(0, 179, 159)",
                }}
              />{" "}
              by the Layer5 Community
            </span>
          </Typography>
        </Item>
      </NoSsr>
    </Fragment>
  );
}

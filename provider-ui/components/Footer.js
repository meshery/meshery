import React, { Fragment } from "react";
import { NoSsr } from '@sistent/sistent'
import { FavoriteIcon, Typography, Paper, styled, accentGrey } from "@sistent/sistent";
const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: accentGrey[20],
  padding: theme.spacing(2),
  color: theme.palette.background.secondary,
  minWidth: "100%",
  border: "0",
}));

const ClickableSpan = styled("span")(({ theme }) => ({
  cursor: "pointer",
  display: "inline",
  verticalAlign: "middle",
  color: theme.palette.text.disabled,
}));

const StyledFavoriteIcon = styled(FavoriteIcon)(({ theme }) => ({
  display: "inline",
  verticalAlign: "top",
  fill: theme.palette.icon.brand
}));

export default function Footer() {
  const handleMesheryCommunityClick = () => {
    if (typeof window !== "undefined") {
      const w = window.open("https://meshery.io", "_blank");
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
            <ClickableSpan onClick={handleMesheryCommunityClick}>
              Built with <StyledFavoriteIcon sx={{ color: "blue" }} /> by the Meshery Community
            </ClickableSpan>
          </Typography>
        </Item>
      </NoSsr>
    </Fragment>
  );
}

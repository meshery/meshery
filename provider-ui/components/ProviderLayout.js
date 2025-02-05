import React from "react";
import PropTypes from "prop-types";
import NoSsr from "@mui/material/NoSsr";
import { styled } from "@layer5/sistent";

const RootContainer = styled("div")(({ theme }) => ({
  padding : "170px 0px",
  textAlign : "center",
  backgroundColor : theme.palette.background.elevatedComponents,
}));

export default function ProviderLayout({ children }) {
  return (
    <>
      <NoSsr>
        <RootContainer data-cy="root">{children}</RootContainer>
      </NoSsr>
    </>
  );
}

ProviderLayout.propTypes = {
  children : PropTypes.node.isRequired,
};

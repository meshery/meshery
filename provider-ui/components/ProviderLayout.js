import React from "react";
import PropTypes from "prop-types";
import NoSsr from "@mui/material/NoSsr";
import { accentGrey, styled } from "@layer5/sistent";

const RootContainer = styled("div")(() => ({
  padding : "170px 0px",
  textAlign : "center",
  backgroundColor : accentGrey[20]
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

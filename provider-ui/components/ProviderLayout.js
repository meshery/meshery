import React from "react";
import PropTypes from "prop-types";
import { NoSsr } from '@layer5/sistent'
import { accentGrey, styled } from "@layer5/sistent";

const RootContainer = styled("div")(() => ({
  padding: "4vh 12vw",
  borderRadius: ".5rem",
  textAlign: "center",
  backgroundColor: accentGrey[20],
  margin: "auto",
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
  children: PropTypes.node.isRequired,
};

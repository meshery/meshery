import React from "react";
import PropTypes from "prop-types";
import { NoSsr } from '@sistent/sistent'
import { accentGrey, styled } from "@sistent/sistent";

const RootContainer = styled("div")(() => ({
  padding: "4vh 12vw",
  borderRadius: ".5rem",
  textAlign: "center",
  backgroundColor: accentGrey[20],
  margin: "auto",
  /* Deep Charcoal base */
  // backgroundColor: "rgb(17, 22, 25)",
  /* Multi-variant radial blend */
  /* Meshery Keppel */
  /* Bright Keppel Accent */
  /* Slate Blue Accent */
  backgroundImage: "radial-gradient(circle at 50% 35%, rgba(17, 22, 25, 1) 5%, transparent 100%)",
  backgroundSize: "cover",
  backgroundPosition: "center",
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

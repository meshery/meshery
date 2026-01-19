import React from 'react';
import PropTypes from 'prop-types';

const GrafanaPanelIframe = ({ src, title }) => <iframe src={src} title={title} loading="lazy" />;

GrafanaPanelIframe.propTypes = {
  src: PropTypes.string.isRequired,
  title: PropTypes.string.isRequired,
};

export default GrafanaPanelIframe;

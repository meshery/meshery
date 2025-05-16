import { FALLBACK_KUBERNETES_IMAGE_PATH, KUBERNETES } from '@/constants/common';
import { componentIcon } from '@layer5/sistent';
import { iconXLarge } from 'css/icons.styles';
import React from 'react';

const GetKubernetesNodeIcon = ({ kind, size, model }) => {
  const imgSrc = componentIcon({
    kind: kind?.toLowerCase(),
    color: 'color',
    model: model ? model : KUBERNETES,
  });

  return (
    <img
      src={imgSrc}
      onError={(event) => {
        event.target.src = FALLBACK_KUBERNETES_IMAGE_PATH;
      }}
      alt={kind}
      style={size ? size : iconXLarge}
    />
  );
};

export default GetKubernetesNodeIcon;

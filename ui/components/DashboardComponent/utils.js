import {
  CUSTOM_RESOURCE_DEFINITION,
  FALLBACK_KUBERNETES_IMAGE_PATH,
  KUBERNETES,
} from '@/constants/common';
import { componentIcon } from '@layer5/sistent';
import { iconXLarge } from 'css/icons.styles';
import React from 'react';

const GetKubernetesNodeIcon = ({ kind, isCRD, size }) => {
  return (
    <img
      src={componentIcon({
        kind: isCRD ? CUSTOM_RESOURCE_DEFINITION : kind?.toLowerCase(),
        color: 'color',
        model: KUBERNETES,
      })}
      onError={(event) => {
        event.target.src = FALLBACK_KUBERNETES_IMAGE_PATH;
      }}
      alt={kind}
      style={size ? size : iconXLarge}
    />
  );
};

export default GetKubernetesNodeIcon;

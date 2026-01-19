import { FALLBACK_KUBERNETES_IMAGE_PATH, KUBERNETES } from '@/constants/common';
import { normalizeStaticImagePath } from '@/utils/fallback';
import { componentIcon } from '@sistent/sistent';
import { iconXLarge } from 'css/icons.styles';
import React from 'react';

const failedIconSources = new Set();

const GetKubernetesNodeIcon = ({ kind, size, model }) => {
  const imgSrc = componentIcon({
    kind: kind?.toLowerCase(),
    color: 'color',
    model: model ? model : KUBERNETES,
  });
  const normalizedImgSrc = normalizeStaticImagePath(imgSrc);
  const resolvedImgSrc =
    normalizedImgSrc && !failedIconSources.has(normalizedImgSrc)
      ? normalizedImgSrc
      : FALLBACK_KUBERNETES_IMAGE_PATH;

  return (
    <img
      src={resolvedImgSrc}
      onError={(event) => {
        if (normalizedImgSrc && normalizedImgSrc !== FALLBACK_KUBERNETES_IMAGE_PATH) {
          failedIconSources.add(normalizedImgSrc);
        }
        event.currentTarget.onerror = null;
        event.currentTarget.src = FALLBACK_KUBERNETES_IMAGE_PATH;
      }}
      alt={kind}
      style={size ? size : iconXLarge}
    />
  );
};

export default GetKubernetesNodeIcon;

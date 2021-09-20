/* eslint-disable react/prop-types */
import React from "react";
import Chip from "@/components/Chip";
import { useTheme } from "@mui/system";

/**
 *
 * @param {{grafana: import("../../mesheryEnvironmentSlice").GrafanaType}} props
 * @returns
 */
export const GrafanaChip = ({ grafana, handleClick }) => {
  const theme = useTheme();
  return (
    <Chip
      label={grafana.grafanaURL}
      onClick={handleClick}
      icon={<img src="/static/img/grafana_icon.svg" height={theme.spacing(2.5)} width={theme.spacing(2.5)} />}
      key="graf-key"
      variant="outlined"
    />
  );
};

/**
 *
 * @param {{prometheus: import("../../mesheryEnvironmentSlice").PrometheusType}} props
 * @returns
 */
export const PrometheusChip = ({ prometheus, handleClick }) => {
  const theme = useTheme();
  return (
    <Chip
      label={prometheus.prometheusURL}
      onClick={handleClick}
      icon={
        <img
          src="/static/img/prometheus_logo_orange_circle.svg"
          height={theme.spacing(2.5)}
          width={theme.spacing(2.5)}
        />
      }
      key="prom-key"
      variant="outlined"
    />
  );
};

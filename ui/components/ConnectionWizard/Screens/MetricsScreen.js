/* eslint-disable no-unused-vars */
import GrafanaIcon from '../icons/GrafanaIcon';
import PrometheusIcon from '../icons/PrometheusIcon';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import MetricsConfig from '../ConfigComponents/Metrics.js';
import ServiceCard from '../ServiceCard.js';
import { Grid } from '@material-ui/core';
import VerticalCarousel from '../../VerticalCarousel/VerticalCarousel';
import MetricsDataPanel from '../DataPanels/Metrics';
import { createRef, useEffect, useRef, useState } from 'react';
import { updateGrafanaConfig, updateProgress, updatePrometheusConfig } from '../../../lib/store.js';
import {
  fetchPromGrafanaScanData,
  verifyGrafanaConnection,
  verifyPrometheusConnection,
} from '../helpers/metrics';
import { ScrollIndicator } from '../ScrollIndicator';

const MetricsScreen = ({ grafana, prometheus, selectedK8sContexts }) => {
  const [isGrafanaConnected, setIsGrafanaConnected] = useState(false);
  const [isPrometheusConnected, setIsPrometheusConnected] = useState(false);
  const [metricsScanUrls, setMetricsScanUrls] = useState({ grafana: [], prometheus: [] });
  const [activeIndex, setActiveIndex] = useState(0);
  const sliderRef = createRef();

  const handleIndicatorClick = (index) => () => {
    sliderRef?.current?.slickGoTo(index, false);
    setActiveIndex(index);
  };

  const metricsComponents = [
    {
      name: 'Prometheus',
      logoComponent: PrometheusIcon,
      configComp: (
        <MetricsConfig
          componentName="Prometheus"
          prometheusScannedUrls={metricsScanUrls.prometheus}
        />
      ),
    },
    {
      name: 'Grafana',
      logoComponent: GrafanaIcon,
      configComp: (
        <MetricsConfig componentName="Grafana" grafanaScannedUrls={metricsScanUrls.grafana} />
      ),
    },
  ];

  const handleAfterSlideChange = (curSlide) => setActiveIndex(curSlide);

  const getConnectionStatus = (name) => {
    if (name === 'Grafana') return isGrafanaConnected;
    if (name === 'Prometheus') return isPrometheusConnected;
  };

  useEffect(() => {
    console.log('Prometheus effect', prometheus);
    verifyPrometheusConnection(prometheus.prometheusURL)
      .then((res) => {
        console.log(res);
        if (typeof res !== 'undefined') setIsPrometheusConnected(true);
        else setIsPrometheusConnected(false);
      })
      .catch((err) => {
        setIsPrometheusConnected(false);
        console.log(err);
      });
  }, [prometheus.ts]);

  useEffect(() => {
    console.log('Greafana effect', grafana);
    verifyGrafanaConnection(grafana.grafanaURL)
      .then((res) => {
        if (typeof res !== 'undefined') setIsGrafanaConnected(true);
        else setIsGrafanaConnected(false);
      })
      .catch((err) => setIsGrafanaConnected(false));
  }, [grafana.ts]);

  useEffect(() => {
    fetchPromGrafanaScanData(selectedK8sContexts)
      .then((res) => setMetricsScanUrls(res))
      .catch(console.log);
  }, []);

  const scrollItems = metricsComponents.map((metricComp) => {
    if (metricComp.name === 'Grafana')
      return {
        activeIcon: '/static/img/grafana_icon.svg',
        inactiveIcon: '/static/img/grafana_icon.svg',
      };
    if (metricComp.name === 'Prometheus')
      return {
        activeIcon: '/static/img/prometheus_logo_orange_circle.svg',
        inactiveIcon: '/static/img/prometheus_logo_orange_circle.svg',
      };
  });

  const itemsToBeRendered = metricsComponents.map((comp) => {
    return (
      <ServiceCard
        serviceInfo={comp}
        isConnected={getConnectionStatus(comp.name)}
        key={comp.uniqueID}
      />
    );
  });

  return (
    <Grid xs={12} item justify="center" alignItems="flex-start" container>
      <Grid
        item
        lg={6}
        sm={12}
        md={12}
        container
        justify="center"
        alignItems="flex-start"
        style={{ paddingLeft: '1rem' }}
      >
        <div
          style={{ height: '18rem', overflow: 'scroll', marginTop: '-1.2rem' }}
          className="hide-scrollbar"
        >
          <ScrollIndicator
            items={scrollItems}
            handleClick={handleIndicatorClick}
            activeIndex={activeIndex}
          />
        </div>
        <VerticalCarousel
          slides={itemsToBeRendered}
          handleAfterSlideChange={handleAfterSlideChange}
          sliderRef={sliderRef}
        />
      </Grid>
      <Grid item lg={6} sm={12} md={12} container justify="center" style={{ paddingRight: '1rem' }}>
        <MetricsDataPanel
          isConnected={
            metricsComponents[activeIndex].name === 'Grafana'
              ? isGrafanaConnected
              : isPrometheusConnected
          }
          componentName={metricsComponents[activeIndex].name}
        />
      </Grid>
    </Grid>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig: bindActionCreators(updatePrometheusConfig, dispatch),
});

const mapStateToProps = (state) => {
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    grafana,
    prometheus,
    selectedK8sContexts,
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(MetricsScreen);

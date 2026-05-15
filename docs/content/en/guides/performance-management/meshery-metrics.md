---
title: Using Metrics in Meshery
description: How to connect and use Prometheus and Grafana metrics in Meshery
categories: [performance]
---

## Connect and use metrics in Meshery

Meshery provides performance reports, including performance test results, node resource metrics etc. so that operators may easily understand the overhead of their cloud native infrastructure control plane and data plane in context of the overhead incurred on nodes running within the cluster. In order to generate performance test reports of infrastructure and their workloads, Meshery uses Grafana and/or Prometheus as visualization and metrics systems, respectively. This guide outlines the requirements necessary for Meshery to connect to these systems. The steps may vary depending upon the infrastructure and its configuration.

In order to pull in these environment metrics, you can also manually configure Meshery to connect with your existing Grafana and/or Prometheus instances through the Meshery dashboard. Once they have been loaded and are displayed on the screen, you may also perform an _ad-hoc_ test to check Meshery's connection status.

<main>
  
  <input id="tab1" type="radio" name="tabs" checked>
  <label for="tab1">Prometheus Metrics</label>
    
  <input id="tab2" type="radio" name="tabs">
  <label for="tab2">Grafana Charts</label>
    
  <input id="tab3" type="radio" name="tabs">
  <label for="tab3">Static Boards</label>
    
  <input id="tab4" type="radio" name="tabs">
  <label for="tab4">Dynamic Boards</label>
    
  <section class="tabbed" id="content1">
    <p>User needs to set the Prometheus URL and API key to create and query boards.</p>
    <img src="/guides/performance-management/images/PrometheusCharts.svg" 
      alt="Prometheus Metrics in Meshery" />

  </section>
    
  <section class="tabbed" id="content2">
    <p>User needs to set the Grafana URL and API key to create and query boards.</p>
    <img src="/guides/performance-management/images/GrafanaBoards.svg"
      alt="Grafana Charts in Meshery" />

  </section>
    
  <section class="tabbed" id="content3">
    <p>
      Static Boards capture the standard performance metrics included in every Meshery performance test irrespective of which load generator is used. A set calculations are made using statistical analysis of the metrics gathered in the static boards. Static Boards Queries Prometheus SDK directly.
    </p>
  </section>
    
  <section class="tabbed" id="content4">
    <p>
      Dynamic Boards can be generated from Prometheus or Grafana. These boards are defined by the user. Grafana SDK is used for these boards.
    </p>

  </section>
    
</main>

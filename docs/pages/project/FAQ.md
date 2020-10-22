---
layout: page
title: FAQ
permalink: project/faq
---

# FAQ 
<div class="center" style="color:gray;position:relative;top:-10px;font-size:1.25em;"> Can't find your question here?<br />Post it on the <a href="https://layer5io.slack.com/archives/C01AFD2D547">#support channel</a>
on the <a hre="http://slack.layer5.io/">Layer5 Slack</a>. </div>


## Why use Meshery?
* Because its an open source, vendor neutral projects that facilitates testing across meshes.
* Because fortio is not packaged into a mesh testing utility, but is only a load-generator unto its own.
* Because regpatrol is closed source, binary is not released, scripted for one mesh, and is produced by a vendor of that mesh.

## Why create Meshery and not use another benchmark tool?
Meshery is purpose built for facilitating benchmarking of service meshes and their workloads. Other benchmark tools are not. There are some other tools used for service mesh benchmarking, like regpatrol. Regpatrol is used by IBM is not open source or available in binary form to use and has the following differences from Meshery:
- Telemetry - regpatrol sources telemetry from the Mixer Prometheus adapter and uses IBM's proprietary node agent.
- Meshery sources from the Mixer Prometheus adapter and uses Prometheus node-exporter.
- Traffic type - regpatrol uses jmeter, which can parse responses and perform functional tests.
- Meshery is using fortio, which is for load-gen and perf-testing only.


<!--Add other questions-->
---
layout: enhanced
title: Component Shape Guide
permalink: extensions/component-shape-guide
language: en
abstract: Kubernetes architecture deployment and architecture diagramming tool for cloud native applications - MeshMap.
display-title: "false"
list: include
type: extensions
category: meshmap
---

## Component Shape Guide

Inside MeshMap, the allocation of specific shapes to signify various purposes creates a coherent and intelligible visual representation of intricate designs.
Currently, the circle is used as the default shape for new components.

Although the usage of the components is divided into categories, some shapes serve as a universal representation of particular components.

Below are all the shapes with their current usage in a general context.

<div class="svg-grid" style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin: 40px;">

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="68" viewBox="0 0 150 136" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M98.0102 15.75C87.8104 -1.91667 62.3107 -1.91668 52.1109 15.75L7.07755 93.75C-3.1223 111.417 9.62751 133.5 30.0272 133.5H120.094C140.494 133.5 153.243 111.417 143.044 93.75L98.0102 15.75Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Service</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 151 150" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="2.56055" y="2.5" width="145" height="145" rx="22.5" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Pod</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="75" viewBox="0 0 151 150" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M43.763 140.376C7.65635 123.089 -7.60209 79.8012 9.68328 43.6899C17.9841 26.3486 32.8327 13.0155 50.9628 6.62363C69.0928 0.231726 89.0191 1.30463 106.358 9.6062C142.465 26.8934 157.723 70.1814 140.438 106.293C123.152 142.404 79.8698 157.664 43.763 140.376Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Container</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="56" viewBox="0 0 151 113" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="2.56055" y="2.5" width="145" height="107.5" fill="white" fill-opacity="0.3" stroke="#666666" stroke-width="5" stroke-dasharray="10 10"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Namespace</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="75" viewBox="0 0 114 151" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="3.31055" y="3" width="107.5" height="145" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Node</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="75" viewBox="0 0 151 151" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="2.56055" y="3" width="145" height="145" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Cluster</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="72" viewBox="0 0 160 155" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M95.6369 7.72668C86.349 0.978614 73.7722 0.978623 64.4843 7.72669L13.8383 44.5231C4.55039 51.2712 0.663959 63.2325 4.21163 74.1511L23.5567 133.689C27.1043 144.608 37.2792 152 48.7597 152H111.362C122.842 152 133.017 144.608 136.565 133.689L155.91 74.151C159.457 63.2324 155.571 51.2712 146.283 44.5231L95.6369 7.72668Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Deployment</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="61" viewBox="0 0 153 123" fill="white" fill-opacity="0.5" xmlns="http://www.w3.org/2000/svg">
<path d="M144.25 46.2377L112.401 10.7933C107.944 5.8331 101.59 3 94.9215 3H26.0605C13.0819 3 2.56055 13.5213 2.56055 26.5V97C2.56055 109.979 13.0818 120.5 26.0605 120.5H94.9791C101.61 120.5 107.932 117.699 112.387 112.787L144.178 77.731C152.273 68.8045 152.304 55.2013 144.25 46.2377Z" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">ConfigMap</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="78" viewBox="0 0 132 156" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M128.766 32.7442C128.751 30.2268 127.557 28.3686 125.291 27.3159C123.649 26.5561 121.966 25.8909 120.291 25.18L98.4409 16.026C88.553 11.8885 78.6681 7.74784 68.7863 3.60416C67 2.85049 65.271 2.76811 63.4575 3.54924C59.0899 5.42579 54.692 7.23522 50.3092 9.07211L29.7103 17.6982C22.2546 20.8227 14.808 23.9452 7.37042 27.0657C4.69543 28.1855 3.35642 30.1505 3.3534 33.0676C3.3534 39.8812 3.37758 46.6947 3.3534 53.5113C3.30806 62.0031 3.31108 70.4888 4.86469 78.8891C5.64182 83.2089 6.69179 87.4742 8.00818 91.6588C12.3122 105.21 19.4835 117.655 29.0252 128.131C38.5668 138.607 50.2507 146.863 63.2671 152.328C64.8781 153.011 66.5042 153.292 68.1757 152.6C76.6233 149.142 84.5435 144.494 91.7005 138.795C99.4826 132.603 106.283 125.248 111.867 116.985C119.741 105.359 125.028 92.1506 127.364 78.2666C128.116 74.0118 128.56 69.7071 128.691 65.387C128.806 60.1785 128.718 54.9638 128.718 49.7521H128.766C128.766 44.0828 128.797 38.4135 128.766 32.7442Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5" stroke-miterlimit="10"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Secret</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="75" viewBox="0 0 124 150" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M3.3606 45.4019C6.77162 49.7691 11.733 53.5667 17.6307 56.6388C29.1298 62.6288 44.8421 66.2589 62.0606 66.2589C79.279 66.2589 94.9914 62.6288 106.49 56.6388C112.388 53.5667 117.35 49.7691 120.761 45.4019V118.12C120.761 125.52 114.953 132.834 104.181 138.445C93.5296 143.994 78.642 147.5 62.0606 147.5C45.4792 147.5 30.5916 143.994 19.9406 138.445C9.16774 132.834 3.3606 125.52 3.3606 118.12V45.4019Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
<path d="M3.3606 31.8795C3.3606 24.4796 9.16774 17.1662 19.9406 11.5545C30.5916 6.00635 45.4792 2.5 62.0606 2.5C78.642 2.5 93.5296 6.00635 104.181 11.5545C114.953 17.1662 120.761 24.4796 120.761 31.8795C120.761 39.2794 114.953 46.5928 104.181 52.2045C93.5296 57.7527 78.642 61.259 62.0606 61.259C45.4792 61.259 30.5916 57.7527 19.9406 52.2045C9.16774 46.5928 3.3606 39.2794 3.3606 31.8795Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Volume</figcaption></figure>

</div>


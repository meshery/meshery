---
layout: default
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
Currently, the circle is used as the default shape for new components. However, if users or contributors have alternative shapes they believe better suit a particular component, they are encouraged to propose them.

Although the usage of the components is divided into categories, some shapes serve as a universal representation of particular components.

Below are all the shapes with their current usage in a general context.

<div class="svg-grid" style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin: 40px;">

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="68" viewBox="0 0 150 136" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M98.0102 15.75C87.8104 -1.91667 62.3107 -1.91668 52.1109 15.75L7.07755 93.75C-3.1223 111.417 9.62751 133.5 30.0272 133.5H120.094C140.494 133.5 153.243 111.417 143.044 93.75L98.0102 15.75Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Triangle</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 151 150" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="2.56055" y="2.5" width="145" height="145" rx="22.5" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Square (rounded edges)</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="75" viewBox="0 0 151 150" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M43.763 140.376C7.65635 123.089 -7.60209 79.8012 9.68328 43.6899C17.9841 26.3486 32.8327 13.0155 50.9628 6.62363C69.0928 0.231726 89.0191 1.30463 106.358 9.6062C142.465 26.8934 157.723 70.1814 140.438 106.293C123.152 142.404 79.8698 157.664 43.763 140.376Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Circle</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="56" viewBox="0 0 151 113" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="2.56055" y="2.5" width="145" height="107.5" fill="white" fill-opacity="0.3" stroke="#666666" stroke-width="5" stroke-dasharray="10 10"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Rectangle (dashed lines)</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="75" viewBox="0 0 114 151" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="3.31055" y="3" width="107.5" height="145" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Rectangle</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="75" viewBox="0 0 151 151" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="2.56055" y="3" width="145" height="145" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Square (sharp edges) </figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="72" viewBox="0 0 160 155" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M95.6369 7.72668C86.349 0.978614 73.7722 0.978623 64.4843 7.72669L13.8383 44.5231C4.55039 51.2712 0.663959 63.2325 4.21163 74.1511L23.5567 133.689C27.1043 144.608 37.2792 152 48.7597 152H111.362C122.842 152 133.017 144.608 136.565 133.689L155.91 74.151C159.457 63.2324 155.571 51.2712 146.283 44.5231L95.6369 7.72668Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Pentagon</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="61" viewBox="0 0 153 123" fill="white" fill-opacity="0.5" xmlns="http://www.w3.org/2000/svg">
<path d="M144.25 46.2377L112.401 10.7933C107.944 5.8331 101.59 3 94.9215 3H26.0605C13.0819 3 2.56055 13.5213 2.56055 26.5V97C2.56055 109.979 13.0818 120.5 26.0605 120.5H94.9791C101.61 120.5 107.932 117.699 112.387 112.787L144.178 77.731C152.273 68.8045 152.304 55.2013 144.25 46.2377Z" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">ConfigMap</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="78" viewBox="0 0 132 156" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M128.766 32.7442C128.751 30.2268 127.557 28.3686 125.291 27.3159C123.649 26.5561 121.966 25.8909 120.291 25.18L98.4409 16.026C88.553 11.8885 78.6681 7.74784 68.7863 3.60416C67 2.85049 65.271 2.76811 63.4575 3.54924C59.0899 5.42579 54.692 7.23522 50.3092 9.07211L29.7103 17.6982C22.2546 20.8227 14.808 23.9452 7.37042 27.0657C4.69543 28.1855 3.35642 30.1505 3.3534 33.0676C3.3534 39.8812 3.37758 46.6947 3.3534 53.5113C3.30806 62.0031 3.31108 70.4888 4.86469 78.8891C5.64182 83.2089 6.69179 87.4742 8.00818 91.6588C12.3122 105.21 19.4835 117.655 29.0252 128.131C38.5668 138.607 50.2507 146.863 63.2671 152.328C64.8781 153.011 66.5042 153.292 68.1757 152.6C76.6233 149.142 84.5435 144.494 91.7005 138.795C99.4826 132.603 106.283 125.248 111.867 116.985C119.741 105.359 125.028 92.1506 127.364 78.2666C128.116 74.0118 128.56 69.7071 128.691 65.387C128.806 60.1785 128.718 54.9638 128.718 49.7521H128.766C128.766 44.0828 128.797 38.4135 128.766 32.7442Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5" stroke-miterlimit="10"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Shield</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;"><svg width="100%" height="75" viewBox="0 0 124 150" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M3.3606 45.4019C6.77162 49.7691 11.733 53.5667 17.6307 56.6388C29.1298 62.6288 44.8421 66.2589 62.0606 66.2589C79.279 66.2589 94.9914 62.6288 106.49 56.6388C112.388 53.5667 117.35 49.7691 120.761 45.4019V118.12C120.761 125.52 114.953 132.834 104.181 138.445C93.5296 143.994 78.642 147.5 62.0606 147.5C45.4792 147.5 30.5916 143.994 19.9406 138.445C9.16774 132.834 3.3606 125.52 3.3606 118.12V45.4019Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
<path d="M3.3606 31.8795C3.3606 24.4796 9.16774 17.1662 19.9406 11.5545C30.5916 6.00635 45.4792 2.5 62.0606 2.5C78.642 2.5 93.5296 6.00635 104.181 11.5545C114.953 17.1662 120.761 24.4796 120.761 31.8795C120.761 39.2794 114.953 46.5928 104.181 52.2045C93.5296 57.7527 78.642 61.259 62.0606 61.259C45.4792 61.259 30.5916 57.7527 19.9406 52.2045C9.16774 46.5928 3.3606 39.2794 3.3606 31.8795Z" fill="white" fill-opacity="0.5" stroke="#666666" stroke-width="5"/>
</svg><figcaption style="margin-top: 10px; font-size: 14px;">Cylinder</figcaption></figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M38.7212 6.1944V33.8058C38.7212 38.9463 30.5523 38.7427 20.5002 38.7427C10.4482 38.7427 2.2793 38.9463 2.2793 33.8058V6.1944C2.2793 1.05386 10.4227 1.25744 20.5002 1.25744C30.5778 1.25744 38.7212 1.05386 38.7212 6.1944Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
  </svg>
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg" style="display: none;">
    <path d="M38.7212 6.1944V33.8058C38.7212 38.9463 30.5523 38.7427 20.5002 38.7427C10.4482 38.7427 2.2793 38.9463 2.2793 33.8058V6.1944C2.2793 1.05386 10.4227 1.25744 20.5002 1.25744C30.5778 1.25744 38.7212 1.05386 38.7212 6.1944Z" fill="white" fill-opacity="0.1" stroke="white" stroke-miterlimit="10"/>
  </svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Barrel</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M1.5 1H39.5V32.0236C39.5 35.8817 36.3817 39 32.5236 39H8.47635C4.61822 39 1.5 35.8817 1.5 32.0236V1Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
  </svg>
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg" style="display: none;">
    <path d="M1.5 1H39.5V32.0236C39.5 35.8817 36.3817 39 32.5236 39H8.47635C4.61822 39 1.5 35.8817 1.5 32.0236V1Z" fill="white" fill-opacity="0.1" stroke="white" stroke-width="1" stroke-miterlimit="10"/>
  </svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">BottomRoundRectangle</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M20.4716 2.13574L27.1532 12.9296L38.5094 15.8442L30.5529 25.4771L31.6154 38L20.4716 34.0727L9.35152 38L10.1306 25.5759L2.45752 15.8442L14.0026 13.0531L20.4716 2.13574Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
  </svg>
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg" style="display: none;">
    <path d="M20.4716 2.13574L27.1532 12.9296L38.5094 15.8442L30.5529 25.4771L31.6154 38L20.4716 34.0727L9.35152 38L10.1306 25.5759L2.45752 15.8442L14.0026 13.0531L20.4716 2.13574Z" fill="white" fill-opacity="0.1" stroke="white" stroke-width="1" stroke-miterlimit="10"/>
  </svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Star</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M37.4998 3L20.5118 36.9998L3.5 3L20.5118 14.3175L37.4998 3Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
  </svg>
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg" style="display: none;">
    <path d="M37.4998 3L20.5118 36.9998L3.5 3L20.5118 14.3175L37.4998 3Z" fill="white" fill-opacity="0.1" stroke="white" stroke-width="1" stroke-miterlimit="10"/>
  </svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Vee</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
  <svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    <rect x="20.3101" y="0.331543" width="27.8166" height="27.8166" rx="4.45066" transform="rotate(44.4462 20.3101 0.331543)" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
  </svg>
  <svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg" style="display: none;">
    <rect x="20.3101" y="0.331543" width="27.8166" height="27.8166" rx="4.45066" transform="rotate(44.4462 20.3101 0.331543)" fill="white" fill-opacity="0.1" stroke="white" stroke-miterlimit="10"/>
  </svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">RoundDiamond</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M13.9786 0.999905H39.6475L26.9956 39H1.35279L13.9786 0.999905Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
  </svg>
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg" style="display: none;">
    <path d="M13.9786 0.999905H39.6475L26.9956 39H1.35279L13.9786 0.999905Z" fill="white" fill-opacity="0.1" stroke="white" stroke-miterlimit="10"/>
  </svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Right Rhomboid</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M20.5126 2L30.3177 20.0127L38.932 38.0001H20.5126H2.06787L10.7075 20.0127L20.5126 2Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Triangle</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M17.4963 7.00738C18.9399 4.12018 23.0601 4.12019 24.5037 7.00739L37.1655 32.3309C38.4677 34.9355 36.5738 38 33.6618 38H8.33824C5.42623 38 3.53226 34.9355 4.83455 32.3309L17.4963 7.00738Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10" />
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Round Triangle</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M20.5003 1.00017L1.50049 20L20.5003 38.9998L39.5001 20L20.5003 1.00017Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Diamond</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M31.8762 2H9.09886L2.22314 8.6296V31.3704L8.59881 38H32.3763L38.7769 31.1452V8.85475L31.8762 2Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Cut Rectangle</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M27.1192 2H13.8808C12.2577 2 10.7867 2.90503 10.0512 4.33799L2.97552 18.0642C2.34149 19.2709 2.34149 20.729 2.97552 21.9358L10.0512 35.662C10.7867 37.095 12.2577 38 13.8808 38H27.1192C28.7423 38 30.2133 37.095 30.9488 35.662L38.0245 21.9358C38.6585 20.729 38.6585 19.2709 38.0245 18.0642L30.9488 4.33799C30.2133 2.90503 28.7423 2 27.1192 2Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Round Hexagon</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M30 1H11L1.5 20L11 39H30L39.5 20L30 1Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Hexagon</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M18.7217 2.39519L7.1149 7.9655C6.03826 8.46733 5.27673 9.42081 5.01413 10.5499L2.09931 23.2713C1.83671 24.3753 2.09931 25.5044 2.83458 26.3826L10.9488 36.5196C11.7103 37.448 12.8658 38 14.1 38H26.9147C28.1489 38 29.3043 37.448 30.0658 36.5196L38.1801 26.3826C38.8891 25.5044 39.1517 24.3502 38.9153 23.2713L36.0005 10.5248C35.7379 9.39571 34.9764 8.46733 33.8997 7.94041L22.2667 2.39519C21.1375 1.86827 19.8246 1.86827 18.6954 2.39519H18.7217Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Round Heptagon</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M20.5 1L5.26817 8.53202L1.5 25.4529L12.0349 39.0001H28.9651L39.5001 25.4529L35.7319 8.53202L20.5 1Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Heptagon</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M26.6242 2H14.3248C13.3572 2 12.4405 2.3719 11.7529 3.04132L3.06951 11.4959C2.38197 12.1653 2 13.0579 2 14V25.9752C2 26.9174 2.38197 27.8099 3.06951 28.4793L11.7784 36.9587C12.4659 37.6281 13.3827 38 14.3503 38H26.6497C27.6173 38 28.5341 37.6281 29.2216 36.9587L37.9305 28.4793C38.618 27.8099 39 26.9174 39 25.9752V14C39 13.0579 38.618 12.1653 37.9305 11.4959L29.1961 3.04132C28.5086 2.3719 27.5919 2 26.6242 2Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">RoundOctagon</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="1.5" y="1" width="38" height="38" rx="5" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">RoundRectangle</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg xmlns="http://www.w3.org/2000/svg" id="Layer_1" data-name="Layer 1" viewBox="0 0 134.95 135.02" width="100%" height="75">
<polygon points="69.49 31.82 69.49 64.07 97.44 47.89 69.49 31.82" fill="#00d3a9"/>
<polygon points="69.49 70.81 69.49 103.22 97.7 87.09 69.49 70.81" fill="#00d3a9"/>
<polygon points="65.47 63.85 65.47 32.09 37.87 47.92 65.47 63.85" fill="#00b39f"/>
<path d="M10.1,103.1a67.79,67.79,0,0,0,21.41,21.55V90.71Z" fill="#00b39f"/>
<polygon points="65.47 103.06 65.47 71.05 37.8 87.07 65.47 103.06" fill="#00b39f"/>
<polygon points="35.54 122.63 63.56 106.61 35.54 90.41 35.54 122.63" fill="#00d3a9"/>
<polygon points="99.61 122.8 99.61 90.63 71.63 106.63 99.61 122.8" fill="#00b39f"/>
<path d="M127,99.37a67.22,67.22,0,0,0,7.91-28.94L105.78,87.11Z" fill="#00b39f"/>
<polygon points="103.64 83.69 131.76 67.61 103.64 51.45 103.64 83.69" fill="#00d3a9"/>
<polygon points="99.61 44.5 99.61 12.52 71.76 28.49 99.61 44.5" fill="#00b39f"/>
<polygon points="99.61 83.55 99.61 51.28 71.7 67.44 99.61 83.55" fill="#00b39f"/>
<polygon points="67.48 135.02 67.49 135.02 67.48 135.02 67.48 135.02" fill="#00b39f"/>
<polygon points="35.54 51.22 35.54 83.73 63.66 67.45 35.54 51.22" fill="#00d3a9"/>
<path d="M65.47,0A67.2,67.2,0,0,0,35.83,7.83l29.64,17Z" fill="#00b39f"/>
<polygon points="35.54 12.3 35.54 44.62 63.68 28.48 35.54 12.3" fill="#00d3a9"/>
<path d="M31.51,10.34A67.89,67.89,0,0,0,10.1,31.89L31.51,44.25Z" fill="#00b39f"/>
<path d="M99.43,8A67.23,67.23,0,0,0,69.49,0V25.15Z" fill="#00d3a9"/>
<path d="M0,69.87A67.27,67.27,0,0,0,8.07,99.63L29.76,87.07Z" fill="#00d3a9"/>
<path d="M8.07,35.37A67.16,67.16,0,0,0,0,65L29.79,47.91Z" fill="#00d3a9"/>
<path d="M35.78,127.13A67.13,67.13,0,0,0,65.47,135V110.15Z" fill="#00b39f"/>
<path d="M124.92,32a67.9,67.9,0,0,0-21.28-21.52V44.3Z" fill="#00d3a9"/>
<path d="M103.64,124.54A68,68,0,0,0,125,102.86L103.64,90.52Z" fill="#00d3a9"/>
<path d="M135,64.81a67.06,67.06,0,0,0-8-29.35L105.49,47.88Z" fill="#00b39f"/>
<path d="M69.49,135a67.12,67.12,0,0,0,29.63-7.83L69.49,110Z" fill="#00d3a9"/>
<polygon points="31.51 83.44 31.51 51.56 3.83 67.43 31.51 83.44" fill="#00b39f"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Polygon</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M20.4711 1.78955L2.37842 15.614L9.28517 37.9999H31.657L38.5888 15.614L20.4711 1.78955Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Pentagon</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M27.2792 13.5545V1H13.9569V13.5545H1.5V26.5232H13.9569V39H27.2792V26.5232H39.5V13.5545H27.2792Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10" stroke-width="1"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Plus</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M34.4009 20L38.5575 37H2.44189L6.62344 20L2.44189 3H38.5575L34.4009 20Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">ConcaveHexagon</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M39.5 20.0876L26.6367 1H1.5V39H26.6367L39.5 20.0876Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Tag</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M37.3812 16.9435L28.3015 3.47042C27.2607 1.92594 25.5201 1 23.6576 1H7.1C4.00721 1 1.5 3.50721 1.5 6.6V33.4C1.5 36.4928 4.00721 39 7.1 39H23.673C25.527 39 27.2608 38.0825 28.3034 36.5495L37.3678 23.2225C38.6556 21.3291 38.6609 18.8424 37.3812 16.9435Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Round Tag</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 41" fill="none" xmlns="http://www.w3.org/2000/svg">
<polygon fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10" points="39.76,40.5 0.5,40.5 0.5,1.24 19.64,20.38"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Left Triangle</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M12.5136 36.6823C3.30016 32.2711 -0.593367 21.2252 3.81738 12.0105C5.93552 7.58552 9.72449 4.18328 14.3508 2.55224C18.9771 0.921198 24.0617 1.19497 28.4862 3.31331C37.6997 7.72454 41.5931 18.7705 37.1824 27.9851C32.7717 37.1997 21.7271 41.0935 12.5136 36.6823Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Circle</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path fill-rule="evenodd" clip-rule="evenodd" d="M8.58192 7.90585C10.1208 5.93621 11.9673 4.40122 14.0188 3.34985C12.1971 9.82411 14.2885 17.0898 19.916 21.4866C25.423 25.789 32.7012 26.1674 38.5375 23.0885C37.9156 25.1274 36.927 27.1385 35.5137 28.9473C29.7037 36.3837 18.9635 37.7024 11.5272 31.8924C4.09071 26.0825 2.77197 15.3423 8.58192 7.90585Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Crescent</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M1.77734 39L39.2232 29.7375V1L1.77734 10.5V39Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Parallelogram</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M39 1H1V39H39V1Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Rectangle</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M13.9785 38.9999H39.6475L26.9956 0.999756H1.35274L13.9785 38.9999Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Rhomboid</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<polygon points="30.33,39.5 0.5,39.5 10.53,10.73 39.5,0.5" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">SlantedParallelogram</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
  <svg width="100%" height="75" viewBox="0 0 41 40" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path stroke-miterlimit="10" d="M30.8738 19.6872C30.6794 19.9805 30.6794 19.9805 30.8738 20.0782C32.0402 21.251 33.2065 22.3261 34.3729 23.4989C36.0252 25.1604 37.6776 26.7242 39.3299 28.3857C39.3557 28.3857 39.3804 28.396 39.3986 28.4143C39.4169 28.4326 39.4271 28.4575 39.4271 28.4834C39.3299 28.5811 39.1355 28.6789 39.0383 28.7766L34.9561 31.7087C33.9841 32.3928 32.915 33.1747 31.943 33.9566C31.8458 34.0543 31.8458 33.9566 31.7486 33.9566C30.7766 32.7837 29.8047 31.7087 28.8327 30.5358C27.472 28.9721 26.1112 27.4083 24.7505 25.7468C23.6813 24.4763 22.6121 23.3034 21.6402 22.0329L19.9879 20.0782H19.9393L14.5935 26.3332L9.24766 32.6274C8.85888 33.0183 8.56729 33.507 8.1785 33.8979C8.08131 33.9957 8.08131 33.9957 7.98411 33.8979C6.33178 32.7251 4.77664 31.5523 3.1243 30.2817C2.34673 29.793 1.56916 29.2066 0.694393 28.6202C0.694393 28.5225 0.694393 28.5225 0.597196 28.5225L0.888785 28.2293C2.54112 26.6655 4.09626 25.1018 5.7486 23.538C6.91495 22.4629 8.08131 21.2901 9.15047 20.215C9.34486 20.0195 9.34486 20.0195 9.15047 19.9218C7.98411 18.749 6.81776 17.6739 5.6514 16.5011C3.99907 14.8005 2.24953 13.2758 0.597196 11.6143L0.5 11.5166C0.597196 11.4189 0.694393 11.3211 0.791589 11.3211C2.73551 9.8551 4.67944 8.48681 6.62336 7.02079C7.10935 6.72758 7.49813 6.33664 7.98411 6.04344C8.08131 5.9457 8.08131 6.04344 8.1785 6.04344L11.5804 10.0506L16.343 15.7192L19.9393 19.9218L19.9976 19.9805L20.085 19.8827C21.4458 18.319 22.8065 16.6575 24.1673 15.0937L27.7636 10.8911L31.1654 6.88396C31.3676 6.60266 31.5954 6.3409 31.8458 6.10208C31.943 6.00434 31.943 6.00434 32.0402 6.10208C33.4981 7.17716 34.9561 8.25225 36.5112 9.32733C37.4832 10.0115 38.4551 10.6956 39.3299 11.3798H39.4271C39.5243 11.3798 39.5243 11.4775 39.4271 11.4775C39.0383 11.7707 38.7467 12.1617 38.3579 12.4549C36.0252 14.7028 33.6925 16.9507 31.457 19.1986C31.2877 19.3895 31.0912 19.5542 30.8738 19.6872Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F"/>
  </svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">XWing</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M8.5 3H32.5L39.5 37H1.5L8.5 3Z" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-miterlimit="10" stroke-width="1"/>
</svg>
  <figcaption style="margin-top: 10px; font-size: 14px;">Trapezoid</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<circle cx="20" cy="20" r="19" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="2"/>
<text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" font-family="Arial, sans-serif" font-size="10" fill="#00B39F">SQL</text>
</svg>
<figcaption style="margin-top: 10px; font-size: 14px;">SqlJob</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<circle cx="20" cy="20" r="19" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="2"/>
<path d="M20 10C22.2091 10 24 11.7909 24 14C24 16.2091 22.2091 18 20 18C17.7909 18 16 16.2091 16 14C16 11.7909 17.7909 10 20 10Z" fill="#00B39F"/>
<path d="M26 28C26 24.6863 23.3137 22 20 22C16.6863 22 14 24.6863 14 28" stroke="#00B39F" stroke-width="2"/>
</svg>
<figcaption style="margin-top: 10px; font-size: 14px;">User</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="2" y="10" width="36" height="20" rx="2" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="2" stroke-dasharray="4 4"/>
<line x1="6" y1="15" x2="34" y2="15" stroke="#00B39F" stroke-width="2"/>
<line x1="6" y1="20" x2="28" y2="20" stroke="#00B39F" stroke-width="2"/>
<line x1="6" y1="25" x2="22" y2="25" stroke="#00B39F" stroke-width="2"/>
</svg>
<figcaption style="margin-top: 10px; font-size: 14px;">TextBox</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<rect x="1.5" y="1" width="38" height="38" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1" stroke-miterlimit="10" stroke-dasharray="4 4"/>
</svg>
<figcaption style="margin-top: 10px; font-size: 14px;">Section</figcaption>
</figure>

<figure style="display: flex; flex-direction: column; align-items: center;">
<svg width="100%" height="75" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
<circle cx="20" cy="20" r="19" fill="#00B39F" fill-opacity="0.1" stroke="#00B39F" stroke-width="1"/>
</svg>
<figcaption style="margin-top: 10px; font-size: 14px;">GenericNode</figcaption>
</figure>

</div>
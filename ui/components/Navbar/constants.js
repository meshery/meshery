import React from "react";

// import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import DashboardIcon from "@mui/icons-material/Dashboard";
import LifecycleIcon from "./drawer-icons/lifecycle_mgmt_svg";
import PerformanceIcon from "./drawer-icons/performance_svg";
import GitHubIcon from "@mui/icons-material/GitHub";
import DescriptionOutlinedIcon from "@mui/icons-material/DescriptionOutlined";
import MailIcon from "@mui/icons-material/Mail";
import ConformanceIcon from "./drawer-icons/conformance_svg";
import SmiIcon from "./drawer-icons/servicemeshinterface-icon-white_svg";
import { drawerIconsStyle } from "./Navbar.styles";
import { faSlack } from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
// import { faDigitalTachograph } from "@fortawesome/free-solid-svg-icons";
// import { faSlack } from "@fortawesome/free-brands-svg-icons";
// import { faExternalLinkAlt } from "@fortawesome/free-solid-svg-icons";

// const ExternalLinkIcon = (
//   <FontAwesomeIcon style={externalLinkIconStyle} icon={MailIcon} transform="shrink-7" />
// );

const ExternalLinkIcon = <MailIcon />;

export const externlinks = [
  {
    id: "doc",
    href: "https://docs.meshery.io",
    title: "Documentation",
    icon: <DescriptionOutlinedIcon style={drawerIconsStyle} />,
    external_icon: ExternalLinkIcon,
  },
  {
    id: "community",
    href: "http://slack.layer5.io",
    title: "Community",
    icon: <FontAwesomeIcon style={{ marginBottom: 2, ...drawerIconsStyle }} icon={faSlack} transform="grow-1" />,
    external_icon: ExternalLinkIcon,
  },
  {
    id: "mailinglist",
    href: "https://meshery.io/subscribe",
    title: "Mailing List",
    icon: <MailIcon style={drawerIconsStyle} />,
    external_icon: ExternalLinkIcon,
  },
  {
    id: "issues",
    href: "https://github.com/meshery/meshery/issues/new/choose",
    title: "Issues",
    icon: <GitHubIcon style={drawerIconsStyle} />,
    external_icon: ExternalLinkIcon,
  },
];

export const navigatorItemsTree = [
  {
    id: "Dashboard",
    icon: <DashboardIcon style={drawerIconsStyle} />,
    href: "/",
    title: "Dashboard",
    show: true,
    link: true,
  },
  {
    id: "Lifecycle",
    icon: <LifecycleIcon style={drawerIconsStyle} />,
    href: "/management",
    title: "Lifecycle",
    show: true,
    link: true,
    children: [
      {
        id: "App_Mesh",
        href: "/management/app-mesh",
        title: "AWS App Mesh",
        link: true,
        show: true,
      },
      {
        id: "Citrix_Service_Mesh",
        href: "/management/citrix",
        title: "Citrix Service Mesh",
        link: true,
        show: true,
      },
      {
        id: "Consul",
        href: "/management/consul",
        title: "Consul",
        link: true,
        show: true,
      },
      {
        id: "Istio",
        href: "/management/istio",
        title: "Istio",
        link: true,
        show: true,
      },
      {
        id: "Kuma",
        href: "/management/kuma",
        title: "Kuma",
        link: true,
        show: true,
      },
      {
        id: "Linkerd",
        href: "/management/linkerd",
        title: "Linkerd",
        link: true,
        show: true,
      },
      {
        id: "Network_Service_Mesh",
        href: "/management/nsm",
        title: "Network Service Mesh",
        link: true,
        show: true,
      },
      {
        id: "NGINX_Service_Mesh",
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />,
        href: "/management/nginx",
        title: "NGINX Service Mesh",
        link: true,
        show: true,
      },
      {
        id: "Octarine",
        href: "/management/octarine",
        title: "Octarine",
        link: true,
        show: true,
      },
      {
        id: "Open_Service_Mesh",
        href: "/management/osm",
        title: "Open Service Mesh",
        link: true,
        show: true,
      },
      {
        id: "Traefik_Mesh",
        href: "/management/traefik-mesh",
        title: "Traefik Mesh",
        link: true,
        show: true,
      },
    ],
  },
  {
    id: "Configuration",
    icon: <img src="/static/img/configuration_trans.svg" style={{ width: "1.21rem" }} />,
    // disabled: true,
    href: "#",
    title: "Configuration",
    show: true,
    link: true,
    children: [
      {
        id: "Applications",
        icon: <img src="/static/img/web-applications.svg" style={{ width: "1.21rem" }} />,
        href: "/configuration/applications",
        title: "Applications",
        show: true,
        link: true,
      },
      {
        id: "Filters",
        icon: <img src="/static/img/web-filters.svg" style={{ width: "1.21rem" }} />,
        href: "/configuration/filters",
        title: "Filters",
        show: true,
        link: true,
      },
      {
        id: "Patterns",
        icon: <img src="/static/img/pattern_trans.svg" style={{ width: "1.21rem" }} />,
        href: "/configuration/patterns",
        title: "Patterns",
        show: false,
        link: true,
      },
    ],
  },
  {
    id: "Performance",
    icon: <PerformanceIcon style={{ transform: "scale(1.3)", ...drawerIconsStyle }} />,
    href: "/performance",
    title: "Performance",
    show: true,
    link: true,
    children: [
      {
        id: "Profiles",
        // icon: <FontAwesomeIcon icon={MailIcon} transform="shrink-2" style={{ verticalAlign: "top" }} />,
        icon: <MailIcon />,
        href: "/performance/profiles",
        title: "Profiles",
        show: true,
        link: true,
      },
    ],
  },
  {
    id: "Settings",
    href: "/settings",
    title: "Settings",
    show: false,
    link: true,
  }, // title is used for comparison in the Header.js file as well
  {
    id: "Conformance",
    icon: <ConformanceIcon style={drawerIconsStyle} />,
    href: "/smi_results", //Temp
    title: "Conformance",
    show: true,
    link: true,
    children: [
      {
        id: "Service Mesh Interface",
        icon: <SmiIcon style={drawerIconsStyle} />,
        href: "/smi_results",
        title: "Service Mesh Interface",
        show: true,
        link: true,
      },
    ],
  },
];

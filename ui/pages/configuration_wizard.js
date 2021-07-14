import ServiceSwitch from "../components/ConfigurationWizard/ServiceSwitchCard.js"
import KubernetesIcon from "../components/ConfigurationWizard/icons/KubernetesIcon.js"
import GrafanaIcon from "../components/ConfigurationWizard/icons/GrafanaIcon.js"
import PrometheusIcon from "../components/ConfigurationWizard/icons/PrometheusIcon.js"
import ServiceSpecificConfig  from "../components/ConfigurationWizard/ServiceSpecificConfig.js"
import KubernetesConfigComp  from "../components/ConfigurationWizard/ServiceSpecifComponents/Kubernetes.js"

const ConfigurationWizard = () => {

  const prometheusserviceInfo = {
    name: "Prometheus",
    logoComponent: PrometheusIcon,
    configComp : () => null
  }
  const kubeserviceInfo = {
    name: "Kubernetes",
    logoComponent: KubernetesIcon,
    configComp :  <ServiceSpecificConfig components={[KubernetesConfigComp]} />

  }
  const grafanaserviceInfo = {
    name: "Kubernetes",
    logoComponent: GrafanaIcon,
    configComp : () => null
  }

  return (
    <>
      <ServiceSwitch serviceInfo={kubeserviceInfo} /> 
    </>
  )
}

export default ConfigurationWizard

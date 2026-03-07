<style>
  td:hover,tr:hover {
      background-color: var(--color-primary-dark);
      cursor:pointer;
    }
    td.details {
      background-color: #fafafa;
      cursor:text;
    }
    .yellowCheckbox{
      width:2.5rem
    }
    .tooltipss{
      position:relative;
      width:fit-content;
      cursor:pointer;
      margin-right: auto;
      margin-left: auto;
    }
    .tooltipss .tooltiptext {
    visibility: hidden;
    width: 120px;
    background-color: #555;
    color: #fff;
    text-align: center;
    border-radius: 6px;
    padding: 5px 0;
    position: absolute;
    z-index: 1;
    bottom: 125%;
    left: 50%;
    margin-left: -60px;
    opacity: 0;
    transition: opacity 0.3s;
    }

  .tooltipss .tooltiptext::after {
    content: "";
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -5px;
    border-width: 5px;
    border-style: solid;
    border-color: #555 transparent transparent transparent;
  }
  .tooltipss:hover .tooltiptext {
    visibility: visible;
    opacity: 1;
}


</style>

<table class="table table-striped" >
  <th>Kubernetes Version</th>

  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="/installation/compatibility-matrix/images/service-meshes/istio.svg" /><a href="{{ .Site.Params.repo }}-istio">meshery-istio</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="/installation/compatibility-matrix/images/service-meshes/linkerd.svg" /><a href="{{ .Site.Params.repo }}-linkerd">meshery-linkerd</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="/installation/compatibility-matrix/images/service-meshes/kuma.svg" /><a href="{{ .Site.Params.repo }}-kuma">meshery-kuma</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="/installation/compatibility-matrix/images/adapters/nighthawk/nighthawk.svg" /><a href="{{ .Site.Params.repo }}-nighthawk">meshery-nighthawk</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="/installation/compatibility-matrix/images/service-meshes/nginx-sm.svg" /><a href="{{ .Site.Params.repo }}-nginx-sm">meshery-nginx-sm</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="/installation/compatibility-matrix/images/service-meshes/traefik-mesh.svg" /><a href="{{ .Site.Params.repo }}-traefik-mesh">meshery-traefik-mesh</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="/installation/compatibility-matrix/images/service-meshes/cilium.svg" /><a href="{{ .Site.Params.repo }}-cilium">meshery-cilium</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="/installation/compatibility-matrix/images/service-meshes/consul.svg" /><a href="{{ .Site.Params.repo }}-consul">meshery-consul</a></th>

{{/* 
Jekyll include data loop - DEPRECATED
Functionality now in compatibility-matrix-kubernetes.html shortcode
This file is no longer being used. All Jekyll liquid syntax is commented out.

Original Jekyll loop:
{% for k8s in {{include.k8s_tests_group}} %}
*/}}

  {{/* DEPRECATED: All Jekyll liquid loop code below is commented out
  Functionality has been migrated to compatibility-matrix-kubernetes.html shortcode
  Original Jekyll loop from here onwards is non-functional in Hugo.
  */}}

</table>

<script>
  function showCompatability () {
      let percentContainer = document.querySelectorAll(".compatibility")
      for(let i = 0 ; i<percentContainer.length;i++){
        let percentage = parseFloat(percentContainer[i].innerHTML);
        if (percentage >= 90.00){
          percentContainer[i].innerHTML = `
            <div class = "tooltipss" style="text-align:center">
              <img src = "/assets/img/passing.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>
          `
        }
        else if(percentage >=1 && percentage<=89.99){
          percentContainer[i].innerHTML = `<div class = "tooltipss" style="text-align:center">
              <img src = "/assets/img/YellowCheck.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>`
        }
        else if(percentage < 0){
           percentContainer[i].innerHTML = `<div class = "tooltipss" style="text-align:center">
              <img src = "/assets/img/na-icon.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">Not Applicable</span>
            </div>`
        }
        else{
           percentContainer[i].innerHTML = `<div class = "tooltipss" style="text-align:center">
              <img src = "/assets/img/failing.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>`
        }
      }
    }
  function clickIcon(serviceMesh){
    location.href = `/installation/compatibility-matrix/${serviceMesh}-past-results`
  }

showCompatability()
</script>


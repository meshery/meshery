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

  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/istio.svg" /><a href="{{ site.repo }}-istio">meshery-istio</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/linkerd.svg" /><a href="{{ site.repo }}-linkerd">meshery-linkerd</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/kuma.svg" /><a href="{{ site.repo }}-kuma">meshery-kuma</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/osm.svg" /><a href="{{ site.repo }}-osm">meshery-osm</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/nginx-sm.svg" /><a href="{{ site.repo }}-nginx-sm">meshery-nginx-sm</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/traefik-mesh.svg" /><a href="{{ site.repo }}-traefik-mesh">meshery-traefik-mesh</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/cilium.svg" /><a href="{{ site.repo }}-cilium">meshery-cilium</a></th>
  <th><img style="height: 1.5rem; vertical-align: text-bottom;" src="{{site.baseurl}}/assets/img/service-meshes/consul.svg" /><a href="{{ site.repo }}-consul">meshery-consul</a></th>

{% for k8s in {{include.k8s_tests_group}} %}

  <tr class = "first-row">
    {% assign successfull_istio = 0 %}
    {% assign successfull_linkerd = 0 %}
    {% assign successfull_cilium = 0 %}
    {% assign successfull_osm = 0 %}
    {% assign successfull_kuma = 0 %}
    {% assign successfull_traefik_mesh = 0 %}
    {% assign successfull_nginx_sm = 0 %}
    {% assign successfull_consul = 0 %}

    {%if k8s.name != "v1.21.0"%}
      <td>{{k8s.name}}</td>
      {% assign k8s_items = k8s.items | group_by: "meshery-component"  %}
      {% for k8s_item in k8s_items %}
        {% if k8s_item.name == "meshery-linkerd" %}
          {% assign linkerd_size = k8s_item.size | times:1.0 %}
          {% assign linkerd_item = k8s_item.items | sort: "meshery-component-version" | reverse %}
          {% for single in linkerd_item limit: 1 %}
            {% if single.overall-status == "passing" %}
              {% assign successfull_linkerd = 1 %}
            {% elsif single.overall-status == "failing" %}
              {% assign successfull_linkerd = 0 %}
            {% else %}
              {% assign successfull_linkerd = 0.5%}
            {% endif %}
          {% endfor %}

        {% elsif k8s_item.name == "meshery-istio" %}
          {% assign istio_size = k8s_item.size | times:1.0 | times:1.0 %}
          {% assign istio_items = k8s_item.items | sort: "meshery-component-version" | reverse %}
          {% for single in istio_items limit: 1 %}
             {% if single.overall-status == "passing" %}
              {% assign successfull_istio = 1 %}
            {% elsif single.overall-status == "failing" %}
              {% assign successfull_istio = 0 %}
            {% else %}
              {% assign successfull_istio = 0.5%}
            {% endif %}
          {% endfor %}

        {% elsif k8s_item.name == "meshery-kuma" %}
          {% assign kuma_size = k8s_item.size | times:1.0 %}
          {% assign kuma_items = k8s_item.items | sort: "meshery-component-version" | reverse  %}
          {% for single in kuma_items limit: 1 %}
            {% if single.overall-status == "passing" %}
              {% assign successfull_kuma = 1 %}
            {% elsif single.overall-status == "failing" %}
              {% assign successfull_kuma = 0 %}
            {% else %}
              {% assign successfull_kuma = 0.5%}
            {% endif %}
          {% endfor %}

        {% elsif k8s_item.name == "meshery-osm" %}
          {%if k8s.name == "v1.20.1" or k8s.name == "v1.21.5" or k8s.name == "v1.20.11" or k8s.name == "v1.22.2" %}
            {% assign osm_size = 0 | times:1.0 %}
          {%else%}
            {% assign osm_size = k8s_item.size | times:1.0 %}
            {% assign osm_items = k8s_item.items | sort: "meshery-component-version" | reverse  %}
          {%endif%}
          {% for single in osm_items limit: 1 %}
             {% if single.overall-status == "passing" %}
              {% assign successfull_osm = 1 %}
            {% elsif single.overall-status == "failing" %}
              {% assign successfull_osm = 0 %}
            {% else %}
              {% assign successfull_osm = 0.5%}
            {% endif %}
          {% endfor %}

        {% elsif k8s_item.name == "meshery-cilium" %}
          {% assign cilium_size = k8s_item.size | times:1.0 %}
          {% assign cilium_items = k8s_item.items | sort: "meshery-component-version" | reverse  %}
          {% for single in cilium_items limit: 1 %}
            {% if single.overall-status == "passing" %}
              {% assign successfull_cilium = 1 %}
            {% elsif single.overall-status == "failing" %}
              {% assign successfull_cilium = 0 %}
            {% else %}
              {% assign successfull_cilium = 0.5%}
            {% endif %}
          {% endfor %}

        {% elsif k8s_item.name == "meshery-nginx-sm"%}
          {% assign nginx_size = k8s_item.size | times:1.0 %}
          {% assign nginx_items = k8s_item.items | sort: "meshery-component-version" | reverse  %}
          {% for single in nginx_items limit: 1 %}
            {% if single.overall-status == "passing" %}
              {% assign successfull_nginx_sm = 1 %}
            {% elsif single.overall-status == "failing" %}
              {% assign successfull_nginx_sm = 0 %}
            {% else %}
              {% assign successfull_nginx_sm = 0.5%}
            {% endif %}
          {% endfor %}

        {% elsif k8s_item.name == "meshery-consul"%}
          {% assign consul_size = k8s_item.size | times:1.0 %}
          {% assign consul_items = k8s_item.items | sort: "meshery-component-version" | reverse  %}
          {% for single in consul_items limit: 1 %}
             {% if single.overall-status == "passing" %}
              {% assign successfull_consul = 1 %}
            {% elsif single.overall-status == "failing" %}
              {% assign successfull_consul = 0 %}
            {% else %}
              {% assign successfull_consul = 0.5%}
            {% endif %}
          {% endfor %}

        {% elsif k8s_item.name == "meshery-traefik-mesh" %}
          {% assign traefik_size = k8s_item.size | times:1.0 %}
          {% assign traefik_items = k8s_item.items | sort: "meshery-component-version" | reverse  %}
          {% for single in traefik_items limit: 1 %}
            {% if single.overall-status == "passing" %}
              {% assign successfull_traefik_mesh = 1 %}
            {% elsif single.overall-status == "failing" %}
              {% assign successfull_traefik_mesh = 0 %}
            {% else %}
              {% assign successfull_traefik_mesh = 0.5%}
            {% endif %}
          {% endfor %}
        {% endif %}
      {% endfor %}

      {%if istio_size and istio_size !=0%}
        {% assign istio_percentage = successfull_istio | divided_by: 1 | times:100 | round:2 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-istio`)" class = "compatibility">{{istio_percentage}}%</td>
      {%else%}
        {% assign istio_percentage = -100.0 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-istio`)" class = "compatibility">{{istio_percentage}}%</td>
      {%endif%}

      {%if linkerd_size and linkerd_size !=0%}
        {% assign linkerd_percentage = successfull_linkerd | divided_by: 1 | times:100 | round:2 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-linkerd`)" class = "compatibility">{{linkerd_percentage}}%</td>
      {%else%}
        {% assign linkerd_percentage = -100.0 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-linkerd`)" class = "compatibility">{{linkerd_percentage}}%</td>
      {%endif%}

      {%if kuma_size and kuma_size !=0%}
       {% assign kuma_percentage = successfull_kuma | divided_by: 1 | times:100 | round:2 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-kuma`)" class = "compatibility">{{kuma_percentage}}%</td>
      {%else%}
        {% assign kuma_percentage = -100.0 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-kuma`)" class = "compatibility">{{kuma_percentage}}%</td>
      {%endif%}

      {%if osm_size and osm_size !=0%}
        {% assign osm_percentage = successfull_osm | divided_by: 1 | times:100 | round:2 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-osm`)" class = "compatibility">{{osm_percentage}}%</td>
      {%else%}
        {% assign osm_percentage = -100.0 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-osm`)" class = "compatibility">{{osm_percentage}}%</td>
      {%endif%}

      {%if nginx_size and nginx_size !=0%}
        {% assign nginx_percentage = successfull_nginx_sm | divided_by: 1 | times:100 | round:2 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-nginx-sm`)" class = "compatibility">{{nginx_percentage}}%</td>
      {%else%}
        {% assign nginx_percentage = -100.0 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-nginx-sm`)" class = "compatibility">{{nginx_percentage}}%</td>
      {%endif%}

      {%if traefik_size and traefik_size !=0%}
        {% assign traefik_percentage = successfull_traefik_mesh | divided_by: 1 | times:100 | round:2 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-traefik-mesh`)" class = "compatibility">{{traefik_percentage}}%</td>
      {%else%}
        {% assign traefik_percentage = -100.0 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-traefik-mesh`)" class = "compatibility">{{traefik_percentage}}%</td>
      {%endif%}

      {%if cilium_size and cilium_size !=0%}
        {% assign cilium_percentage = successfull_cilium | divided_by: 1 | times:100 | round:2 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-cilium`)" class = "compatibility">{{cilium_percentage}}%</td>
      {%else%}
        {% assign cilium_percentage = -100.0 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-cilium`)" class = "compatibility">{{cilium_percentage}}%</td>
      {%endif%}

      {%if consul_size and consul_size !=0%}
        {% assign consul_percentage = successfull_consul | divided_by: 1 | times:100 | round:2 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-consul`)" class = "compatibility">{{consul_percentage}}%</td>
      {%else%}
        {% assign consul_percentage = -100.0 %}
        <td style="text-align:center" onclick = "clickIcon(`meshery-consul`)" class = "compatibility">{{consul_percentage}}%</td>
      {%endif%}

      {%assign consul_size = 0 %}
      {%assign cilium_size = 0 %}
      {%assign traefik_size = 0 %}
      {%assign nginx_size = 0 %}
      {%assign linkerd_size = 0 %}
      {%assign istio_size = 0 %}
      {%assign kuma_size = 0 %}
      {%assign osm_size = 0 %}
    {%endif%}

  </tr>
{% endfor %}

</table>

<script>
  function showCompatability () {
      let percentContainer = document.querySelectorAll(".compatibility")
      for(let i = 0 ; i<percentContainer.length;i++){
        let percentage = parseFloat(percentContainer[i].innerHTML);
        if (percentage >= 90.00){
          percentContainer[i].innerHTML = `
            <div class = "tooltipss" style="text-align:center">
              <img src = "{{site.baseurl}}/assets/img/passing.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>
          `
        }
        else if(percentage >=1 && percentage<=89.99){
          percentContainer[i].innerHTML = `<div class = "tooltipss" style="text-align:center">
              <img src = "{{site.baseurl}}/assets/img/YellowCheck.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>`
        }
        else if(percentage < 0){
           percentContainer[i].innerHTML = `<div class = "tooltipss" style="text-align:center">
              <img src = "{{site.baseurl}}/assets/img/na-icon.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">Not Applicable</span>
            </div>`
        }
        else{
           percentContainer[i].innerHTML = `<div class = "tooltipss" style="text-align:center">
              <img src = "{{site.baseurl}}/assets/img/failing.svg" class = "yellowCheckbox" >
              <span class = "tooltiptext">${percentage}%</span>
            </div>`
        }
      }
    }
  function clickIcon(serviceMesh){
    location.href = `{{site.baseurl}}/installation/compatibility-matrix/${serviceMesh}-past-results`
  }

showCompatability()
</script>


---
layout: page
title: Genel açıklama
permalink: tr/overview
language: tr
lang: tr
categories: tr
---

# Meshery'ye Giriş
Hizmet ağı yönetim düzlemi, farklı hizmet ağlarını benimser, işler ve geliştirir.
Meshery, hizmet ağlarının işlevselliği ve performansı hakkında öğrenmeyi kolaylaştırır ve bir hizmet ağı veya birden çok hizmet ağı içinde çalışan uygulamalardan metriklerin toplanmasını ve görselleştirilmesini içerir.
Meshery şu üst düzey işlevleri sağlar:

1. Hizmet ağlarının performans yönetimi.
2. Hizmet ağları yapılandırma yönetimi.
     - Yapılandırma en iyi uygulamaları.
3. Hizmet ağlarının yaşam döngüsü yönetimi.
4. Birlikte çalışabilirlik ve hizmet ağlarının ilişkilendirilmesi.

<iframe class="container" width="560" height="315" src="https://www.youtube.com/embed/CFj1O_uyhhs" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>

<div style="text-align:center;width:100%"><emphasis>Service Mesh Day 2019'da sunulmuştur</emphasis></div>


<h2>Meshery hangi zorlukları çözer?</h2>
<b>Hizmet ağı yönetimi - bir veya birden çok hizmet ağı.</b>

Bu performans sorularının yanıtlanması gerektiğinde, bunlar ölçüm için kullanılan belirli iş yükü ve altyapıya göre özneldir. Bu zorlukla karşı karşıya kalan Elçi projesi, örneğin, performans verilerini yayınlamayı reddediyor çünkü bu tür testler şunlar olabilir:
- Şaşkın
- Yanlış yorumlandı.

Farklı iş yüklerinin (uygulamaların) bir permütasyonu altında performans ve yüksek veri tüketimi ihtiyacının ötesinde, altyapı kaynaklarının türleri ve boyutları, çok işlevli bir proje ihtiyacı ve eşler arasındaki karşılaştırmalar, davranış farklılıklarının karşılaştırılmasını kolaylaştırmak için istenir. hizmet ağları ve kullanımlarının seçimi arasında. Bireysel projeler, diğer hizmet ağı yarışmalarından test sonuçlarının yayınlanması için ayrılmıştır. Bağımsız, tarafsız ve güvenilir bir analiz gereklidir.

Meshery'nin amacı, hizmet ağı performansını tek tip olarak değerlendirmek için satıcı ve projeden bağımsız bir yardımcı program olmaktır. Hizmet ağı projeleri ve proxy hizmetleri arasında (ve şaşırtıcı bir şekilde, tek bir proje içinde), farklı sayıda araç ve sonuç vardır. Meshery, performans ölçümleri ve değerlendirmeleri sağlayarak ekosisteminiz için verimli bir araç seti seçmenizi sağlar.

1. Meshery'den yararlanarak, eş hizmet ağlarının performans karşılaştırmasını elde edebilirsiniz.
2. Bir hizmet ağının performansını sürümden yayına kadar izler.
3. Hizmet ağları arasındaki davranış farklılıklarını anlayın.
4. Uygulamanızın performansını sürümden sürüme izler.

## Meshery Meraklılar ve Operatörler içindir
İster 0. Günden dağıtım veya bakımın 2. Gününe kadar dağıtım seçeneği olsun, Meshery her koşul için yararlı yeteneklere sahiptir. Meshery ile bir proje için hedeflenen kitle, ekosistemindeki hizmet ağından yararlanan herhangi bir operatör türü olabilir; mikro hizmetler platformuna bağımlı olan geliştiriciler, geliştiriciler, mühendisler, karar vericiler, mimarlar ve kuruluşlar dahil.

## Meshery, performans yönetimi içindir: pazar testleri ve karşılaştırmaları
Meshery, kullanıcıların bir hizmet ağı uygulamanın değerini, onu yürütmenin içerdiği yüksek yüke karşı tartmalarına yardımcı olur. Meshery, iş yükünüzün, altyapınızın ve hizmet ağı yapılandırmanızın permütasyonlarında gözlemlendiği gibi istek gecikmesi ve performansının istatistiksel analizini sağlar.
Meshery, gecikme ve performans talep etmenin yanı sıra küme düğümlerindeki CPU ve bellek ek yükünü de izler. Veri düzlemini ve kontrol düzlemini farklı altyapı ve iş yükleriyle ölçün.

<a href="https://raw.githubusercontent.com/layer5io/meshery/master/assets/img/readme/meshery_lifecycle_management.png"><img alt="Layer5 Service Mesh Community" src="{{ site.baseurl }}{% link assets/img/readme/meshery_lifecycle_management.png %}"  width="100%" align="center"/></a>
Marka performansını belirleyin ve zaman içinde değişen performansı bir standart olarak izleyin.

## Meshery, herhangi bir ağ hizmeti içindir
Altyapı çeşitliliği her kuruluş için bir gerçektir. İster tek bir ağ hizmeti veya birden çok hizmet ağı türü çalıştırıyor olun, Meshery çeşitli altyapıyı (veya bir ağın eksikliğini) desteklediğini görecektir.



- **Servis ağı adaptörleri mevcuttur ** - Mashery tarafından desteklenen servis ağı adaptörleri.

| Platform    &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;| Durum        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "stable" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }})                     |       {{ adapter.project_status }} |
{% endif -%}
{% endfor %}
<br>
- **Hizmet ağı bağdaştırıcıları devam ediyor ** - Topluluğa gönderilen ve geliştirme hizmeti ağ bağdaştırıcıları

| Platform    &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;| Durum        |
| :------------ | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "beta" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}
<br>
- **Destek arayan hizmet ağ bağdaştırıcıları ** - Topluluktan yardım arayan hizmet ağ bağdaştırıcıları.

| Platform        &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;     | Durum        |
| :------------               | :------------ |
{% for adapter in site.adapters -%}
{% if adapter.project_status == "alpha" -%}
| [{{ adapter.name }}]({{ site.baseurl }}{{ adapter.url }}) | {{ adapter.project_status }} |
{% endif -%}
{% endfor %}

## Topluluk
Bu proje topluluk tarafından yapılır ve her türlü işbirliğine açığız! [Fork here on Github](https://github.com/layer5io/meshery)

* Katılmak [haftalık topluluk toplantısı](https://docs.google.com/document/d/1c07UO9dS7_tFD-ClCWHIrEzRnzUJoFQ10EzfJTpS7FY/edit?usp=sharing) [Cuma 10: 00-11: 00 Merkezi Saat](/assets/projects/meshery/Meshery-Community-Meeting.ics). 
  * Gözlemlemek[topluluk kayıtları](https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0) ve abone olun [topluluk takvimi](https://bit.ly/2SbrRhe).

* Giriş [topluluk sürücüsü](https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA) (erişim talep etmek).

# FAQ 

## Neden Meshery kullanıyorsunuz?
* Açık kaynaklı bir proje olduğu için, ağlarda testi kolaylaştıran projeden bağımsız sağlayıcı.
* Fortio bir ağ testi uygulamasında paketlenmediğinden, yalnızca bir yük üreticisidir.
* Regpatrol kapalı kaynak olduğundan, ikili kod serbest bırakılmaz, bir ağ için gömülür ve bu ağın sağlayıcısı tarafından üretilir.

## Neden Meshery oluşturup başka bir markalama aracı kullanmıyorsunuz?
Meshery, diğer araçlardan farklı olarak hizmet ağlarının ve bunların iş yüklerinin karşılaştırılmasını kolaylaştırmak için özel olarak tasarlanmıştır. Ancak, regpatrol gibi hizmet ağlarının karşılaştırmalı değerlendirilmesi için kullanılan başka araçlar da vardır. Regpatrol, açık kaynaklı olmayan ve ikili biçimde kullanılamayan IBM tarafından kullanılır, Meshery ile aşağıdaki farklılıkları vardır:
- Telemetri: regpatrol, Mixer Prometheus adaptöründen telemetriyi alır ve IBM'in tescilli düğüm aracısını kullanır.
- Meshery, Prometheus Mikser adaptöründen elde edilir ve Prometheus düğüm dışa aktarıcısını kullanır.
- Trafik türü: regpatrol, yanıtları analiz edebilen ve işlevsel testler gerçekleştirebilen JMeter'ı kullanır.
- Meshery, yalnızca yapı ve performans testi için olan fortio kullanıyor.

# Kaynaklar

## Meshery Sunumları

- [O'Reilly OSCON 2020](https://conferences.oreilly.com/oscon/oscon-or)
- [O'Reilly Infrastructure & Ops 2020](https://conferences.oreilly.com/infrastructure-ops/io-ca/public/schedule/speaker/226795)
- [InnoTech Dallas 2020](https://innotechdallas2020.sched.com/event/aN7E/a-management-plane-for-service-meshes)
- [KubeCon EU 2020](https://kccnceu20.sched.com/event/Zetg/discreetly-studying-the-effects-of-individual-traffic-control-functions-lee-calcote-layer5?iframe=no&w=100%&sidebar=yes&bg=no)
- DockerCon 2020 ([deck](https://calcotestudios.com/talks/decks/slides-dockercon-2020-service-meshing-with-docker-desktop-and-webassembly.html), [video](https://www.youtube.com/watch?v=5BrbbKZOctw&list=PL3A-A6hPO2IN_HSU0pSfijBboiHggs5mC&index=4&t=0s))
- [Open Source 101 at Home](https://calcotestudios.com/talks/decks/slides-open-source-101-at-home-solving-the-service-mesh-adopters-dilemma.html)
- [Docker Captains Roundtable 2020](https://calcotestudios.com/talks/decks/slides-docker-captains-2020-meshery-the-multi-service-mesh-manager.html)
- [Cloud Native Austin 2020](https://www.meetup.com/Cloud-Native-Austin/events/267784090/)
- NSMCon 2019 talk ([video](https://www.youtube.com/watch?v=4xKixsDTtdM), [deck](https://calcotestudios.com/talks/decks/slides-nsmcon-kubecon-na-2019-adopting-network-service-mesh-with-meshery.html))
- [Service Mesh Day 2019](https://youtu.be/CFj1O_uyhhs)
- [DockerCon 2019 Open Source Summit](https://www.docker.com/dockercon/2019-videos?watch=open-source-summit-service-mesh)
- KubeCon EU 2019 ([video](https://www.youtube.com/watch?v=LxP-yHrKL4M&list=PLYjO73_1efChX9NuRaU7WocTbgrfvCoPE), [deck](https://calcotestudios.com/talks/decks/slides-kubecon-eu-2019-service-meshes-at-what-cost.html))
- [KubeCon EU 2019 Istio Founders Meetup](https://calcotestudios.com/talks/decks/slides-istio-meetup-kubecon-eu-2019-istio-at-scale-large-and-small.html)
- [Cloud Native Rejekts EU 2019](https://calcotestudios.com/talks/decks/slides-cloud-native-rejekts-2019-evaluating-service-meshes.html)
- [Container World 2019](https://calcotestudios.com/talks/decks/slides-container-world-2019-service-meshes-but-at-what-cost.html)
- Solving the Service Mesh Adopter’s Dilemma ([deck](https://calcotestudios.com/talks/decks/slides-open-source-101-at-home-solving-the-service-mesh-adopters-dilemma.html), [event](https://opensource101.com/sessions/solving-the-service-mesh-adopters-dilemma/),[video](https://www.youtube.com/watch?v=Q1zSWbO0RmI&list=PL3A-A6hPO2IN_HSU0pSfijBboiHggs5mC&index=2&t=0s))

## Başka kaynaklar
- [Hizmet Ağları karşılaştırması](https://layer5.io/landscape)
- [Servis ağları araçları](https://layer5.io/landscape#tools)
- [Hizmet Ağı ile ilgili kitaplar](https://layer5.io/books)
- [Hizmet Ağı Çalıştayları](https://layer5.io/workshops)

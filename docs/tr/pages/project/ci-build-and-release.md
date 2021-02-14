---
layout: default
title: Derleme ve Yayınlama (CI)
permalink: tr/project/build-and-release
language: tr
lang: tr
categories: tr
---

Meshery’nin derleme ve yayınlama sistemi, her biri farklı olaylarla tetiklenen farklı iş akışları şeklinde organize edilmiş birçok aracı bir araya getirir. Meshery’nin derleme ve yayınlama sistemi bir programa göre değil, olay odaklı çalışır. GitHub Eylemleri, Meshery’nin CI iş akışlarını tanımlamak için kullanılır. Yeni Meshery yapıları ve çeşitli bileşenleri, genellikle ilgili ana dallarıyla ilişkili olarak, itme, bırakma ve diğer benzer olaylarda otomatik olarak oluşturulur.

## Eserler

Bugün, Meshery ve Meshery bağdaştırıcıları Docker Hub'da bulunan Docker kapsayıcı görüntüleri olarak yayınlandı. Meshery bağdaştırıcıları işlem dışı bağdaştırıcılardır (yani ana Meshery ikili dosyasında derlenmemişlerdir) ve bu nedenle bağımsız derleme yapılarıdır. Docker görüntüleri oluşturma, git commit SHA ile etiketleme ve Docker Hub'a gönderme işlemi yapılıyor GitHub Eylemlerini kullanarak otomatik olarak.

### Yapı Kod Depoları

Derleme süreçlerinde üretilen yapılar, farklı genel depolarda ve farklı formatlarda yayınlanır ve saklanır.

| Location      | Project       | Repository    |
| ------------- | ------------- | ------------- |
| Docker Hub    | Meshery       | [https://hub.docker.com/r/layer5/meshery](https://hub.docker.com/r/layer5/meshery) |
| GitHub        | mesheryctl    | [https://github.com/layer5io/meshery/releases](https://github.com/layer5io/meshery/releases) |
| Docker Hub    | Meshery Adapter for \<service-mesh\> | https://hub.docker.com/r/layer5/meshery-\<service-mesh\> |
| Docs          | Meshery Documentation | [https://docs.meshery.io](https://docs.meshery.io) |
| GitHub        | [Service Mesh Performance](https://smp-spec.io) | [https://github.com/layer5io/service-mesh-performance](https://github.com/layer5io/service-mesh-performance) |

## Sırlar

İş akışının bazı bölümleri, görevlerini yerine getirmek için sırlar gerektirir. Bu sırlar, ilgili havuzlarda tanımlanır ve çalışma zamanı sırasında iş akışlarına erişilebilir. Şu anda tanımlanmış sırlar şunları içerir:

- `DOCKER_USERNAME`: Docker Hub kullanıcısının görüntüleri itmek için doğru ayrıcalıklara sahip kullanıcı adı
- `DOCKER_PASSWORD`: Docker Hub kullanıcısı için şifre
- `GO_VERSION`: 9 Aralık 2020 itibariyle 1.15
- `IMAGE_NAME`: Docker kapsayıcı görüntülerinin her biri için uygun görüntü adı. Hepsi "layer5io" kuruluşunun altındadır.
- `SLACK_BOT_TOKEN`: Meshery deposuna verilen yeni GitHub yıldızlarının bildirimi için kullanılır.
- CYPRESS_RECORD_KEY`: Cypress üzerindeki Layer5 hesabıyla entegrasyon için kullanılır.
- `GLOBAL_TOKEN`: Hiçbiri Sağlayıcı için performans testi sonuçlarının güvenli bir şekilde iletilmesi için kullanılır.

Docker Hub kullanıcısı "mesheryci", Docker Hub'daki "ciusers" ekibine aittir ve bu otomatikleştirilmiş yapıların itildiği hizmet hesabı olarak hareket eder. Yeni bir Docker Hub deposu her oluşturulduğunda, ciusers ekibine "Yönetici" (Docker Hub deposunda README'yi güncellemek için) izinleri vermemiz gerekir.

## Kontroller ve Testler
Meshery’nin CI iş akışı, birleştirme ve / veya herhangi bir dalı taahhüt etme sırasında çeşitli kontroller (aşağıdaki kısmi liste) içerir ve bozuk kodun ana birimle birleştirilmesini önlemek için istekleri ana dala çeker.

Toplu olarak, Meshery depoları genellikle aşağıdaki eylemlerden oluşan kaydetme ve çekme istekleri için CI iş akışına sahip olacaktır:

- Tüy bırakmayan kontrol (golint)
- Statik analiz kontrolü (statik kontrol)
- Veteriner (govet)
- Güvenlik kontrolleri (gosec)
- Birim testleri (testlere gidin)
- İnşa et (inşa et)
- İkili dosyaları GoReleaser aracılığıyla serbest bırakın (yalnızca Meshery deposundaki mesheryctl için)
- Docker derleme, etiketleme ve itme

## Otomatik Derlemeler

Tüm Meshery GitHub depoları GitHub Eylemleri ile yapılandırılır. Herhangi bir havuzun ana dalına karşı bir çekme isteği her gönderildiğinde, o deponun GitHub Eylemleri çağrılır (PR birleştirilmiş olsun ya da olmasın). Meshery havuzunda tanımlanan iş akışları genellikle (ancak her zaman değil) aşağıdaki eylemleri gerçekleştirir:

1. Docker konteyner görüntüsü oluşturmak için bir Docker derlemesini tetikleyin
1. iki Docker etiketi oluşturun:
   1. git merge SHA'yı içeren bir etiket
   1. söz konusu sürümün git etiketini içeren bir etiket (varsa)
1. Bu iki etiketin her birini yeni kapsayıcı resmine ve en son etikete atayın.
1. Yeni Docker etiketlerini ve görüntüsünü Docker Hub'a gönderin.

### Yapı `mesheryctl`

Özel bir durum olarak, meshery deposu her derleme sırasında üretilen ek bir yapı içerir. Bu yapı, yürütülebilir bir ikili olarak oluşturulmuş mesheryctl'dir. Farklı platform mimarileri ve işletim sistemlerinin bir kombinasyonu için mesheryctl oluşturma işini kolaylaştırmak için [GoReleaser] (https://goreleaser.com) kullanıyoruz. Daldan bağımsız olarak, meshery havuzuna yapılan her git commit ve git push için, GoReleaser işletim sistemini ve arşive özgü ikili dosyaları çalıştıracak ve oluşturacaktır (ancak bunları GitHub'da yayınlamayacaktır). Mesheryctl ikili dosyaları her seferinde bir çekme isteği ana ile birleştirildiğinde oluşturulsa da, yalnızca kararlı kanal yapıları yayınlanır (kalıcıdır).

### "mesheryctl" yi GitHub'a serbest bırakma

Yalnızca semantik bir sürüm numarası içeren bir git etiketi (ana dalda bir kaydetme) mevcut olduğunda GoReleaser yürütür, arşivleri oluşturur ve ayrıca arşivleri [Meshery'nin GitHub sürümlerinde] (https://github.com/layer5io / meshery / release) otomatik olarak. GoReleaser, aşağıdaki işletim sistemi, ARCH kombinasyonu için yapılar oluşturacak şekilde yapılandırılmıştır:

- Darwin - i386, x86_64
- Linux - i386, x86_64
- Windows - i386, x86_64
- FreeBSD - i386, x86_64

Yapılar, tüm işletim sistemleri için tar.gz arşivi olarak sunulacaktır. mesheryctl, yaygın olarak kullanılan paket yöneticileri için paketler halinde paketlenmiştir: homebrew ve scoop.

#### Homebrew

GoReleaser, mesheryctl için bir demleme formülünün oluşturulmasını kolaylaştırır. [Homebrew-tap] (https://github.com/layer5io/homebrew-tap) deposu, Layer5’in brew formüllerinin bulunduğu yerdir.

#### Kepçe

GoReleaser, mesheryctl için bir Scoop uygulamasının oluşturulmasını kolaylaştırır. [Scoop-buck] (https://github.com/layer5io/scoop-bucket) deposu, Layer5’in Scoop paketinin konumudur.

## Sürüm Oluşturma

Meshery, Meshery Adapter ve Performance Benchmark Specification sürümleri için yaygın olarak kullanılan anlamsal sürümlemeyi takip ediyoruz. MAJOR.MINOR.PATCH.BUILD sürüm numarası verildiğinde, şunları artırın:

- MAJOR sürümü - nadiren uyumsuz API değişiklikleri potansiyeline sahip önemli değişiklikler.
- KÜÇÜK sürüm - geriye doğru uyumlu bir şekilde işlevsellik ekleyin.
- PATCH sürümü - çoğunlukla hata ve güvenlik düzeltmeleri için.
- AlPHA / BETA / RC - gelecek bir sürümün erken testini kolaylaştırmak için kullanılır.

### Bileşen Sürüm Oluşturma

Meshery, bir sunucu, adaptörler, UI ve CLI dahil olmak üzere bir dizi bileşenden oluşur. Bir uygulama olarak Meshery, bu farklı fonksiyonel bileşenlerin bir bileşimidir. Meshery’nin tüm bileşenleri genel olarak toplu bir birim olarak (birlikte) konuşlandırılırken, her bileşen bağımsız olarak, gevşek bir şekilde bağlanmalarına ve işlevsellik üzerinde bağımsız olarak yinelemelerine izin verecek şekilde versiyonlanır. Bazı bileşenlerin eşzamanlı olarak yükseltilmesi gerekirken, diğerleri bağımsız olarak yükseltilebilir. Daha fazla bilgi için bkz. [Meshery Yükseltme] (/ guide / upgrade).

GitHub yayın etiketleri, anlamsal bir sürüm numarası içerecektir. Anlamsal sürüm numaralarının, anlamsal sürüm numarasıyla ana şubede ilgili bir kesinleştirme etiketlenmesiyle manuel olarak yönetilmesi gerekecektir (örnek: v1.2.3).

## Yayın Süreci

Meshery sürümlerinin dokümantasyonu bir sürümler tablosu ve sürüm notları içerir ve her sürümde güncellenmelidir.

### Otomatik Yayınlar

Sürümler, bir sürümü yayınlayan sürüm ekibinin bir üyesi tarafından manuel olarak tetiklenir. Sürüm adları ve sürüm etiketlerinin yayıncı ekip üyesi tarafından atanması gerekir. GitHub Action iş akışları, gerekli adımları çalıştırmayı ve tüm yapıtları (örneğin, ikili ve docker görüntüleri) yayınlamayı tetikleyecek ve ilgilenecektir.

### İş Akışı Tetikleyicileri

Aşağıdaki olaylar bir veya daha fazla iş akışını tetikleyecektir:

1. Etiketli Sürüm
1. Taahhüt ana şubeye itildi
1. PR açıldı veya PR şubesine itildi
1. PR, ana şubeyle birleştirildi

### Sürüm notları

GitHub Eylemlerinin kullanımı otomatikleştirilmiş derlemeleri kolaylaştırırken, ReleaseDrafter otomatik sürüm notlarını ve sürüm oluşturmayı kolaylaştırmaya yardımcı oluyor.
### Sürüm Notları Oluşturma

ReleaseDrafter, bir GitHub etiketi ve yayın taslağı oluşturur. ReleaseDrafter eylemi tetiklenecek ve yapılandırma kurulumuna göre otomatik olarak sürüm notları hazırlayacaktır. ReleaseDrafter, önceki sürümden sonra ana hale bir taahhüt yapılır yapılmaz sürümleri taslak haline getirir. GitHub Action, ReleaseDrafter, anlamsal sürümlerle uyumludur ve önceki sürüm sürümüne bakarak anlamsal sürüm numarasını otomatik olarak artırmak için kullanılır.

#### Otomatik Sürüm Notları Yayınlama

Meshery Docs'a sürüm notlarının yayınlanması otomatikleştirilmiştir. Bir yayınlama olayıyla tetiklenen bir iş akışı, Meshery deposunu teslim alır, otomatik olarak hazırlanan sürüm notlarını Meshery Docs içindeki bir Jekyll koleksiyonuna kopyalar ve bir çekme talebi oluşturur.

#### Otomatik Çekme İsteği Etiketleyici

GitHub Sorun etiketleme botu, hangi dosyaların hangi dizinlerde değiştiğine bağlı olarak sorunlara otomatik olarak etiket atayacak şekilde yapılandırılmıştır. Örneğin, "/ docs / **" klasöründeki dosyalarda değişiklikler olan bir çekme isteği, "alan / dokümanlar" etiketini alacaktır. "Alan / belgeler" etiketinin varlığı, Meshery Belgelerinin belge derlemelerini ve Netlify derlemelerini tetiklemek için kullanılır. İş akışlarını tetiklemek için benzer etiketler atanır ve kullanılır veya hangi iş akışlarının veya bir iş akışında hangi adımların çalıştırılacağını belirlemek için iş akışlarında koşullu bayraklar olarak kullanılır.

## Yayın Kanalları

Meshery ve bileşenlerine yönelik yapıların yapıları iki farklı sürüm kanalı altında yayınlanır, böylece hem Meshery kullanıcıları hem de Meshery geliştiricileri için gelişmiş kontroller sağlanabilir. İki yayın kanalı, * uç * ve * kararlı * yayın kanallarıdır.

Kararlı sürümlere kıyasla, uç sürümler çok daha sık meydana gelir. Uç sürümler, ana sürüm için birleştirme kararlı bir sürüm değilse, ana sürüm için her birleştirme ile yapılır. İş akışında bir GitHub yayın etiketi de mevcut olduğunda, her bir birleştirmeyle ana olarak kararlı sürümler yapılır.

### Sabit Kanal

Aşağıda, sürüm kanallarının ve bunları ayırt etmek için kullanılan docker etiketlerinin bir örneği verilmiştir. En son etiket, yalnızca kararlı sürüm kanalındaki görüntülere uygulanacaktır. İşte iki farklı görüntüye sahip iki sürüm.

** En Son Kararlı Görüntü **

- layer5 / meshery: stabil-en son
- layer5 / meshery: kararlı-v0.4.1
- layer5 / meshery: kararlı-324vdgb (sha)

** Daha Eski Sabit Görüntü **

- layer5 / meshery: kararlı-v0.4.0
- layer5 / meshery: kararlı-289d02 (sha)

Oluşturulan her docker görüntüsü, kenar etiketlerini veya kararlı sekmeleri alır. Hangi resim etiketi setinin atandığı, bir yayın etiketinin mevcut olup olmadığına göre belirlenir. Diğer bir deyişle, kararlı kanal docker görüntüleri yalnızca bir yayın etiketi (ör. V0.4.1) varlığında "kararlı" etiketleri alır.

### Edge Kanalı

Kenar yayın kanalı genellikle daha az test edilmiş, daha az "pişmiş" kod içerir. "Edge" in birincil nedeni, katkıda bulunanların ve ileri düzey kullanıcıların özelliklere daha erken erişmesine izin vermektir. Bazı özellikler, tolerans ve sabırla kullanıcıların bunları denemesine izin vererek en iyi şekilde kolaylaştırılan testlere ihtiyaç duyar.

Kararlı ve uç sürümlerin her ikisi de aynı Docker Hub deposunda yayınlanır. Docker Hub depoları, yayın kanallarını görüntü etiketine göre farklılaştırır. Aşağıdaki Docker görüntüleri etiketleme kuralı izlenir:

** En Son Kenar Görüntüsü **

- layer5 / meshery: en son kenar
- layer5 / meshery: edge-289d02 (sha)

** Eski Kenar Resmi **

- layer5 / meshery: edge-324vdgb (sha)


### Meshery Yayın Kanalları Arasında Geçiş Yapma

Kullanıcılar, boş zamanlarında yayın kanalları arasında geçiş yapma yetkisine sahiptir.

#### mesheryctl Kullanarak Yayın Kanallarını Değiştirme

Kullanıcılar, mesheryctl'yi yayın kanalları arasında geçiş yapmak için kullanabilir, ör. `mesheryctl sistem kanalı [kararlı | kenar] '. Alternatif olarak, kullanıcılar meshery.yaml / Kubernetes manifest dosyalarındaki docker görüntü etiketlerini güncelleyerek kanallar arasında manuel olarak geçiş yapabilir. Bu komut, farklı Docker kapsayıcı görüntüleri için yayın kanalına uygun etiketlere sahip bir meshery.yml (bir docker-compose dosyası) oluşturur.

#### Meshery Kullanıcı Arayüzünde Yayın Kanalı ve Sürüm Bilgilerini Görüntüleme

Kullanıcılara, Meshery kullanıcı arayüzünün Tercihler alanında, Meshery dağıtımının yayın kanalı aboneliğinin yeni ayarı gösterilir, böylece insanlar isterlerse kanallar arasında geçiş yapmak için kullanıcı arayüzünü alternatif olarak kullanabilir. Meshery bağdaştırıcılarının sürüm numaraları da kullanıcı arayüzünde gösterilir.

## Yayın Temposu

Meshery projesinin küçük sürümleri, bu zamanlar arasında isteğe bağlı olarak yapılan yama sürümleriyle sık sık (ortalama aylık olarak) yayınlanır. Projenin henüz hata düzeltmeleriyle sürdürülen uzun vadeli sürümleri yok. En son sürümde gerektiğinde hata düzeltmeleri ve yamalar yayınlanacaktır.

### Sürüm Desteği

Layer5'ten genel topluluk desteği ve ticari destek mevcuttur. Ayrı olarak, üçüncü şahıslar ve ortaklar daha uzun vadeli destek çözümleri sunabilir.

#### Pre v1.0

Proje, mimarideki değişimler için esnekliği korurken işlevsellik, kalite ve benimsemeye odaklanır.

#### Yayın v1.0

1.0 sürümü bir kez yapıldığında, yaklaşık ayda bir kez, proje yöneticileri bu günlük yapılardan birini alıp bir dizi ek yeterlilik testinden geçirecek ve yapıyı Kararlı sürüm olarak etiketleyecektir. Yaklaşık üç ayda bir, proje yöneticileri bu Kararlı sürümlerden birini alır, birkaç test daha yapar ve yapıyı Uzun Süreli Destek (LTS) sürümü olarak etiketler. Son olarak, bir LTS sürümünde yanlış bir şey bulursak, yamalar çıkarırız.

Farklı tipler (Günlük, Sabit, LTS), farklı ürün kalite seviyelerini ve Meshery ekibinden farklı destek seviyelerini temsil eder. Bu bağlamda destek, kritik sorunlar için yama sürümleri üreteceğimiz ve teknik yardım sunacağımız anlamına gelir.

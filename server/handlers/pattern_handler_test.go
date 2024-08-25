package handlers

// import (
// 	"fmt"
// 	"os"
// 	"path"
// 	"testing"

// 	"github.com/gofrs/uuid"
// 	"github.com/layer5io/meshery/server/models"
// 	"github.com/layer5io/meshkit/logger"
// 	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
// 	"github.com/layer5io/meshkit/validator"
// 	"github.com/meshery/schemas/models/v1beta1/pattern"
// 	"github.com/spf13/viper"
// )

// var testPatternFile = `
// name: GuestBook App
// version: 0.0.11
// services:
//   frontend-dtrds:
//     annotations: {}
//     apiVersion: v1
//     dependsOn: []
//     id: null
//     isAnnotation: null
//     labels: {}
//     model: kubernetes
//     name: frontend
//     namespace: default
//     settings:
//       spec:
//         ports:
//         - port: 80
//           protocol: TCP
//         - port: 80
//           protocol: TCP
//         selector:
//           app: guestbook
//           tier: frontend
//     traits:
//       meshmap:
//         edges:
//         - data:
//             id: 499a0982-138b-48f0-bfc0-4b2c9f711510
//             metadata:
//               port: 80
//               protocol: TCP
//             source: 4ae55da0-82eb-405e-9410-75623cee1043
//             subType: Network
//             target: 43a7b935-6109-4f5c-84fb-5e3454f9c6d1
//           style:
//             control-point-distances: null
//             control-point-weights: "0.5"
//             curve-style: bezier
//             haystack-radius: "0"
//             line-color: rgb(153,153,153)
//             line-style: dotted
//             opacity: "1"
//             segment-distances: 20px
//             segment-weights: "0.5"
//             target-arrow-shape: vee
//             taxi-direction: auto
//             taxi-turn: 50%
//             taxi-turn-min-distance: 10px
//             width: 1.5px
//         fieldRefData: {}
//         id: 4ae55da0-82eb-405e-9410-75623cee1043
//         label: frontend
//         meshmodel-metadata:
//           genealogy: ""
//           isCustomResource: false
//           isNamespaced: true
//           logoURL: https://github.com/cncf/artwork/blob/master/projects/kubernetes/icon/white/kubernetes-icon-white.svg
//           model: kubernetes
//           modelDisplayName: Kubernetes
//           primaryColor: '#326CE5'
//           published: true
//           secondaryColor: '#7aa1f0'
//           shape: round-triangle
//           styleOverrides:
//             background-fit: none
//             background-position-y: 4.5
//             height: 16
//             padding: 12
//             width: 17
//             z-index: 3
//           subCategory: Scheduling & Orchestration
//           svgColor: ui/public/static/img/meshmodels/kubernetes/color/apiservice-color.svg
//           svgComplete: ""
//           svgWhite: ui/public/static/img/meshmodels/kubernetes/white/apiservice-white.svg
//         position:
//           posX: 710
//           posY: 570
//         whiteboardData:
//           style:
//             z-index: 15
//     type: Service
//     version: ""
//   frontend-mzdsv:
//     annotations: {}
//     apiVersion: apps/v1
//     dependsOn: []
//     id: null
//     isAnnotation: null
//     labels: {}
//     model: kubernetes
//     name: frontend
//     namespace: ""
//     settings:
//       spec:
//         replicas: 3
//         selector:
//           match Labels:
//             app: guestbook
//             tier: frontend
//         template:
//           metadata:
//             labels:
//               app: guestbook
//               tier: frontend
//           spec:
//             containers:
//             - env:
//               - name: GET_HOSTS_FROM
//                 value: dns
//               image: gcr.io/google_samples/gb-frontend:v5
//               name: php-redis
//               ports:
//               - container Port: 80
//               resources:
//                 requests:
//                   cpu: 100m
//                   memory: 100Mi
//     traits:
//       meshmap:
//         edges:
//         - data:
//             id: 499a0982-138b-48f0-bfc0-4b2c9f711510
//             metadata:
//               port: 80
//               protocol: TCP
//             source: 4ae55da0-82eb-405e-9410-75623cee1043
//             subType: Network
//             target: 43a7b935-6109-4f5c-84fb-5e3454f9c6d1
//           style:
//             control-point-distances: null
//             control-point-weights: "0.5"
//             curve-style: bezier
//             haystack-radius: "0"
//             line-color: rgb(153,153,153)
//             line-style: dotted
//             opacity: "1"
//             segment-distances: 20px
//             segment-weights: "0.5"
//             target-arrow-shape: vee
//             taxi-direction: auto
//             taxi-turn: 50%
//             taxi-turn-min-distance: 10px
//             width: 1.5px
//         fieldRefData: {}
//         id: 43a7b935-6109-4f5c-84fb-5e3454f9c6d1
//         label: frontend
//         meshmodel-metadata:
//           styleOverrides:
//             z-index: 18
//         position:
//           posX: 710
//           posY: 690
//         whiteboardData:
//           style:
//             z-index: 9
//     type: Deployment
//     version: ""
// `

// func TestPatternFileConversionFromV1Alpha1ToV1Beta1(t *testing.T) {

// 	v1beta1PatternFile := &pattern.PatternFile{}
// 	home, _ := os.UserHomeDir()

// 	viper.SetDefault("USER_DATA_FOLDER", path.Join(home, ".meshery", "config"))

// 	dbHandler := models.GetNewDBInstance()
// 	regManager, err := meshmodel.NewRegistryManager(dbHandler)
// 	if err != nil {
// 		t.Fatal(err)
// 		return
// 	}
// 	log, err := logger.New("pattern_handler-tests", logger.Options{
// 		LogLevel: 3,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 		return
// 	}
// 	handler := Handler{
// 		registryManager: regManager,
// 		log:             log,
// 	}
// 	t.Run("TestPatternFileConversionFromV1Alpha1ToV1Beta1", func(t *testing.T) {
// 		// var beta1PatternFileStr string
// 		var err error
// 		v1beta1PatternFile, _, err = handler.convertV1alpha2ToV1beta1(testPatternFile, uuid.Nil)
// 		if err != nil {
// 			fmt.Println(err)
// 			t.Fatal(err)
// 			return
// 		}

// 		// b, _ := json.Marshal(v1beta1PatternFile)
// 		// _ = os.WriteFile("test-beta1.json", b, 0655)
// 		// _ = os.WriteFile("test-beta1.yml", []byte(beta1PatternFileStr), 0655)

// 	})

// 	t.Run("TestValidateConvertedPatternFile", func(t *testing.T) {
// 		cueSchema, err := validator.GetSchemaFor("")
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		err = validator.Validate(cueSchema, v1beta1PatternFile)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	})
// }

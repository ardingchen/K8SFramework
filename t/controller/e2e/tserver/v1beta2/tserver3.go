package v1beta2

import (
	"context"
	"e2e/scaffold"
	"github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	k8sCoreV1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	patchTypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	tarsCrdV1Beta2 "k8s.tars.io/crd/v1beta2"
	tarsMetaTools "k8s.tars.io/meta/tools"
	tarsMetaV1Beta2 "k8s.tars.io/meta/v1beta2"
	"time"
)

var _ = ginkgo.Describe("try create tars server and check filed", func() {
	opts := &scaffold.Options{
		Name:      "default",
		K8SConfig: scaffold.GetK8SConfigFile(),
		SyncTime:  1500 * time.Millisecond,
	}

	s := scaffold.NewScaffold(opts)

	var tsLayout *tarsCrdV1Beta2.TServer
	var Resource = "test-testserver"
	var App = "Test"
	var Server = "TestServer"
	var Template = "tt.cpp"
	var FirstObj = "FirstObj"
	var SecondObj = "SecondObj"

	ginkgo.BeforeEach(func() {
		ttLayout := &tarsCrdV1Beta2.TTemplate{
			ObjectMeta: k8sMetaV1.ObjectMeta{
				Name:      Template,
				Namespace: s.Namespace,
			},
			Spec: tarsCrdV1Beta2.TTemplateSpec{
				Content: "tt.cpp content",
				Parent:  Template,
			},
		}
		_, err := s.CRDClient.CrdV1beta2().TTemplates(s.Namespace).Create(context.TODO(), ttLayout, k8sMetaV1.CreateOptions{})
		assert.Nil(ginkgo.GinkgoT(), err)
		time.Sleep(s.Opts.SyncTime)

		tsLayout = &tarsCrdV1Beta2.TServer{
			ObjectMeta: k8sMetaV1.ObjectMeta{
				Name:      Resource,
				Namespace: s.Namespace,
			},
			Spec: tarsCrdV1Beta2.TServerSpec{
				App:       App,
				Server:    Server,
				SubType:   tarsCrdV1Beta2.TARS,
				Important: 5,
				Tars: &tarsCrdV1Beta2.TServerTars{
					Template:    Template,
					Profile:     "",
					AsyncThread: 3,
					Servants: []*tarsCrdV1Beta2.TServerServant{
						{
							Name:       FirstObj,
							Port:       10000,
							Thread:     3,
							Connection: 1000,
							Capacity:   1000,
							Timeout:    1000,
							IsTars:     true,
							IsTcp:      true,
						},
						{
							Name:       SecondObj,
							Port:       10001,
							Thread:     3,
							Connection: 1000,
							Capacity:   1000,
							Timeout:    1000,
							IsTars:     true,
							IsTcp:      true,
						},
					},
				},
				K8S: tarsCrdV1Beta2.TServerK8S{
					AbilityAffinity: tarsCrdV1Beta2.None,
					NodeSelector:    []k8sCoreV1.NodeSelectorRequirement{},
					ImagePullPolicy: k8sCoreV1.PullAlways,
					LauncherType:    tarsCrdV1Beta2.Background,
				},
			},
		}
		_, err = s.CRDClient.CrdV1beta2().TServers(s.Namespace).Create(context.TODO(), tsLayout, k8sMetaV1.CreateOptions{})
		assert.Nil(ginkgo.GinkgoT(), err)
		time.Sleep(s.Opts.SyncTime)
	})

	ginkgo.It("check filed value", func() {
		tserver, err := s.CRDClient.CrdV1beta2().TServers(s.Namespace).Get(context.TODO(), Resource, k8sMetaV1.GetOptions{})
		assert.Nil(ginkgo.GinkgoT(), err)
		assert.NotNil(ginkgo.GinkgoT(), tserver)

		expectedLabels := map[string]string{
			tarsMetaV1Beta2.TServerAppLabel:  App,
			tarsMetaV1Beta2.TServerNameLabel: Server,
			tarsMetaV1Beta2.TemplateLabel:    Template,
			tarsMetaV1Beta2.TSubTypeLabel:    string(tarsCrdV1Beta2.TARS),
		}
		assert.True(ginkgo.GinkgoT(), scaffold.CheckLeftInRight(expectedLabels, tserver.Labels))
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.Important)
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.Tars)
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.Tars.Servants)
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.Tars.Template)
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.Tars.AsyncThread)
		assert.Nil(ginkgo.GinkgoT(), tserver.Spec.Tars.Ports)
		assert.Nil(ginkgo.GinkgoT(), tserver.Spec.Normal)
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.K8S)
		assert.False(ginkgo.GinkgoT(), tserver.Spec.K8S.HostIPC)
		assert.False(ginkgo.GinkgoT(), tserver.Spec.K8S.HostNetwork)
		assert.Nil(ginkgo.GinkgoT(), tserver.Spec.K8S.HostPorts)
		assert.False(ginkgo.GinkgoT(), tserver.Spec.K8S.DaemonSet)
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.K8S.Replicas)
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.K8S.NodeSelector)
		assert.NotNil(ginkgo.GinkgoT(), tserver.Spec.K8S.ImagePullPolicy)

		expectedReadinessGates := tarsMetaV1Beta2.TPodReadinessGate
		assert.Equal(ginkgo.GinkgoT(), expectedReadinessGates, tserver.Spec.K8S.ReadinessGate)
	})

	ginkgo.It("try remove immutable filed", func() {
		removeFileds := map[string]interface{}{
			"/spec/app":     nil,
			"/spec/server":  nil,
			"/spec/subType": nil,
			"/spec/tars":    nil,
			"/spec/k8s":     nil,
		}
		for k := range removeFileds {
			jsonPath := tarsMetaTools.JsonPatch{
				{
					OP:   tarsMetaTools.JsonPatchRemove,
					Path: k,
				},
			}
			bs, _ := json.Marshal(jsonPath)
			_, err := s.CRDClient.CrdV1beta2().TConfigs(s.Namespace).Patch(context.TODO(), Resource, patchTypes.JSONPatchType, bs, k8sMetaV1.PatchOptions{})
			assert.NotNil(ginkgo.GinkgoT(), err)
		}
	})

	ginkgo.It("try update immutable filed", func() {
		updateFileds := map[string]string{
			"/spec/app":     scaffold.RandStringRunes(3),
			"/spec/server":  scaffold.RandStringRunes(5),
			"/spec/subType": scaffold.RandStringRunes(5),
		}
		for k, v := range updateFileds {
			jsonPath := tarsMetaTools.JsonPatch{
				{
					OP:    tarsMetaTools.JsonPatchReplace,
					Path:  k,
					Value: v,
				},
			}
			bs, _ := json.Marshal(jsonPath)
			_, err := s.CRDClient.CrdV1beta2().TConfigs(s.Namespace).Patch(context.TODO(), Resource, patchTypes.JSONPatchType, bs, k8sMetaV1.PatchOptions{})
			assert.NotNil(ginkgo.GinkgoT(), err)
		}
	})

	ginkgo.It("try update subType", func() {
		jsonPath := tarsMetaTools.JsonPatch{
			{
				OP:    tarsMetaTools.JsonPatchReplace,
				Path:  "/spec/subType",
				Value: tarsCrdV1Beta2.Normal,
			},
			{
				OP:   tarsMetaTools.JsonPatchRemove,
				Path: "/spec/tars",
			},
			{
				OP:   tarsMetaTools.JsonPatchAdd,
				Path: "/spec/normal",
				Value: &tarsCrdV1Beta2.TServerNormal{
					Ports: []*tarsCrdV1Beta2.TServerPort{},
				},
			},
		}
		bs, _ := json.Marshal(jsonPath)
		_, err := s.CRDClient.CrdV1beta2().TConfigs(s.Namespace).Patch(context.TODO(), Resource, patchTypes.JSONPatchType, bs, k8sMetaV1.PatchOptions{})
		assert.NotNil(ginkgo.GinkgoT(), err)
	})
})

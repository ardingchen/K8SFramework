package v1beta2

import (
	"fmt"
	"hash/crc32"
	k8sAdmissionV1 "k8s.io/api/admission/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	tarsCrdV1beta2 "k8s.tars.io/crd/v1beta2"
	tarsMetaTools "k8s.tars.io/meta/tools"
	"time"
)

func mutatingCreateTConfig(requestAdmissionView *k8sAdmissionV1.AdmissionReview) ([]byte, error) {
	tconfig := &tarsCrdV1beta2.TConfig{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, tconfig)

	var jsonPatch tarsMetaTools.JsonPatch

	if tconfig.Labels == nil {
		labels := map[string]string{}
		jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
			OP:    tarsMetaTools.JsonPatchAdd,
			Path:  "/metadata/labels",
			Value: labels,
		})
	}

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1ServerApp",
		Value: tconfig.App,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1ServerName",
		Value: tconfig.Server,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1ConfigName",
		Value: tconfig.ConfigName,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1PodSeq",
		Value: tconfig.PodSeq,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1Activated",
		Value: fmt.Sprintf("%t", tconfig.Activated),
	})

	versionString := fmt.Sprintf("%s-%x", time.Now().Format("20060102030405"), crc32.ChecksumIEEE([]byte(tconfig.Name)))
	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/version",
		Value: versionString,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1Version",
		Value: versionString,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/updateTime",
		Value: k8sMetaV1.Now().ToUnstructured(),
	})

	if jsonPatch != nil {
		return json.Marshal(jsonPatch)
	}
	return nil, nil
}

func mutatingUpdateTConfig(requestAdmissionView *k8sAdmissionV1.AdmissionReview) ([]byte, error) {
	tconfig := &tarsCrdV1beta2.TConfig{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, tconfig)

	var jsonPatch tarsMetaTools.JsonPatch

	if tconfig.Labels == nil {
		labels := map[string]string{}
		jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
			OP:    tarsMetaTools.JsonPatchAdd,
			Path:  "/metadata/labels",
			Value: labels,
		})
	}

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1ServerApp",
		Value: tconfig.App,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1ServerName",
		Value: tconfig.Server,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1ConfigName",
		Value: tconfig.ConfigName,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1PodSeq",
		Value: tconfig.PodSeq,
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1Activated",
		Value: fmt.Sprintf("%t", tconfig.Activated),
	})

	jsonPatch = append(jsonPatch, tarsMetaTools.JsonPatchItem{
		OP:    tarsMetaTools.JsonPatchAdd,
		Path:  "/metadata/labels/tars.io~1Version",
		Value: tconfig.Version,
	})

	if jsonPatch != nil {
		return json.Marshal(jsonPatch)
	}
	return nil, nil
}

package v1beta3

import (
	"fmt"
	k8sAdmissionV1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/json"
	tarsCrdV1beta3 "k8s.tars.io/crd/v1beta3"
	tarsMetaV1beta3 "k8s.tars.io/meta/v1beta3"
	"strings"
	"tarscontroller/controller"
)

func validTServer(newTServer, oldTServer *tarsCrdV1beta3.TServer, clients *controller.Clients, informers *controller.Informers) error {

	if oldTServer != nil {
		if newTServer.Spec.App != oldTServer.Spec.App {
			return fmt.Errorf(tarsMetaV1beta3.FiledImmutableError, "tserver", ".spec.app")
		}

		if newTServer.Spec.Server != oldTServer.Spec.Server {
			return fmt.Errorf(tarsMetaV1beta3.FiledImmutableError, "tserver", ".spec.server")
		}

		if newTServer.Spec.SubType != oldTServer.Spec.SubType {
			return fmt.Errorf(tarsMetaV1beta3.FiledImmutableError, "tserver", ".spec.subType")
		}

		if oldTServer.Spec.Tars == nil {
			if newTServer.Spec.Tars != nil {
				return fmt.Errorf(tarsMetaV1beta3.FiledImmutableError, "tserver", ".spec.tars")
			}
		}

		if oldTServer.Spec.Normal == nil {
			if newTServer.Spec.Normal != nil {
				return fmt.Errorf(tarsMetaV1beta3.FiledImmutableError, "tserver", ".spec.normal")
			}
		}
	}

	namespace := newTServer.Namespace

	if newTServer.Name != strings.ToLower(newTServer.Spec.App)+"-"+strings.ToLower(newTServer.Spec.Server) {
		return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", "unexpected resource name")
	}

	if len(newTServer.Name) >= tarsMetaV1beta3.MaxTServerName {
		return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", "length of resource name should less then 59")
	}

	portNames := map[string]interface{}{}
	portValues := map[int32]interface{}{}

	if newTServer.Spec.Tars != nil {

		for _, servant := range newTServer.Spec.Tars.Servants {
			portName := strings.ToLower(servant.Name)
			portValue := servant.Port

			if portValue == tarsMetaV1beta3.NodeServantPort {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("servants port value should not equal %d", tarsMetaV1beta3.NodeServantPort))
			}

			if _, ok := portNames[portName]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate servant name value %s", servant.Name))
			}

			if _, ok := portValues[portValue]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate port value %d", servant.Port))
			}

			portNames[portName] = nil
			portValues[portValue] = nil
		}

		for _, port := range newTServer.Spec.Tars.Ports {
			portName := strings.ToLower(port.Name)
			portValue := port.Port

			if portValue == tarsMetaV1beta3.NodeServantPort {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("port value should not equal %d", tarsMetaV1beta3.NodeServantPort))
			}

			if _, ok := portNames[portName]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate port name value %s", port.Name))
			}

			if _, ok := portValues[portValue]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate port value %d", port.Port))
			}
			portNames[portName] = nil
			portValues[portValue] = nil
		}

		templateName := newTServer.Spec.Tars.Template
		_, err := informers.TTemplateInformer.Lister().ByNamespace(namespace).Get(templateName)
		if err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf(tarsMetaV1beta3.ResourceGetError, "ttemplate", namespace, templateName, err.Error())
			}
			return fmt.Errorf(tarsMetaV1beta3.ResourceNotExistError, "ttemplate", namespace, templateName)
		}
	} else if newTServer.Spec.Normal != nil {
		for _, port := range newTServer.Spec.Normal.Ports {
			portName := strings.ToLower(port.Name)
			portValue := port.Port

			if _, ok := portNames[portName]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate port name value %s", port.Name))
			}

			if _, ok := portValues[portValue]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate port value %d", port.Port))
			}
			portNames[portName] = nil
			portValues[portValue] = nil
		}
	}

	if newTServer.Spec.K8S.HostPorts != nil {

		hostPortPorts := map[int32]interface{}{}
		hostPortNameRefs := map[string]interface{}{}

		for _, hostPort := range newTServer.Spec.K8S.HostPorts {
			nameRef := strings.ToLower(hostPort.NameRef)
			if _, ok := portNames[nameRef]; !ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("port name %s not exist", hostPort.NameRef))
			}

			if _, ok := hostPortNameRefs[nameRef]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate nameRef value %s", hostPort.NameRef))
			}

			if _, ok := hostPortPorts[hostPort.Port]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate port value %d", hostPort.Port))
			}

			hostPortNameRefs[nameRef] = nil
			hostPortPorts[hostPort.Port] = nil
		}
	}

	if newTServer.Spec.K8S.Mounts != nil {
		mountsNames := map[string]interface{}{}

		for i := range newTServer.Spec.K8S.Mounts {

			mount := &newTServer.Spec.K8S.Mounts[i]

			if _, ok := mountsNames[mount.Name]; ok {
				return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("duplicate .mounts.name value %s", mount.Name))
			}

			if mount.Source.TLocalVolume != nil || mount.Source.PersistentVolumeClaimTemplate != nil {
				if newTServer.Spec.K8S.DaemonSet {
					return fmt.Errorf(tarsMetaV1beta3.ResourceInvalidError, "tserver", fmt.Sprintf("can not use TLocalVolue and PersistentVolumeClaimTemplate when .daemonSet value is true"))
				}
			}

			mountsNames[mount.Name] = nil
		}
	}
	return nil
}

func validCreateTServer(clients *controller.Clients, informers *controller.Informers, view *k8sAdmissionV1.AdmissionReview) error {
	newTServer := &tarsCrdV1beta3.TServer{}
	_ = json.Unmarshal(view.Request.Object.Raw, newTServer)
	return validTServer(newTServer, nil, clients, informers)
}

func validUpdateTServer(clients *controller.Clients, informers *controller.Informers, view *k8sAdmissionV1.AdmissionReview) error {
	newTServer := &tarsCrdV1beta3.TServer{}
	_ = json.Unmarshal(view.Request.Object.Raw, newTServer)

	oldTServer := &tarsCrdV1beta3.TServer{}
	_ = json.Unmarshal(view.Request.OldObject.Raw, oldTServer)

	return validTServer(newTServer, oldTServer, clients, informers)
}

func validDeleteTServer(clients *controller.Clients, informers *controller.Informers, view *k8sAdmissionV1.AdmissionReview) error {
	return nil
}

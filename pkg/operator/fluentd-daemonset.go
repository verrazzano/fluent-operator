package operator

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
)

func MakeFluentdDaemonSet(fd fluentdv1alpha1.Fluentd) *appsv1.DaemonSet {

	ports := makeFluentdPorts(fd)

	labels := map[string]string{
		"app.kubernetes.io/name":      fd.Name,
		"app.kubernetes.io/instance":  "fluentd",
		"app.kubernetes.io/component": "fluentd",
	}

	if len(fd.Labels) > 0 {
		for k, v := range fd.Labels {
			if _, ok := labels[k]; !ok {
				labels[k] = v
			}
		}
	}

	defaultFsGroup := DefaultFsGroup

	daemonSet := appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fd.Name,
			Namespace:   fd.Namespace,
			Labels:      labels,
			Annotations: fd.Annotations,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fd.Name,
					Namespace:   fd.Namespace,
					Labels:      labels,
					Annotations: fd.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: fd.Name,
					ImagePullSecrets:   fd.Spec.ImagePullSecrets,
					Volumes: []corev1.Volume{
						{
							Name: SecretVolName,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: fmt.Sprintf("%s-config", fd.Name),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "fluentd",
							Image:           fd.Spec.Image,
							Args:            fd.Spec.Args,
							ImagePullPolicy: fd.Spec.ImagePullPolicy,
							Ports:           ports,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      SecretVolName,
									ReadOnly:  true,
									MountPath: FluentdMountPath,
								},
							},
							Resources: fd.Spec.Resources,
							Env: []corev1.EnvVar{
								{
									Name:  "BUFFER_PATH",
									Value: BufferMountPath,
								},
							},
							SecurityContext: fd.Spec.ContainerSecurityContext,
						},
					},
					NodeSelector: fd.Spec.NodeSelector,
					Tolerations:  fd.Spec.Tolerations,
					Affinity:     fd.Spec.Affinity,
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: &defaultFsGroup,
					},
				},
			},
		},
	}

	if fd.Spec.RuntimeClassName != "" {
		daemonSet.Spec.Template.Spec.RuntimeClassName = &fd.Spec.RuntimeClassName
	}

	if fd.Spec.PriorityClassName != "" {
		daemonSet.Spec.Template.Spec.PriorityClassName = fd.Spec.PriorityClassName
	}

	if fd.Spec.Volumes != nil {
		daemonSet.Spec.Template.Spec.Volumes = append(daemonSet.Spec.Template.Spec.Volumes, fd.Spec.Volumes...)
	}

	if fd.Spec.VolumeMounts != nil {
		daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts, fd.Spec.VolumeMounts...)
	}

	if fd.Spec.EnvVars != nil {
		daemonSet.Spec.Template.Spec.Containers[0].Env = append(daemonSet.Spec.Template.Spec.Containers[0].Env, fd.Spec.EnvVars...)
	}

	if fd.Spec.SecurityContext != nil {
		daemonSet.Spec.Template.Spec.SecurityContext = fd.Spec.SecurityContext
	}

	if fd.Spec.SchedulerName != "" {
		daemonSet.Spec.Template.Spec.SchedulerName = fd.Spec.SchedulerName
	}

	if fd.Spec.PositionDB != (corev1.VolumeSource{}) {
		daemonSet.Spec.Template.Spec.Volumes = append(daemonSet.Spec.Template.Spec.Volumes, corev1.Volume{
			Name:         "positions",
			VolumeSource: fd.Spec.PositionDB,
		})
		daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      "positions",
			MountPath: "/fluent-bit/tail",
		})
	}
	// Mount host or emptydir VolumeSource
	if fd.Spec.BufferVolume != nil && !fd.Spec.BufferVolume.DisableBufferVolume {
		bufferVolName := fmt.Sprintf("%s-buffer", fd.Name)
		bufferpv := fd.Spec.BufferVolume

		if bufferpv.HostPath != nil {
			daemonSet.Spec.Template.Spec.Volumes = append(daemonSet.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: bufferVolName,
				VolumeSource: corev1.VolumeSource{
					HostPath: bufferpv.HostPath,
				},
			})

			daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      bufferVolName,
				MountPath: BufferMountPath,
			})
			return &daemonSet
		}

		if bufferpv.EmptyDir != nil {
			daemonSet.Spec.Template.Spec.Volumes = append(daemonSet.Spec.Template.Spec.Volumes, corev1.Volume{
				Name: bufferVolName,
				VolumeSource: corev1.VolumeSource{
					EmptyDir: bufferpv.EmptyDir,
				},
			})

			daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
				Name:      bufferVolName,
				MountPath: BufferMountPath,
			})

			return &daemonSet
		}
	}

	if fd.Spec.BufferVolume == nil || !fd.Spec.BufferVolume.DisableBufferVolume {
		daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts = append(daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      fmt.Sprintf("%s-buffer-pvc", fd.Name),
			MountPath: BufferMountPath,
		})
	}
	return &daemonSet
}

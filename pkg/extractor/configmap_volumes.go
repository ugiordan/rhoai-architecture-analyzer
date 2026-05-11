package extractor

// ConfigMapVolumeRef represents an explicit link between a ConfigMap and the
// container that mounts it. This surfaces the ConfigMap consumption pattern
// that env-based refs (EnvFromConfigmaps) miss.
type ConfigMapVolumeRef struct {
	DeploymentName string   `json:"deployment_name"`
	ContainerName  string   `json:"container_name"`
	VolumeName     string   `json:"volume_name"`
	MountPath      string   `json:"mount_path"`
	ConfigMapName  string   `json:"configmap_name"`
	ReadOnly       bool     `json:"read_only,omitempty"`
	SubPath        string   `json:"sub_path,omitempty"`
	Source         string   `json:"source"`
}

// extractConfigMapVolumes correlates volume mounts with their ConfigMap sources
// across all deployments. This runs after deployment extraction and creates
// explicit traversable links from ConfigMap -> container -> mount path.
func extractConfigMapVolumes(arch *ComponentArchitecture) []ConfigMapVolumeRef {
	var refs []ConfigMapVolumeRef

	for _, dep := range arch.Deployments {
		allContainers := make([]struct {
			containers []Container
			initPrefix string
		}, 0, 2)
		allContainers = append(allContainers,
			struct {
				containers []Container
				initPrefix string
			}{dep.Containers, ""},
			struct {
				containers []Container
				initPrefix string
			}{dep.InitContainers, "init:"},
		)

		for _, group := range allContainers {
			for _, c := range group.containers {
				for _, vm := range c.VolumeMounts {
					cmName, _ := vm["configMap"].(string)
					if cmName == "" {
						continue
					}
					mountPath, _ := vm["mountPath"].(string)
					volName, _ := vm["name"].(string)
					readOnly, _ := vm["readOnly"].(bool)
					subPath, _ := vm["subPath"].(string)

					containerName := group.initPrefix + c.Name

					refs = append(refs, ConfigMapVolumeRef{
						DeploymentName: dep.Name,
						ContainerName:  containerName,
						VolumeName:     volName,
						MountPath:      mountPath,
						ConfigMapName:  cmName,
						ReadOnly:       readOnly,
						SubPath:        subPath,
						Source:         dep.Source,
					})
				}
			}
		}
	}

	return refs
}

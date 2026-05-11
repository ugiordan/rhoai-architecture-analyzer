package extractor

import (
	"testing"
)

func TestExtractConfigMapVolumes(t *testing.T) {
	tests := []struct {
		name     string
		arch     *ComponentArchitecture
		wantLen  int
		validate func(t *testing.T, refs []ConfigMapVolumeRef)
	}{
		{
			name: "configmap volume mount in container",
			arch: &ComponentArchitecture{
				Deployments: []Deployment{
					{
						Name:   "my-operator",
						Source: "config/manager/deployment.yaml",
						Containers: []Container{
							{
								Name:  "manager",
								Image: "controller:latest",
								VolumeMounts: []map[string]interface{}{
									{
										"name":      "config-vol",
										"mountPath": "/etc/config",
										"configMap": "operator-config",
										"readOnly":  true,
									},
								},
							},
						},
					},
				},
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []ConfigMapVolumeRef) {
				ref := refs[0]
				if ref.DeploymentName != "my-operator" {
					t.Errorf("expected deployment 'my-operator', got %q", ref.DeploymentName)
				}
				if ref.ConfigMapName != "operator-config" {
					t.Errorf("expected configmap 'operator-config', got %q", ref.ConfigMapName)
				}
				if ref.MountPath != "/etc/config" {
					t.Errorf("expected mountPath '/etc/config', got %q", ref.MountPath)
				}
				if !ref.ReadOnly {
					t.Error("expected readOnly true")
				}
			},
		},
		{
			name: "secret volume mount ignored",
			arch: &ComponentArchitecture{
				Deployments: []Deployment{
					{
						Name: "my-app",
						Containers: []Container{
							{
								Name: "app",
								VolumeMounts: []map[string]interface{}{
									{
										"name":      "secret-vol",
										"mountPath": "/etc/secret",
										"secret":    "my-secret",
									},
								},
							},
						},
					},
				},
			},
			wantLen: 0,
		},
		{
			name: "init container with configmap",
			arch: &ComponentArchitecture{
				Deployments: []Deployment{
					{
						Name:   "my-app",
						Source: "deployment.yaml",
						InitContainers: []Container{
							{
								Name: "setup",
								VolumeMounts: []map[string]interface{}{
									{
										"name":      "init-config",
										"mountPath": "/init",
										"configMap": "setup-config",
									},
								},
							},
						},
					},
				},
			},
			wantLen: 1,
			validate: func(t *testing.T, refs []ConfigMapVolumeRef) {
				if refs[0].ContainerName != "init:setup" {
					t.Errorf("expected container 'init:setup', got %q", refs[0].ContainerName)
				}
			},
		},
		{
			name: "multiple configmaps across containers",
			arch: &ComponentArchitecture{
				Deployments: []Deployment{
					{
						Name:   "complex-app",
						Source: "deploy.yaml",
						Containers: []Container{
							{
								Name: "main",
								VolumeMounts: []map[string]interface{}{
									{"name": "cfg1", "mountPath": "/cfg1", "configMap": "config-a"},
									{"name": "cfg2", "mountPath": "/cfg2", "configMap": "config-b"},
								},
							},
							{
								Name: "sidecar",
								VolumeMounts: []map[string]interface{}{
									{"name": "cfg3", "mountPath": "/cfg3", "configMap": "config-c"},
								},
							},
						},
					},
				},
			},
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refs := extractConfigMapVolumes(tt.arch)
			if len(refs) != tt.wantLen {
				t.Errorf("expected %d refs, got %d", tt.wantLen, len(refs))
			}
			if tt.validate != nil && len(refs) == tt.wantLen {
				tt.validate(t, refs)
			}
		})
	}
}

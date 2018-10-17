package kaniko

import (
	"fmt"
	"strings"
	"time"

	"github.com/covexo/devspace/pkg/devspace/builder/docker"
	"github.com/covexo/devspace/pkg/devspace/registry"

	"github.com/covexo/devspace/pkg/devspace/kubectl"
	synctool "github.com/covexo/devspace/pkg/devspace/sync"
	"github.com/covexo/devspace/pkg/util/ignoreutil"
	"github.com/covexo/devspace/pkg/util/log"
	"github.com/covexo/devspace/pkg/util/randutil"
	"github.com/docker/docker/api/types"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/util/interrupt"
)

// Builder holds the necessary information to build and push docker images
type Builder struct {
	RegistryURL      string
	ImageName        string
	ImageTag         string
	PreviousImageTag string
	BuildNamespace   string

	allowInsecureRegistry bool
	kubectl               *kubernetes.Clientset
}

// NewBuilder creates a new kaniko.Builder instance
func NewBuilder(registryURL, imageName, imageTag, lastImageTag, buildNamespace string, kubectl *kubernetes.Clientset, allowInsecureRegistry bool) (*Builder, error) {
	return &Builder{
		RegistryURL:           registryURL,
		ImageName:             imageName,
		ImageTag:              imageTag,
		PreviousImageTag:      lastImageTag,
		BuildNamespace:        buildNamespace,
		allowInsecureRegistry: allowInsecureRegistry,
		kubectl:               kubectl,
	}, nil
}

// Authenticate authenticates kaniko for pushing to the RegistryURL (if username == "", it will try to get login data from local docker daemon)
func (b *Builder) Authenticate(username, password string, checkCredentialsStore bool) (*types.AuthConfig, error) {
	email := "noreply@devspace-cloud.com"

	if len(username) == 0 {
		dockerBuilder, dockerBuilderErr := docker.NewBuilder(b.RegistryURL, b.ImageName, b.ImageTag, false)
		if dockerBuilderErr != nil {
			return nil, dockerBuilderErr
		}

		authConfig, err := dockerBuilder.Authenticate(username, password, true)
		if err != nil {
			return nil, err
		}
		username = authConfig.Username
		email = authConfig.Email

		if authConfig.Password != "" {
			password = authConfig.Password
		} else {
			password = authConfig.IdentityToken
		}
	}
	return nil, registry.CreatePullSecret(b.kubectl, b.BuildNamespace, b.RegistryURL, username, password, email)
}

// BuildImage builds a dockerimage within a kaniko pod
func (b *Builder) BuildImage(contextPath, dockerfilePath string, options *types.ImageBuildOptions) error {
	pullSecretName := registry.GetRegistryAuthSecretName(b.RegistryURL)
	randString, _ := randutil.GenerateRandomString(12)
	buildID := strings.ToLower(randString)
	buildPod := &k8sv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "devspace-build-",
			Labels: map[string]string{
				"devspace-build-id": buildID,
			},
		},
		Spec: k8sv1.PodSpec{
			Containers: []k8sv1.Container{
				{
					Name:            "kaniko",
					Image:           "gcr.io/kaniko-project/executor:debug-72e088fda562a73859371354b1fc20e6fd7adea8",
					ImagePullPolicy: k8sv1.PullIfNotPresent,
					Command: []string{
						"/busybox/sleep",
					},
					Args: []string{
						"36000",
					},
					VolumeMounts: []k8sv1.VolumeMount{
						{
							Name:      pullSecretName,
							MountPath: "/root/.docker",
						},
					},
				},
			},
			Volumes: []k8sv1.Volume{
				{
					Name: pullSecretName,
					VolumeSource: k8sv1.VolumeSource{
						Secret: &k8sv1.SecretVolumeSource{
							SecretName: pullSecretName,
							Items: []k8sv1.KeyToPath{
								{
									Key:  k8sv1.DockerConfigJsonKey,
									Path: "config.json",
								},
							},
						},
					},
				},
			},
			RestartPolicy: k8sv1.RestartPolicyOnFailure,
		},
	}

	deleteBuildPod := func() {
		gracePeriod := int64(3)

		deleteErr := b.kubectl.Core().Pods(b.BuildNamespace).Delete(buildPod.Name, &metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriod,
		})

		if deleteErr != nil {
			log.Errorf("Failed to delete build pod: %s", deleteErr.Error())
		}
	}

	intr := interrupt.New(nil, deleteBuildPod)

	err := intr.Run(func() error {
		buildPodCreated, buildPodCreateErr := b.kubectl.Core().Pods(b.BuildNamespace).Create(buildPod)

		if buildPodCreateErr != nil {
			return fmt.Errorf("Unable to create build pod: %s", buildPodCreateErr.Error())
		}

		readyWaitTime := 2 * 60 * time.Second
		readyCheckInterval := 5 * time.Second
		buildPodReady := false

		log.StartWait("Waiting for kaniko build pod to start")

		for readyWaitTime > 0 {
			buildPod, _ = b.kubectl.Core().Pods(b.BuildNamespace).Get(buildPodCreated.Name, metav1.GetOptions{})

			if len(buildPod.Status.ContainerStatuses) > 0 && buildPod.Status.ContainerStatuses[0].Ready {
				buildPodReady = true
				break
			}

			time.Sleep(readyCheckInterval)
			readyWaitTime = readyWaitTime - readyCheckInterval
		}

		log.StopWait()
		log.Done("Kaniko build pod started")

		if !buildPodReady {
			return fmt.Errorf("Unable to start build pod")
		}
		ignoreRules, ignoreRuleErr := ignoreutil.GetIgnoreRules(contextPath)

		if ignoreRuleErr != nil {
			return fmt.Errorf("Unable to parse .dockerignore files: %s", ignoreRuleErr.Error())
		}

		buildContainer := &buildPod.Spec.Containers[0]

		log.StartWait("Uploading files to build container")
		err := synctool.CopyToContainer(b.kubectl, buildPod, buildContainer, contextPath, "/src", ignoreRules)

		if err != nil {
			return fmt.Errorf("Error uploading files to container: %s", err.Error())
		}
		err = synctool.CopyToContainer(b.kubectl, buildPod, buildContainer, dockerfilePath, "/src", ignoreRules)

		if err != nil {
			return fmt.Errorf("Error uploading files to container: %s", err.Error())
		}
		log.StopWait()
		log.Done("Uploaded files to container")

		log.StartWait("Building container image")

		imageDestination := b.ImageName + ":" + b.ImageTag

		if b.RegistryURL != "" {
			imageDestination = strings.TrimSuffix(b.RegistryURL, "/") + "/" + imageDestination
		}
		containerBuildPath := "/src"
		exitChannel := make(chan error)
		kanikoBuildCmd := []string{
			"/kaniko/executor",
			"--dockerfile=" + containerBuildPath + "/Dockerfile",
			"--context=dir://" + containerBuildPath,
			"--destination=" + imageDestination,
			"--single-snapshot",
		}

		if !options.NoCache {
			kanikoBuildCmd = append(kanikoBuildCmd, "--cache=true", "--cache-repo="+b.PreviousImageTag)
		}

		if b.allowInsecureRegistry {
			kanikoBuildCmd = append(kanikoBuildCmd, "--insecure", "--skip-tls-verify")
		}

		stdin, stdout, stderr, execErr := kubectl.Exec(b.kubectl, buildPod, buildContainer.Name, kanikoBuildCmd, false, exitChannel)
		stdin.Close()

		if execErr != nil {
			return fmt.Errorf("Failed to start image building: %s", execErr.Error())
		}

		lastKanikoOutput := formatKanikoOutput(stdout, stderr)
		exitError := <-exitChannel

		log.StopWait()

		if exitError != nil {
			return fmt.Errorf("Error: %s, Last Kaniko Output: %s", exitError.Error(), lastKanikoOutput)
		}

		log.Done("Done building image")

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// PushImage is required to implement builder.Interface
func (b *Builder) PushImage() error {
	return nil
}

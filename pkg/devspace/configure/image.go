package configure

import (
	contextpkg "context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/devspace-cloud/devspace/pkg/devspace/builder/helper"
	"github.com/devspace-cloud/devspace/pkg/devspace/cloud"
	cloudconfig "github.com/devspace-cloud/devspace/pkg/devspace/cloud/config"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/configutil"
	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	v1 "github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"
	"github.com/devspace-cloud/devspace/pkg/devspace/docker"
	"github.com/devspace-cloud/devspace/pkg/devspace/registry"
	"github.com/devspace-cloud/devspace/pkg/util/log"
	"github.com/devspace-cloud/devspace/pkg/util/ptr"
	"github.com/devspace-cloud/devspace/pkg/util/survey"
	"github.com/pkg/errors"
)

// DefaultImageName is the default image name
const DefaultImageName = "devspace"
const dockerHubHostname = "hub.docker.com"

// GetImageConfigFromImageName returns an image config based on the image
func GetImageConfigFromImageName(imageName, dockerfile, context string) *latest.ImageConfig {
	// Figure out tag
	imageTag := ""
	splittedImage := strings.Split(imageName, ":")
	imageName = splittedImage[0]
	if len(splittedImage) > 1 {
		imageTag = splittedImage[1]
	} else if dockerfile == "" {
		imageTag = "latest"
	}

	retImageConfig := &latest.ImageConfig{
		Image:            &imageName,
		CreatePullSecret: ptr.Bool(true),
	}

	if imageTag != "" {
		retImageConfig.Tag = &imageTag
	}
	if dockerfile == "" {
		retImageConfig.Build = &latest.BuildConfig{
			Disabled: ptr.Bool(true),
		}
	} else {
		if dockerfile != helper.DefaultDockerfilePath {
			retImageConfig.Dockerfile = &dockerfile
		}
		if context != "" && context != helper.DefaultContextPath {
			retImageConfig.Context = &context
		}
	}

	return retImageConfig
}

// GetImageConfigFromDockerfile gets the image config based on the configured cloud provider or asks the user where to push to
func GetImageConfigFromDockerfile(config *latest.Config, dockerfile, context string, cloudProvider *string, checkRegistryAuth bool) (*latest.ImageConfig, error) {
	var (
		dockerUsername = ""
		retImageConfig = &latest.ImageConfig{}
	)

	// Get docker client
	client, err := docker.NewClient(config, true, log.GetInstance())
	if err != nil {
		return nil, fmt.Errorf("Cannot create docker client: %v", err)
	}

	// Check if docker is installed
	_, err = client.Ping(contextpkg.Background())
	if err != nil {
		// Check if docker cli is installed
		runErr := exec.Command("docker").Run()
		if runErr == nil {
			log.Warn("Docker daemon not running. Start Docker daemon to build images with Docker instead of using the kaniko fallback.")
		}
	}

	// If not kaniko get docker hub credentials
	if cloudProvider == nil {
		log.StartWait("Checking Docker credentials")
		dockerAuthConfig, err := docker.GetAuthConfig(client, "", true)
		log.StopWait()

		if err == nil {
			dockerUsername = dockerAuthConfig.Username
		}
	}

	registryURL, err := getRegistryURL(config, cloudProvider, checkRegistryAuth)
	if err != nil {
		return nil, err
	}

	// Image name to use
	defaultImageName := ""

	// Is docker hub?
	if registryURL == "hub.docker.com" {
		defaultImageName = survey.Question(&survey.QuestionOptions{
			Question:          "Which image name do you want to use on Docker Hub?",
			DefaultValue:      dockerUsername + "/devspace",
			ValidationMessage: "Please enter a valid docker image name (e.g. myregistry.com/user/repository)",
			ValidationFunc: func(name string) error {
				_, err := registry.GetStrippedDockerImageName(name)
				return err
			},
		})
		defaultImageName, _ = registry.GetStrippedDockerImageName(defaultImageName)
	} else if regexp.MustCompile("^(.+\\.)?gcr.io$").Match([]byte(registryURL)) { // Is google registry?
		project, err := exec.Command("gcloud", "config", "get-value", "project").Output()
		gcloudProject := "myGCloudProject"

		if err == nil {
			gcloudProject = strings.TrimSpace(string(project))
		}

		defaultImageName = survey.Question(&survey.QuestionOptions{
			Question:          "Which image name do you want to push to?",
			DefaultValue:      registryURL + "/" + gcloudProject + "/devspace",
			ValidationMessage: "Please enter a valid docker image name (e.g. myregistry.com/user/repository)",
			ValidationFunc: func(name string) error {
				_, err := registry.GetStrippedDockerImageName(name)
				return err
			},
		})
		defaultImageName, _ = registry.GetStrippedDockerImageName(defaultImageName)
	} else if cloudProvider != nil {
		// Is DevSpace Cloud?
		defaultImageName = registryURL + "/${DEVSPACE_USERNAME}/" + DefaultImageName
	} else {
		if dockerUsername == "" {
			dockerUsername = "myuser"
		}

		defaultImageName = survey.Question(&survey.QuestionOptions{
			Question:          "Which image name do you want to push to?",
			DefaultValue:      registryURL + "/" + dockerUsername + "/devspace",
			ValidationMessage: "Please enter a valid docker image name (e.g. myregistry.com/user/repository)",
			ValidationFunc: func(name string) error {
				_, err := registry.GetStrippedDockerImageName(name)
				return err
			},
		})
		defaultImageName, _ = registry.GetStrippedDockerImageName(defaultImageName)
	}

	// Set image name
	retImageConfig.Image = &defaultImageName

	// Set image specifics
	if dockerfile != "" && dockerfile != helper.DefaultDockerfilePath {
		retImageConfig.Dockerfile = &dockerfile
	}
	if context != "" && context != helper.DefaultContextPath {
		retImageConfig.Context = &context
	}

	retImageConfig.CreatePullSecret = ptr.Bool(true)

	return retImageConfig, nil
}

func getRegistryURL(config *latest.Config, cloudProvider *string, checkRegistryAuth bool) (string, error) {
	cloudRegistryHostname, err := getCloudRegistryHostname(cloudProvider)
	if err != nil {
		return "", err
	}

	useDockerHub := "Use " + dockerHubHostname
	useDevSpaceRegistry := "Use " + cloudRegistryHostname + " (free, private Docker registry)"
	useOtherRegistry := "Use other registry"
	registryUsernameHint := " => you are logged in as %s"
	registryDefaultOption := useDevSpaceRegistry
	registryLoginHint := "Please login via `docker login%s` and try again."

	dockerClient, err := docker.NewClient(config, false, log.GetInstance())
	if err != nil {
		return "", fmt.Errorf("Error creating docker client: %v", err)
	}

	authConfig, err := docker.GetAuthConfig(dockerClient, dockerHubHostname, true)
	if err == nil && authConfig.Username != "" {
		useDockerHub = useDockerHub + fmt.Sprintf(registryUsernameHint, authConfig.Username)
		registryDefaultOption = useDockerHub
	}

	authConfig, err = docker.GetAuthConfig(dockerClient, cloudRegistryHostname, true)
	if err == nil && authConfig.Username != "" {
		useDevSpaceRegistry = useDevSpaceRegistry + fmt.Sprintf(registryUsernameHint, authConfig.Username)
		registryDefaultOption = useDevSpaceRegistry
	}

	registryOptions := []string{useDockerHub, useDevSpaceRegistry, useOtherRegistry}

	selectedRegistry := survey.Question(&survey.QuestionOptions{
		Question:     "Which registry do you want to use for storing your Docker images?",
		DefaultValue: registryDefaultOption,
		Options:      registryOptions,
	})
	var registryURL string

	if selectedRegistry == useDockerHub {
		registryURL = dockerHubHostname
	} else if selectedRegistry == useDevSpaceRegistry {
		registryURL = cloudRegistryHostname
		registryLoginHint = fmt.Sprintf(registryLoginHint, " "+cloudRegistryHostname)
	} else {
		registryURL = survey.Question(&survey.QuestionOptions{
			Question:     "Please enter the registry URL without image name:",
			DefaultValue: "my.registry.tld/username",
		})
		registryURL = strings.Trim(registryURL, "/ ")
		registryLoginHint = fmt.Sprintf(registryLoginHint, " "+registryURL)
	}
	log.WriteString("\n")

	if checkRegistryAuth {
		if dockerClient != nil {
			log.StartWait("Checking registry authentication")
			authConfig, err := docker.Login(dockerClient, registryURL, "", "", true, false, false)
			log.StopWait()

			if err != nil || authConfig.Username == "" {
				if registryURL == dockerHubHostname {
					log.Warn("You are not logged in to Docker Hub")
					log.Warn("Please make sure you have a https://hub.docker.com account")
					log.Warn("Installing docker is NOT required. You simply need a Docker Hub account\n")

					for {
						dockerUsername := survey.Question(&survey.QuestionOptions{
							Question:               "What is your Docker Hub username?",
							DefaultValue:           "",
							ValidationRegexPattern: "^.*$",
						})

						dockerPassword := survey.Question(&survey.QuestionOptions{
							Question:               "What is your Docker Hub password? (will only be sent to Docker Hub)",
							DefaultValue:           "",
							ValidationRegexPattern: "^.*$",
							IsPassword:             true,
						})

						_, err = docker.Login(dockerClient, registryURL, dockerUsername, dockerPassword, false, true, true)
						if err != nil {
							log.Warn(err)
							continue
						}

						break
					}
				} else if selectedRegistry == useDevSpaceRegistry {
					return "", loginDevSpaceCloud(*cloudProvider)
				} else {
					return "", fmt.Errorf("Registry authentication failed for %s.\n         %s", registryURL, registryLoginHint)
				}
			}
		} else {
			return "", errors.New("Unable to verify registry authentication.\n         Please install and start Docker, or use dscr.io as registry")
		}
	}
	return registryURL, nil
}

func getCloudRegistryHostname(cloudProvider *string) (string, error) {
	var registryURL string

	if cloudProvider == nil || *cloudProvider == "" {
		// prevents EnsureLoggedIn call in GetProvider
		// TODO: remove this hard-coded exception once the registry URL can be retrieved from DevSpace Cloud without login
		registryURL = "dscr.io"
	} else {
		// Get default registry
		provider, err := cloud.GetProvider(cloudProvider, log.GetInstance())
		if err != nil {
			return "", fmt.Errorf("Error login into cloud provider: %v", err)
		}

		registries, err := provider.GetRegistries()
		if err != nil {
			return "", fmt.Errorf("Error retrieving registries: %v", err)
		}
		if len(registries) > 0 {
			registryURL = registries[0].URL
		} else {
			registryURL = "hub.docker.com"
		}
	}

	return registryURL, nil
}

func loginDevSpaceCloud(cloudProvider string) error {
	// Get provider configuration
	providerConfig, err := cloudconfig.ParseProviderConfig()
	if err != nil {
		return err
	}

	// Set default cloud provider, if none is provided
	if cloudProvider == "" {
		cloudProvider = cloudconfig.DevSpaceCloudProviderName
	}

	// Choose cloud provider
	if providerConfig.Default != "" {
		cloudProvider = providerConfig.Default
	} else if len(providerConfig.Providers) > 1 {
		options := []string{}
		for _, provider := range providerConfig.Providers {
			options = append(options, provider.Name)
		}

		cloudProvider = survey.Question(&survey.QuestionOptions{
			Question: "Select a cloud provider",
			Options:  options,
		})
	}

	// Ensure user is logged in
	return cloud.EnsureLoggedIn(providerConfig, cloudProvider, log.GetInstance())
}

// AddImage adds an image to the devspace
func AddImage(nameInConfig, name, tag, contextPath, dockerfilePath, buildEngine string) error {
	config := configutil.GetBaseConfig()

	imageConfig := &v1.ImageConfig{
		Image: &name,
	}

	if tag != "" {
		imageConfig.Tag = &tag
	}
	if contextPath != "" {
		imageConfig.Context = &contextPath
	}
	if dockerfilePath != "" {
		imageConfig.Dockerfile = &dockerfilePath
	}

	if buildEngine == "docker" {
		if imageConfig.Build == nil {
			imageConfig.Build = &v1.BuildConfig{}
		}

		imageConfig.Build.Docker = &v1.DockerConfig{}
	} else if buildEngine == "kaniko" {
		if imageConfig.Build == nil {
			imageConfig.Build = &v1.BuildConfig{}
		}

		imageConfig.Build.Kaniko = &v1.KanikoConfig{}
	} else if buildEngine != "" {
		log.Errorf("BuildEngine %v unknown. Please select one of docker|kaniko", buildEngine)
	}

	if config.Images == nil {
		images := make(map[string]*v1.ImageConfig)
		config.Images = &images
	}

	(*config.Images)[nameInConfig] = imageConfig

	err := configutil.SaveLoadedConfig()
	if err != nil {
		return fmt.Errorf("Couldn't save config file: %s", err.Error())
	}

	return nil
}

//RemoveImage removes an image from the devspace
func RemoveImage(removeAll bool, names []string) error {
	config := configutil.GetBaseConfig()

	if len(names) == 0 && removeAll == false {
		return fmt.Errorf("You have to specify at least one image")
	}

	newImageList := make(map[string]*v1.ImageConfig)

	if !removeAll && config.Images != nil {

	ImagesLoop:
		for nameInConfig, imageConfig := range *config.Images {
			for _, deletionName := range names {
				if deletionName == nameInConfig {
					continue ImagesLoop
				}
			}

			newImageList[nameInConfig] = imageConfig
		}
	}

	config.Images = &newImageList

	err := configutil.SaveLoadedConfig()
	if err != nil {
		return fmt.Errorf("Couldn't save config file: %v", err)
	}

	return nil
}

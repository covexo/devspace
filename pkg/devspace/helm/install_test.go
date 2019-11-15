package helm

import ()

/*type checkDependenciesTestCase struct {
	name string

	dependenciesInChart        []*chart.Chart
	dependenciesInRequirements []*helmchartutil.Dependency

	expectedErr string
}

func TestCheckDependencies(t *testing.T) {
	testCases := []checkDependenciesTestCase{
		checkDependenciesTestCase{
			name:                       "Matching dependencies in chart and requirements",
			dependenciesInChart:        []*chart.Chart{&chart.Chart{Metadata: &chart.Metadata{Name: "MatchingDep"}}},
			dependenciesInRequirements: []*helmchartutil.Dependency{&helmchartutil.Dependency{Name: "MatchingDep"}},
		},
		checkDependenciesTestCase{
			name:                       "Requirements has dependency and that the chart has not",
			dependenciesInChart:        []*chart.Chart{&chart.Chart{Metadata: &chart.Metadata{Name: "ChartDep"}}},
			dependenciesInRequirements: []*helmchartutil.Dependency{&helmchartutil.Dependency{Name: "ReqDep"}},
			expectedErr:                "found in requirements.yaml, but missing in charts/ directory: ReqDep",
		},
	}

	for _, testCase := range testCases {
		ch := &chart.Chart{
			Dependencies: testCase.dependenciesInChart,
		}
		reqs := &helmchartutil.Requirements{
			Dependencies: testCase.dependenciesInRequirements,
		}

		err := checkDependencies(ch, reqs)

		if testCase.expectedErr == "" {
			assert.NilError(t, err, "Error checking dependencies in testCase %s", testCase.name)
		} else {
			assert.Error(t, err, testCase.expectedErr, "Wrong or no error checking dependencies in testCase %s", testCase.name)
		}
	}
}

func TestInstallChart(t *testing.T) {
	t.Skip("You're too slow")
	config := createFakeConfig()

	// Create the fake client.
	kubeClient := &kubectl.Client{
		Client: fake.NewSimpleClientset(),
	}
	helmClient := &helm.FakeClient{}

	client, err := create(config, configutil.TestNamespace, helmClient, kubeClient, log.GetInstance())
	if err != nil {
		t.Fatal(err)
	}

	helmConfig := &latest.HelmConfig{
		Chart: &latest.ChartConfig{
			Name: "stable/nginx-ingress",
		},
	}

	err = client.UpdateRepos(log.GetInstance())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.InstallChart("my-release", "", &map[interface{}]interface{}{}, helmConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Upgrade
	_, err = client.InstallChart("my-release", "", &map[interface{}]interface{}{}, helmConfig)
	if err != nil {
		t.Fatal(err)
	}
}

type analyzeErrorTestCase struct {
	name string

	inputErr    error
	namespace   string
	createdPods []*k8sv1.Pod

	expectedErr string
}

func TestAnalyzeError(t *testing.T) {
	testCases := []analyzeErrorTestCase{
		analyzeErrorTestCase{
			name:        "Test analyze no-timeout error",
			inputErr:    errors.Errorf("Some error"),
			expectedErr: "Some error",
		},
		analyzeErrorTestCase{
			name:      "Test analyze timeout error",
			inputErr:  errors.Errorf("timed out waiting"),
			namespace: "testNS",
		},
	}

	for _, testCase := range testCases {
		config := createFakeConfig()

		// Create the fake client.
		kubeClient := &kubectl.Client{
			Client: fake.NewSimpleClientset(),
		}
		helmClient := &helm.FakeClient{}

		for _, pod := range testCase.createdPods {
			_, err := kubeClient.Client.CoreV1().Pods(testCase.namespace).Create(pod)
			assert.NilError(t, err, "Error creating testPod in testCase %s", testCase.name)
		}

		client, err := create(config, configutil.TestNamespace, helmClient, kubeClient, &log.DiscardLogger{})
		if err != nil {
			t.Fatal(err)
		}

		err = client.analyzeError(testCase.inputErr, testCase.namespace)
		if testCase.expectedErr == "" {
			assert.NilError(t, err, "Error analyzing error in testCase %s", testCase.name)
		} else {
			assert.Error(t, err, testCase.expectedErr, "Wrong or no error returned in testCase %s", testCase.name)
		}
	}
}*/

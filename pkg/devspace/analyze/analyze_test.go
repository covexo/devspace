package analyze

import ()

/*func TestAnalyze(t *testing.T) {
	kubeClient := &kubectl.Client{
		Client: fake.NewSimpleClientset(),
	}

	//Analyze empty
	err := Analyze(kubeClient, "testNS", true, &log.DiscardLogger{})
	assert.NilError(t, err, "Error while analyzing")

}

func TestCreateReport(t *testing.T) {
	kubeClient := &kubectl.Client{
		Client: fake.NewSimpleClientset(),
	}

	_, err := kubeClient.Client.CoreV1().Namespaces().Create(&k8sv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testNS",
		},
	})
	assert.NilError(t, err, "Error creating namespace")
	_, err = kubeClient.Client.CoreV1().Pods("testNS").Create(&k8sv1.Pod{
		Status: k8sv1.PodStatus{
			Reason: "Error",
		},
	})
	assert.NilError(t, err, "Error creating pod")

	reports, err := CreateReport(kubeClient, "testNS", false)
	assert.NilError(t, err, "Error while creating a report")
	assert.Equal(t, 1, len(reports), "Wrong number of problems reported")
	assert.Equal(t, true, strings.Contains(reports[0].Problems[0], "Pod"), "Report does not address pods")
	assert.Equal(t, true, strings.Contains(reports[0].Problems[0], "Error"), "Report does not address the pod status")

	_, err = kubeClient.Client.CoreV1().Pods("testNS").Update(&k8sv1.Pod{
		Status: k8sv1.PodStatus{
			Reason:    "Running",
			StartTime: &metav1.Time{Time: time.Now().Add(-MinimumPodAge / 10 * 9)},
		},
	})
	assert.NilError(t, err, "Error fixing pod")
	_, err = kubeClient.Client.AppsV1().ReplicaSets("testNS").Create(&v1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ReplicaSet with errors",
		},
		Spec: v1.ReplicaSetSpec{
			Replicas: ptr.Int32(1),
		},
		Status: v1.ReplicaSetStatus{
			Replicas: 2,
		},
	})
	assert.NilError(t, err, "Error creating replicaSet")

	reports, err = CreateReport(kubeClient, "testNS", false)
	assert.NilError(t, err, "Error while creating a report")
	assert.Equal(t, 0, len(reports), "Problems reported when only the ReplicaSets have problems.")

	err = kubeClient.Client.AppsV1().ReplicaSets("testNS").Delete("ReplicaSet with errors", &metav1.DeleteOptions{})
	assert.NilError(t, err, "Error deleting replicaSet")
	_, err = kubeClient.Client.AppsV1().StatefulSets("testNS").Create(&v1.StatefulSet{
		Spec: v1.StatefulSetSpec{
			Replicas: ptr.Int32(1),
		},
		Status: v1.StatefulSetStatus{
			Replicas:        2,
			ReadyReplicas:   2,
			CurrentReplicas: 2,
		},
	})
	assert.NilError(t, err, "Error creating statefulSet")
	reports, err = CreateReport(kubeClient, "testNS", false)
	assert.NilError(t, err, "Error while creating a report")
	assert.Equal(t, 0, len(reports), "Problems reported when only the StatefulSets have problems.")

	// Delete test namespace
	err = kubeClient.Client.CoreV1().Namespaces().Delete("testNS", &metav1.DeleteOptions{})
	if err != nil {
		t.Fatalf("Error deleting namespace: %v", err)
	}

}

func TestReportToString(t *testing.T) {
	report := []*ReportItem{
		&ReportItem{
			Name: "testReport",
			Problems: []string{
				"Somethings wrong, I guess...",
			},
		},
	}

	expectedString := `
` + ansi.Color(`  ================================================================================
                         testReport (1 potential issue(s))                        
  ================================================================================
`, "green+b")
	expectedString = expectedString + `Somethings wrong, I guess...
`
	assert.Equal(t, expectedString, ReportToString(report), "Report wrong translated")
}*/

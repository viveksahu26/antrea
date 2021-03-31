// Copyright 2021 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e

import (
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1a1 "github.com/vmware-tanzu/antrea/pkg/apis/core/v1alpha2"
	secv1alpha1 "github.com/vmware-tanzu/antrea/pkg/apis/security/v1alpha1"
)

func testInvalidCGIPBlockWithPodSelector(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup created with ipblock and podSelector")
	cgName := "ipb-pod"
	pSel := &metav1.LabelSelector{MatchLabels: map[string]string{"pod": "x"}}
	cidr := "10.0.0.10/32"
	ipb := &secv1alpha1.IPBlock{CIDR: cidr}
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: cgName,
		},
		Spec: corev1a1.GroupSpec{
			PodSelector: pSel,
			IPBlock:     ipb,
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
}

func testInvalidCGIPBlockWithNSSelector(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup created with ipblock and namespaceSelector")
	cgName := "ipb-ns"
	nSel := &metav1.LabelSelector{MatchLabels: map[string]string{"ns": "y"}}
	cidr := "10.0.0.10/32"
	ipb := &secv1alpha1.IPBlock{CIDR: cidr}
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: cgName,
		},
		Spec: corev1a1.GroupSpec{
			NamespaceSelector: nSel,
			IPBlock:           ipb,
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
}

func testInvalidCGServiceRefWithPodSelector(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup created with serviceReference and podSelector")
	cgName := "svcref-pod-selector"
	pSel := &metav1.LabelSelector{MatchLabels: map[string]string{"pod": "x"}}
	svcRef := &corev1a1.ServiceReference{
		Namespace: "y",
		Name:      "test-svc",
	}
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: cgName,
		},
		Spec: corev1a1.GroupSpec{
			PodSelector:      pSel,
			ServiceReference: svcRef,
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
}

func testInvalidCGServiceRefWithNSSelector(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup created with serviceReference and namespaceSelector")
	cgName := "svcref-ns-selector"
	nSel := &metav1.LabelSelector{MatchLabels: map[string]string{"ns": "y"}}
	svcRef := &corev1a1.ServiceReference{
		Namespace: "y",
		Name:      "test-svc",
	}
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: cgName,
		},
		Spec: corev1a1.GroupSpec{
			NamespaceSelector: nSel,
			ServiceReference:  svcRef,
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
}

func testInvalidCGServiceRefWithIPBlock(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup created with ipblock and namespaceSelector")
	cgName := "ipb-svcref"
	cidr := "10.0.0.10/32"
	ipb := &secv1alpha1.IPBlock{CIDR: cidr}
	svcRef := &corev1a1.ServiceReference{
		Namespace: "y",
		Name:      "test-svc",
	}
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: cgName,
		},
		Spec: corev1a1.GroupSpec{
			ServiceReference: svcRef,
			IPBlock:          ipb,
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
}

func testInvalidCGChildGroupDoesNotExist(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup childGroup does not exist")
	cgName := "child-group-not-exist"
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: cgName,
		},
		Spec: corev1a1.GroupSpec{
			ChildGroups: []corev1a1.ClusterGroupReference{corev1a1.ClusterGroupReference("some-non-existing-cg")},
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
}

var testChildCGName = "test-child-cg"

func createChildCGForTest(t *testing.T) {
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: testChildCGName,
		},
		Spec: corev1a1.GroupSpec{
			PodSelector: &metav1.LabelSelector{},
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err != nil {
		failOnError(err, t)
	}
}

func cleanupChildCGForTest(t *testing.T) {
	if err := k8sUtils.DeleteCG(testChildCGName); err != nil {
		failOnError(err, t)
	}
}

func testInvalidCGChildGroupWithPodSelector(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup created with childGroups and podSelector")
	cgName := "child-group-pod-selector"
	pSel := &metav1.LabelSelector{MatchLabels: map[string]string{"pod": "x"}}
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: cgName,
		},
		Spec: corev1a1.GroupSpec{
			PodSelector: pSel,
			ChildGroups: []corev1a1.ClusterGroupReference{corev1a1.ClusterGroupReference(testChildCGName)},
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
}

func testInvalidCGChildGroupWithServiceReference(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup created with childGroups and ServiceReference")
	cgName := "child-group-svcref"
	svcRef := &corev1a1.ServiceReference{
		Namespace: "y",
		Name:      "test-svc",
	}
	cg := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name: cgName,
		},
		Spec: corev1a1.GroupSpec{
			ServiceReference: svcRef,
			ChildGroups:      []corev1a1.ClusterGroupReference{corev1a1.ClusterGroupReference(testChildCGName)},
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
}

func testInvalidCGMaxNestedLevel(t *testing.T) {
	invalidErr := fmt.Errorf("clustergroup created with childGroup which has childGroups itself")
	cgName1, cgName2 := "cg-nested-1", "cg-nested-2"
	cg1 := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{Name: cgName1},
		Spec: corev1a1.GroupSpec{
			ChildGroups: []corev1a1.ClusterGroupReference{corev1a1.ClusterGroupReference(testChildCGName)},
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg1); err != nil {
		// Above creation of CG must succeed as it is a valid spec.
		failOnError(err, t)
	}
	cg2 := &corev1a1.ClusterGroup{
		ObjectMeta: metav1.ObjectMeta{Name: cgName2},
		Spec: corev1a1.GroupSpec{
			ChildGroups: []corev1a1.ClusterGroupReference{corev1a1.ClusterGroupReference(cgName1)},
		},
	}
	if _, err := k8sUtils.CreateOrUpdateCG(cg2); err == nil {
		// Above creation of CG must fail as it is an invalid spec.
		failOnError(invalidErr, t)
	}
	// cleanup cg-nested-1
	if err := k8sUtils.DeleteCG(cgName1); err != nil {
		failOnError(err, t)
	}
}

func TestClusterGroup(t *testing.T) {
	data, err := setupTest(t)
	if err != nil {
		t.Fatalf("Error when setting up test: %v", err)
	}
	defer teardownTest(t, data)
	skipIfAntreaPolicyDisabled(t, data)
	initialize(t, data)

	t.Run("TestGroupClusterGroupValidate", func(t *testing.T) {
		t.Run("Case=IPBlockWithPodSelectorDenied", func(t *testing.T) { testInvalidCGIPBlockWithPodSelector(t) })
		t.Run("Case=IPBlockWithNamespaceSelectorDenied", func(t *testing.T) { testInvalidCGIPBlockWithNSSelector(t) })
		t.Run("Case=ServiceRefWithPodSelectorDenied", func(t *testing.T) { testInvalidCGServiceRefWithPodSelector(t) })
		t.Run("Case=ServiceRefWithNamespaceSelectorDenied", func(t *testing.T) { testInvalidCGServiceRefWithNSSelector(t) })
		t.Run("Case=ServiceRefWithIPBlockDenied", func(t *testing.T) { testInvalidCGServiceRefWithIPBlock(t) })
		t.Run("Case=InvalidChildGroupName", func(t *testing.T) { testInvalidCGChildGroupDoesNotExist(t) })
	})
	t.Run("TestGroupClusterGroupValidateChildGroup", func(t *testing.T) {
		createChildCGForTest(t)
		t.Run("Case=ChildGroupWithPodSelectorDenied", func(t *testing.T) { testInvalidCGChildGroupWithPodSelector(t) })
		t.Run("Case=ChildGroupWithPodServiceReferenceDenied", func(t *testing.T) { testInvalidCGChildGroupWithServiceReference(t) })
		t.Run("Case=ChildGroupExceedMaxNestedLevel", func(t *testing.T) { testInvalidCGMaxNestedLevel(t) })
		cleanupChildCGForTest(t)
	})
	failOnError(k8sUtils.CleanCGs(), t)
}
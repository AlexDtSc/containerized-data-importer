package tests_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/ginkgo/extensions/table"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"kubevirt.io/containerized-data-importer/tests/framework"
	"kubevirt.io/containerized-data-importer/tests/utils"

	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/datavolumecontroller/v1alpha1"
)

var _ = Describe("DataVolume tests", func() {

	f, err := framework.NewFramework("dv-func-test", framework.Config{})
	if err != nil {
		Fail("Unable to create framework struct")
	}

	Describe("Verify DataVolume", func() {
		table.DescribeTable("with http import source should", func(url string, phase cdiv1.DataVolumePhase, dataVolumeName string) {
			dataVolume := utils.NewDataVolumeWithHttpImport(dataVolumeName, "1Gi", url)

			By(fmt.Sprintf("creating new datavolume %s", dataVolume.Name))
			dataVolume, err := utils.CreateDataVolumeFromDefinition(f.CdiClient, f.Namespace.Name, dataVolume)
			Expect(err).To(BeNil())

			By(fmt.Sprintf("waiting for datavolume to match phase %s", string(phase)))
			utils.WaitForDataVolumePhase(f.CdiClient, f.Namespace.Name, phase, dataVolume.Name)

			// verify PVC was created
			By("verifying pvc was created")
			_, err = f.K8sClient.CoreV1().PersistentVolumeClaims(dataVolume.Namespace).Get(dataVolume.Name, metav1.GetOptions{})
			Expect(err).To(BeNil())

			err = utils.DeleteDataVolume(f.CdiClient, f.Namespace.Name, dataVolume)
			Expect(err).To(BeNil())

		},
			table.Entry("succeed when given valid url", utils.TinyCoreIsoURL, cdiv1.Succeeded, "dv-phase-test-1"),
			table.Entry("fail due to invalid DNS entry", "http://i-made-this-up.kube-system/tinyCore.iso", cdiv1.Failed, "dv-phase-test-2"),
			table.Entry("fail due to file not found", utils.TinyCoreIsoURL+"not.real.file", cdiv1.Failed, "dv-phase-test-3"),
		)
	})
})

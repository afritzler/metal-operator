/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"

	metalv1alpha1 "github.com/afritzler/metal-operator/api/v1alpha1"
)

var _ = Describe("ServerClaim Controller", func() {
	ns := SetupTest()

	var server *metalv1alpha1.Server

	BeforeEach(func(ctx SpecContext) {
		By("Creating an Endpoints object")
		endpoint := &metalv1alpha1.Endpoint{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-",
			},
			Spec: metalv1alpha1.EndpointSpec{
				// emulator BMC mac address
				MACAddress: "23:11:8A:33:CF:EA",
				IP:         metalv1alpha1.MustParseIP("127.0.0.1"),
			},
		}
		Expect(k8sClient.Create(ctx, endpoint)).To(Succeed())
		DeferCleanup(k8sClient.Delete, endpoint)

		By("Ensuring that the BMC resource has been created for an endpoint")
		bmc := &metalv1alpha1.BMC{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("bmc-%s", endpoint.Name),
			},
		}
		Eventually(Get(bmc)).Should(Succeed())

		By("Creating a Server object")
		By("Ensuring that the Server resource has been created")
		server = &metalv1alpha1.Server{
			ObjectMeta: metav1.ObjectMeta{
				Name: GetServerNameFromBMCandIndex(0, bmc),
			},
		}
		Eventually(Get(server)).Should(Succeed())
	})

	It("should successfully claim a server in available state", func(ctx SpecContext) {
		By("Creating an Ignition secret")
		ignitionSecret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-",
			},
			Data: map[string][]byte{
				"foo": []byte("bar"),
			},
		}
		Expect(k8sClient.Create(ctx, ignitionSecret)).To(Succeed())
		DeferCleanup(k8sClient.Delete, ignitionSecret)

		By("Creating a ServerClaim")
		claim := &metalv1alpha1.ServerClaim{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-",
			},
			Spec: metalv1alpha1.ServerClaimSpec{
				Power:             metalv1alpha1.PowerOn,
				ServerRef:         &v1.LocalObjectReference{Name: server.Name},
				IgnitionSecretRef: &v1.LocalObjectReference{Name: ignitionSecret.Name},
				Image:             "foo:bar",
			},
		}
		Expect(k8sClient.Create(ctx, claim)).To(Succeed())

		By("Patching the Server to available state")
		Eventually(UpdateStatus(server, func() {
			server.Status.State = metalv1alpha1.ServerStateAvailable
		})).Should(Succeed())

		By("Ensuring that the Server has the correct claim ref")
		Eventually(Object(server)).Should(SatisfyAll(
			HaveField("Spec.ServerClaimRef.Name", claim.Name),
			HaveField("Spec.Power", metalv1alpha1.PowerOn),
			HaveField("Status.State", metalv1alpha1.ServerStateReserved),
		))

		By("Ensuring that the ServerClaim is bound")
		Eventually(Object(claim)).Should(SatisfyAll(
			HaveField("Finalizers", ContainElement(ServerClaimFinalizer)),
			HaveField("Status.Phase", metalv1alpha1.PhaseBound),
		))

		By("Ensuring that the ServerBootConfiguration has been created")
		config := &metalv1alpha1.ServerBootConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns.Name,
				Name:      claim.Name,
			},
		}
		Eventually(Object(config)).Should(SatisfyAll(
			HaveField("OwnerReferences", ContainElement(metav1.OwnerReference{
				APIVersion:         "metal.ironcore.dev/v1alpha1",
				Kind:               "ServerClaim",
				Name:               claim.Name,
				UID:                claim.UID,
				Controller:         ptr.To(true),
				BlockOwnerDeletion: ptr.To(true),
			})),
			HaveField("Spec.ServerRef.Name", server.Name),
			HaveField("Spec.Image", "foo:bar"),
			HaveField("Spec.IgnitionSecretRef.Name", ignitionSecret.Name),
		))

		By("Ensuring that the server has a correct boot configuration ref")
		Eventually(Object(server)).Should(SatisfyAll(
			HaveField("Spec.BootConfigurationRef", &v1.ObjectReference{
				APIVersion: "metal.ironcore.dev/v1alpha1",
				Kind:       "ServerBootConfiguration",
				Namespace:  ns.Name,
				Name:       config.Name,
				UID:        config.UID,
			}),
		))

		By("Deleting the ServerClaim")
		Expect(k8sClient.Delete(ctx, claim)).To(Succeed())

		By("Ensuring that the Server is available")
		Eventually(Object(server)).Should(SatisfyAll(
			HaveField("Spec.ServerClaimRef", BeNil()),
			HaveField("Spec.BootConfigurationRef", BeNil()),
			HaveField("Spec.Power", metalv1alpha1.PowerOff),
			HaveField("Status.State", metalv1alpha1.ServerStateAvailable),
		))
	})

	It("should not claim a server in a non-available state", func(ctx SpecContext) {
		By("Creating a ServerClaim")
		claim := &metalv1alpha1.ServerClaim{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-",
			},
			Spec: metalv1alpha1.ServerClaimSpec{
				Power:     metalv1alpha1.PowerOn,
				ServerRef: &v1.LocalObjectReference{Name: server.Name},
				Image:     "foo:bar",
			},
		}
		Expect(k8sClient.Create(ctx, claim)).To(Succeed())
		DeferCleanup(k8sClient.Delete, claim)

		By("Patching the Server to available state")
		Eventually(UpdateStatus(server, func() {
			server.Status.State = metalv1alpha1.ServerStateInitial
		})).Should(Succeed())

		By("Ensuring that the Server has no claim ref")
		Eventually(Object(server)).Should(SatisfyAll(
			HaveField("Spec.ServerClaimRef", BeNil()),
			HaveField("Status.State", metalv1alpha1.ServerStateInitial),
		))

		By("Ensuring that the ServerClaim is bound")
		Eventually(Object(claim)).Should(SatisfyAll(
			HaveField("Finalizers", ContainElement(ServerClaimFinalizer)),
			HaveField("Status.Phase", metalv1alpha1.PhaseUnbound),
		))

		By("Ensuring that the ServerBootConfiguration has not been created")
		config := &metalv1alpha1.ServerBootConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns.Name,
				Name:      claim.Name,
			},
		}
		Eventually(Get(config)).Should(Satisfy(apierrors.IsNotFound))
	})

	It("should not claim a server with set claim ref", func(ctx SpecContext) {
		By("Patching the Server to available state")
		Eventually(Update(server, func() {
			server.Spec.ServerClaimRef = &v1.ObjectReference{
				APIVersion: "metal.ironcore.dev/v1alpha1",
				Kind:       "ServerClaim",
				Namespace:  ns.Name,
				Name:       "foo",
				UID:        "12345",
			}
		})).Should(Succeed())

		By("Creating a ServerClaim")
		claim := &metalv1alpha1.ServerClaim{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    ns.Name,
				GenerateName: "test-",
			},
			Spec: metalv1alpha1.ServerClaimSpec{
				Power:     metalv1alpha1.PowerOn,
				ServerRef: &v1.LocalObjectReference{Name: server.Name},
				Image:     "foo:bar",
			},
		}
		Expect(k8sClient.Create(ctx, claim)).To(Succeed())
		DeferCleanup(k8sClient.Delete, claim)

		By("Ensuring that the Server has no claim ref")
		Eventually(Object(server)).Should(SatisfyAll(
			HaveField("Spec.ServerClaimRef", &v1.ObjectReference{
				APIVersion: "metal.ironcore.dev/v1alpha1",
				Kind:       "ServerClaim",
				Namespace:  ns.Name,
				Name:       "foo",
				UID:        "12345",
			}),
			HaveField("Status.State", metalv1alpha1.ServerStateReserved),
		))

		By("Ensuring that the ServerClaim is bound")
		Eventually(Object(claim)).Should(SatisfyAll(
			HaveField("Finalizers", ContainElement(ServerClaimFinalizer)),
			HaveField("Status.Phase", metalv1alpha1.PhaseUnbound),
		))

		By("Ensuring that the ServerBootConfiguration has not been created")
		config := &metalv1alpha1.ServerBootConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns.Name,
				Name:      claim.Name,
			},
		}
		Eventually(Get(config)).Should(Satisfy(apierrors.IsNotFound))
	})
})

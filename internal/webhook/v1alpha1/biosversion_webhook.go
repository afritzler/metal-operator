// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	metalv1alpha1 "github.com/ironcore-dev/metal-operator/api/v1alpha1"
	metalutil "github.com/ironcore-dev/metal-operator/internal/util"
)

// nolint:unused
// log is for logging in this package.
var versionLog = logf.Log.WithName("biosversion-resource")

// SetupBIOSVersionWebhookWithManager registers the webhook for BIOSVersion in the manager.
func SetupBIOSVersionWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &metalv1alpha1.BIOSVersion{}).
<<<<<<< HEAD
		WithValidator(&BIOSVersionValidator{}).
=======
		WithValidator(&BIOSVersionCustomValidator{Client: mgr.GetClient()}).
>>>>>>> tmp-original-15-09-26-02-54
		Complete()
}

// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-metal-ironcore-dev-v1alpha1-biosversion,mutating=false,failurePolicy=fail,sideEffects=None,groups=metal.ironcore.dev,resources=biosversions,verbs=create;update;delete,versions=v1alpha1,name=vbiosversion-v1alpha1.kb.io,admissionReviewVersions=v1

// BIOSVersionValidator struct is responsible for validating the BIOSVersion resource
// when it is created, updated, or deleted.
<<<<<<< HEAD
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type BIOSVersionValidator struct {
	// TODO(user): Add more fields as needed for validation
}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type BIOSVersion.
func (v *BIOSVersionValidator) ValidateCreate(_ context.Context, obj *metalv1alpha1.BIOSVersion) (admission.Warnings, error) {
	biosversionlog.Info("Validation for BIOSVersion upon creation", "name", obj.GetName())

	// TODO(user): fill in your validation logic upon object creation.
=======
type BIOSVersionCustomValidator struct {
	client.Client
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type BIOSVersion.
func (v *BIOSVersionCustomValidator) ValidateCreate(ctx context.Context, obj *metalv1alpha1.BIOSVersion) (admission.Warnings, error) {
	versionLog.Info("Validation for BIOSVersion upon creation", "name", obj.GetName())
>>>>>>> tmp-original-15-09-26-02-54

	versions := &metalv1alpha1.BIOSVersionList{}
	if err := v.List(ctx, versions); err != nil {
		return nil, fmt.Errorf("failed to list BIOSVersion: %w", err)
	}
	if err := checkForDuplicateBIOSVersionRefToServer(versions, obj); err != nil {
		return nil, err
	}
	return nil, nil
}

<<<<<<< HEAD
// ValidateUpdate implements admission.Validator so a webhook will be registered for the type BIOSVersion.
func (v *BIOSVersionValidator) ValidateUpdate(_ context.Context, oldObj, newObj *metalv1alpha1.BIOSVersion) (admission.Warnings, error) {
	biosversionlog.Info("Validation for BIOSVersion upon update", "name", newObj.GetName())
=======
// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type BIOSVersion.
func (v *BIOSVersionCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj *metalv1alpha1.BIOSVersion) (admission.Warnings, error) {
	versionLog.Info("Validation for BIOSVersion upon update", "name", newObj.GetName())
>>>>>>> tmp-original-15-09-26-02-54

	// Block updates while the referenced ServerMaintenance is InMaintenance.
	if !ShouldAllowForceUpdateInProgress(newObj) && oldObj.Spec.ServerMaintenanceRef != nil {
		active, err := metalutil.IsAnyServerMaintenanceActive(ctx, v.Client, []metalv1alpha1.ObjectReference{*oldObj.Spec.ServerMaintenanceRef})
		if err != nil {
			return nil, fmt.Errorf("failed to check maintenance state: %w", err)
		}
		if active {
			msg := fmt.Errorf("BIOSVersion %s is under active maintenance, unable to update", oldObj.Name)
			return nil, apierrors.NewInvalid(
				schema.GroupKind{Group: newObj.GroupVersionKind().Group, Kind: newObj.Kind},
				newObj.GetName(), field.ErrorList{field.Forbidden(field.NewPath("spec"), msg.Error())})
		}
	}

	versions := &metalv1alpha1.BIOSVersionList{}
	if err := v.List(ctx, versions); err != nil {
		return nil, fmt.Errorf("failed to list BIOSVersion: %w", err)
	}

	if err := checkForDuplicateBIOSVersionRefToServer(versions, newObj); err != nil {
		return nil, err
	}
	return nil, nil
}

<<<<<<< HEAD
// ValidateDelete implements admission.Validator so a webhook will be registered for the type BIOSVersion.
func (v *BIOSVersionValidator) ValidateDelete(_ context.Context, obj *metalv1alpha1.BIOSVersion) (admission.Warnings, error) {
	biosversionlog.Info("Validation for BIOSVersion upon deletion", "name", obj.GetName())

	// TODO(user): fill in your validation logic upon object deletion.
=======
// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type BIOSVersion.
func (v *BIOSVersionCustomValidator) ValidateDelete(ctx context.Context, obj *metalv1alpha1.BIOSVersion) (admission.Warnings, error) {
	versionLog.Info("Validation for BIOSVersion upon deletion", "name", obj.GetName())
>>>>>>> tmp-original-15-09-26-02-54

	// Block deletion while the referenced ServerMaintenance is InMaintenance.
	if !ShouldAllowForceDeleteInProgress(obj) && obj.Spec.ServerMaintenanceRef != nil {
		active, err := metalutil.IsAnyServerMaintenanceActive(ctx, v.Client, []metalv1alpha1.ObjectReference{*obj.Spec.ServerMaintenanceRef})
		if err != nil {
			return nil, fmt.Errorf("failed to check maintenance state: %w", err)
		}
		if active {
			return nil, apierrors.NewBadRequest("BIOSVersion is under active maintenance, unable to delete")
		}
	}
	return nil, nil
}

func checkForDuplicateBIOSVersionRefToServer(versions *metalv1alpha1.BIOSVersionList, version *metalv1alpha1.BIOSVersion) error {
	if version.Spec.ServerRef == nil {
		return nil
	}

	for _, bv := range versions.Items {
		if version.Name == bv.Name {
			continue
		}
		if bv.Spec.ServerRef == nil {
			continue
		}
		if version.Spec.ServerRef.Name == bv.Spec.ServerRef.Name {
			fldErr := field.Duplicate(field.NewPath("spec").Child("serverRef").Child("name"), version.Spec.ServerRef.Name)
			fldErr.Detail = fmt.Sprintf("server (%s) referred in %s is duplicate of server (%s) referred in %s",
				version.Spec.ServerRef.Name,
				version.Name,
				bv.Spec.ServerRef.Name,
				bv.Name,
			)
			return apierrors.NewInvalid(
				schema.GroupKind{Group: version.GroupVersionKind().Group, Kind: version.Kind},
				version.GetName(), field.ErrorList{fldErr})
		}
	}
	return nil
}

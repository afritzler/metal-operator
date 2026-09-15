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

// log is for logging in this package.
var settingsLog = logf.Log.WithName("biossettings-resource")

// SetupBIOSSettingsWebhookWithManager registers the webhook for BIOSSettings in the manager.
func SetupBIOSSettingsWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &metalv1alpha1.BIOSSettings{}).
<<<<<<< HEAD
		WithValidator(&BIOSSettingsValidator{}).
=======
		WithValidator(&BIOSSettingsCustomValidator{Client: mgr.GetClient()}).
>>>>>>> tmp-original-15-09-26-02-54
		Complete()
}

// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-metal-ironcore-dev-v1alpha1-biossettings,mutating=false,failurePolicy=fail,sideEffects=None,groups=metal.ironcore.dev,resources=biossettings,verbs=create;update;delete,versions=v1alpha1,name=vbiossettings-v1alpha1.kb.io,admissionReviewVersions=v1

// BIOSSettingsValidator struct is responsible for validating the BIOSSettings resource
// when it is created, updated, or deleted.
<<<<<<< HEAD
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type BIOSSettingsValidator struct {
	// TODO(user): Add more fields as needed for validation
}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type BIOSSettings.
func (v *BIOSSettingsValidator) ValidateCreate(_ context.Context, obj *metalv1alpha1.BIOSSettings) (admission.Warnings, error) {
	biossettingslog.Info("Validation for BIOSSettings upon creation", "name", obj.GetName())
=======
type BIOSSettingsCustomValidator struct {
	client.Client
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type BIOSSettings.
func (v *BIOSSettingsCustomValidator) ValidateCreate(ctx context.Context, obj *metalv1alpha1.BIOSSettings) (admission.Warnings, error) {
	settingsLog.Info("Validation for BIOSSettings upon creation", "name", obj.GetName())
>>>>>>> tmp-original-15-09-26-02-54

	settingsList := &metalv1alpha1.BIOSSettingsList{}
	if err := v.List(ctx, settingsList); err != nil {
		return nil, fmt.Errorf("failed to list BIOSSettings: %w", err)
	}

	if err := checkForDuplicateBIOSSettingsRefToServer(settingsList, obj); err != nil {
		return nil, err
	}
	return nil, nil
}

<<<<<<< HEAD
// ValidateUpdate implements admission.Validator so a webhook will be registered for the type BIOSSettings.
func (v *BIOSSettingsValidator) ValidateUpdate(_ context.Context, oldObj, newObj *metalv1alpha1.BIOSSettings) (admission.Warnings, error) {
	biossettingslog.Info("Validation for BIOSSettings upon update", "name", newObj.GetName())
=======
// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type BIOSSettings.
func (v *BIOSSettingsCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj *metalv1alpha1.BIOSSettings) (admission.Warnings, error) {
	settingsLog.Info("Validation for BIOSSettings upon update", "name", newObj.GetName())
>>>>>>> tmp-original-15-09-26-02-54

	// Block updates while the referenced ServerMaintenance is InMaintenance.
	if !ShouldAllowForceUpdateInProgress(newObj) && oldObj.Spec.ServerMaintenanceRef != nil {
		active, err := metalutil.IsAnyServerMaintenanceActive(ctx, v.Client, []metalv1alpha1.ObjectReference{*oldObj.Spec.ServerMaintenanceRef})
		if err != nil {
			return nil, fmt.Errorf("failed to check maintenance state: %w", err)
		}
		if active {
			msg := fmt.Errorf("BIOSSettings %s is under active maintenance, unable to update", oldObj.Name)
			return nil, apierrors.NewInvalid(
				schema.GroupKind{Group: newObj.GroupVersionKind().Group, Kind: newObj.Kind},
				newObj.GetName(), field.ErrorList{field.Forbidden(field.NewPath("spec"), msg.Error())})
		}
	}

	settingsList := &metalv1alpha1.BIOSSettingsList{}
	if err := v.List(ctx, settingsList); err != nil {
		return nil, fmt.Errorf("failed to list BIOSSettings: %w", err)
	}

	if err := checkForDuplicateBIOSSettingsRefToServer(settingsList, newObj); err != nil {
		return nil, err
	}
	return nil, nil
}

<<<<<<< HEAD
// ValidateDelete implements admission.Validator so a webhook will be registered for the type BIOSSettings.
func (v *BIOSSettingsValidator) ValidateDelete(_ context.Context, obj *metalv1alpha1.BIOSSettings) (admission.Warnings, error) {
	biossettingslog.Info("Validation for BIOSSettings upon deletion", "name", obj.GetName())
=======
// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type BIOSSettings.
func (v *BIOSSettingsCustomValidator) ValidateDelete(ctx context.Context, obj *metalv1alpha1.BIOSSettings) (admission.Warnings, error) {
	settingsLog.Info("Validation for BIOSSettings upon deletion", "name", obj.GetName())
>>>>>>> tmp-original-15-09-26-02-54

	// Block deletion while the referenced ServerMaintenance is InMaintenance.
	if !ShouldAllowForceDeleteInProgress(obj) && obj.Spec.ServerMaintenanceRef != nil {
		active, err := metalutil.IsAnyServerMaintenanceActive(ctx, v.Client, []metalv1alpha1.ObjectReference{*obj.Spec.ServerMaintenanceRef})
		if err != nil {
			return nil, fmt.Errorf("failed to check maintenance state: %w", err)
		}
		if active {
			return nil, apierrors.NewBadRequest("BIOSSettings is under active maintenance, unable to delete")
		}
	}

	return nil, nil
}

func checkForDuplicateBIOSSettingsRefToServer(settingsList *metalv1alpha1.BIOSSettingsList, settings *metalv1alpha1.BIOSSettings) error {
	for _, bs := range settingsList.Items {
		if settings.Name == bs.Name {
			continue
		}
		if settings.Spec.ServerRef.Name == bs.Spec.ServerRef.Name {
			fldErr := field.Duplicate(field.NewPath("spec").Child("serverRef"), settings.Spec.ServerRef.Name)
			fldErr.Detail = fmt.Sprintf("server (%s) referred in %s is duplicate of server (%s) referred in %s",
				settings.Spec.ServerRef.Name,
				settings.Name,
				bs.Spec.ServerRef.Name,
				bs.Name)
			return apierrors.NewInvalid(
				schema.GroupKind{Group: settings.GroupVersionKind().Group, Kind: settings.Kind},
				settings.GetName(), field.ErrorList{fldErr})
		}
	}
	return nil
}

// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"context"
	"encoding/base64"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"

	metalv1alpha1 "github.com/afritzler/metal-operator/api/v1alpha1"
	"github.com/afritzler/metal-operator/bmc"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetBMCClientFromBMCName(ctx context.Context, c client.Client, bmcName string, insecure bool) (bmc.BMC, error) {
	bmc := &metalv1alpha1.BMC{}
	if err := c.Get(ctx, client.ObjectKey{Name: bmcName}, bmc); err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("BMC %q not found", bmcName)
		}
		return nil, err
	}
	return GetBMCClientFromBMC(ctx, c, bmc, insecure)
}

func GetBMCClientFromBMC(ctx context.Context, c client.Client, bmcObj *metalv1alpha1.BMC, insecure bool) (bmc.BMC, error) {
	endpoint := &metalv1alpha1.Endpoint{}
	if err := c.Get(ctx, client.ObjectKey{Name: bmcObj.Spec.EndpointRef.Name}, endpoint); err != nil {
		return nil, fmt.Errorf("failed to get Endpoints for BMC: %w", err)
	}

	bmcSecret := &metalv1alpha1.BMCSecret{}
	if err := c.Get(ctx, client.ObjectKey{Name: bmcObj.Spec.BMCSecretRef.Name}, bmcSecret); err != nil {
		return nil, fmt.Errorf("failed to get BMC secret: %w", err)
	}

	protocol := "https"
	if insecure {
		protocol = "http"
	}

	var bmcClient bmc.BMC
	switch bmcObj.Spec.Protocol.Name {
	case ProtocolRedfish:
		bmcAddress := fmt.Sprintf("%s://%s:%d", protocol, endpoint.Spec.IP, bmcObj.Spec.Protocol.Port)
		username, password, err := GetBMCCredentialsFromSecret(bmcSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to get credentials from BMC secret: %w", err)
		}
		bmcClient, err = bmc.NewRedfishBMCClient(ctx, bmcAddress, username, password, true)
		if err != nil {
			return nil, fmt.Errorf("failed to create Redfish client: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported BMC protocol %s", bmcObj.Spec.Protocol.Name)
	}
	return bmcClient, nil
}

func GetBMCCredentialsFromSecret(secret *metalv1alpha1.BMCSecret) (string, string, error) {
	username, ok := secret.Data["username"]
	if !ok {
		return "", "", fmt.Errorf("no username found in the BMC secret")
	}
	password, ok := secret.Data["password"]
	if !ok {
		return "", "", fmt.Errorf("no password found in the BMC secret")
	}
	return base64.StdEncoding.EncodeToString(username), base64.StdEncoding.EncodeToString(password), nil
}

func GetServerNameFromBMCandIndex(index int, bmc *metalv1alpha1.BMC) string {
	return fmt.Sprintf("compute-%d-%s", index, bmc.Name)
}

func GetBMCNameFromEndpoint(endpoint *metalv1alpha1.Endpoint) string {
	return fmt.Sprintf("bmc-%s", endpoint.Name)
}

func GetBMCSecretNameFromEndpoint(endpoint *metalv1alpha1.Endpoint) string {
	return fmt.Sprintf("bmc-%s", endpoint.Name)
}

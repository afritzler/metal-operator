# BIOSSettingsSet

BIOSSettingsSet represents a set of BIOSSettings to perform operations for all selected physical servers through labels. It updates the BIOS settings on all selected physical servers' BIOS through BIOSSettings.

## Key Points

- BIOSSettingsSet uses a label selector to select the Servers to create BIOSSettings for.
- BIOSSettingsSet creates BIOSSettings for each server which matches the label.
    - Only one BIOSSettings can be active per Server at a time.
- BIOSSettingsSet monitors changes to Server resources and creates/deletes BIOSSettings.

## Workflow

1. BIOSSettingsSet filters Servers matching the provided label.
2. BIOSSettingsSet creates a BIOSSettings CRD for each Server selected.
3. BIOSSettingsSet monitors the created BIOSSettings and updates the status.
4. BIOSSettingsSet creates or deletes BIOSSettings based on the changes to Server CRDs.

## Example

```yaml
apiVersion: metal.ironcore.dev/v1alpha1
kind: BIOSSettingsSet
metadata:
  name: biossettingsset-sample
spec:
  biosVersionTemplate:
    version: "U59 v2.34 (10/04/2024)"
    serverMaintenancePolicy: OwnerApproval
    settings:
      foo: bar
  serverSelector:
    matchLabels:
      manufacturer: "dell"
```

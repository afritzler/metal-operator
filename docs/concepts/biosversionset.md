# BIOSVersionSet

BIOSVersionSet represents a set of BIOSVersion to perform operations for all selected physical servers through labels. It updates the BIOS version on all selected physical servers' BIOS through BIOSVersion.

## Key Points

- BIOSVersionSet uses a label selector to select the Servers to create BIOSVersion for.
- BIOSVersionSet creates BIOSVersion for each server which matches the label.
    - Only one BIOSVersion can be active per Server at a time.
- BIOSVersionSet monitors changes to Server resources and creates/deletes BIOSVersion.

## Workflow

1. BIOSVersionSet filters Servers matching the provided label.
2. BIOSVersionSet creates a BIOSVersion CRD for each Server selected.
3. BIOSVersionSet monitors the created BIOSVersion and updates the status.
4. BIOSVersionSet creates or deletes BIOSVersion based on the changes to Server CRDs.

## Example

```yaml
apiVersion: metal.ironcore.dev/v1alpha1
kind: BIOSVersionSet
metadata:
  name: biosversionset-sample
spec:
  biosVersionTemplate:
    version: "U59 v2.34 (10/04/2024)"
    image:
      URI: "https://foo-2.34_10_04_2024.signed.flash"
      transferProtocol: "HTTPS"
    updatePolicy: Normal
    serverMaintenancePolicy: OwnerApproval
  serverSelector:
    matchLabels:
      manufacturer: "dell"
```

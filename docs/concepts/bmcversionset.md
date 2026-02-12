# BMCVersionSet

BMCVersionSet represents a set of BMCVersion to perform operations for all selected physical BMCs through labels. It updates the BMC version on all selected physical servers' BMC through BMCVersion.

## Key Points

- BMCVersionSet uses a label selector to select the BMCs to create BMCVersion for.
- BMCVersionSet creates BMCVersion for each BMC which matches the label.
    - Only one BMCVersion can be active per BMC at a time.
- BMCVersionSet monitors changes to BMC resources and creates/deletes BMCVersion.

## Workflow

1. BMCVersionSet filters BMCs matching the provided label.
2. BMCVersionSet creates a BMCVersion CRD for each BMC selected.
3. BMCVersionSet monitors the created BMCVersion and updates the status.
4. BMCVersionSet creates or deletes BMCVersion based on the changes to BMC CRDs.

## Example

```yaml
apiVersion: metal.ironcore.dev/v1alpha1
kind: BMCVersionSet
metadata:
  name: bmcversionset-sample
spec:
  bmcVersionTemplate:
    version: "U59 v2.34 (10/04/2024)"
    image:
      URI: "https://foo-2.34_10_04_2024.signed.flash"
      transferProtocol: "HTTPS"
    updatePolicy: Normal
    serverMaintenancePolicy: OwnerApproval
  bmcSelector:
    matchLabels:
      manufacturer: "dell"
```

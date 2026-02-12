# BMCVersion

BMCVersion represents a BMC version upgrade operation for a physical server's manager. It updates the BMC version on a physical server's BMC.

## Key Points

- BMCVersion maps a BMC version required for a given server's BMC.
    - BMCVersion spec contains the required details to upgrade the BMC to the required version.
- Only one BMCVersion can be active per BMC at a time.
- BMCVersion starts the version upgrade of the BMC using the Redfish `SimpleUpgrade` API.
- BMCVersion handles reboots of the BMC.
- BMCVersion requests Maintenance if `serverMaintenancePolicy` is set to OwnerApproval.
- Once BMCVersion moves to Failed state, it stays in this state unless manually moved out.

## Workflow

1. A separate operator (e.g., BMCVersionSet) or user creates a BMCVersion resource referencing a specific BMC.
2. The provided settings are checked against the current BMC version.
3. If the version is the same as on the server's BMC, the state is moved to Completed.
4. If OwnerApproval `serverMaintenancePolicy` type is requested and ServerMaintenance is not provided already, it requests one per Server managed by the BMC and waits for all Servers to enter Maintenance state.
5. BMCVersion issues the BMC upgrade using the Redfish `SimpleUpgrade` API and monitors the upgrade task created by the API.
6. BMCVersion moves to Failed state:
    - If `SimpleUpgrade` is issued but unable to get the task to monitor the progress of the BMC upgrade.
    - If the upgrade task created by `SimpleUpgrade` fails and does not reach completed state.
    - If the BMC version requested is lower than the current BMC version.
7. BMCVersion moves to reboot the BMC once the upgrade task has been completed.
8. BMCVersion verifies the BMC version post reboot, removes the ServerMaintenance resource if created by itself, and transitions to Completed state.
9. Any further update to the BMCVersion spec will restart the process.

## Example

```yaml
apiVersion: metal.ironcore.dev/v1alpha1
kind: BMCVersion
metadata:
  name: biosversion-sample
spec:
  version: 2.10.3
  image:
    URI: "http://foo.com/dell-idrac-bmc-2.10.3.bin"
    transferProtocol: "HTTP"
    imageSecretRef:
      name: sample-secret
  updatePolicy: Force
  bmcRef:
    name: BMC-sample
  serverMaintenancePolicy: Enforced
```

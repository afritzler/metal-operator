# BIOSVersion

BIOSVersion represents a BIOS version upgrade operation for a physical server (compute system). It updates the BIOS version on a physical server's BIOS.

## Key Points

- BIOSVersion maps a BIOS version required for a given server's BIOS.
    - BIOSVersion spec contains the required details to upgrade the BIOS to the required version.
- Only one BIOSVersion can be active per Server at a time.
- BIOSVersion starts the version upgrade of the BIOS using the Redfish `SimpleUpgrade` API.
- BIOSVersion handles reboots of the server using a ServerMaintenance resource.
- Once BIOSVersion moves to Failed state, it stays in this state unless manually moved out.

## Workflow

1. A separate operator (e.g., BIOSVersionSet) or user creates a BIOSVersion resource referencing a
   specific Server.
2. The provided BIOS version is checked against the current BIOS version.
3. If the version is the same as on the server's BIOS, the state is moved to Completed.
4. If ServerMaintenance is not provided already, it requests one and waits for the Server to enter Maintenance state.
    - The `policy` used by ServerMaintenance is provided through the `serverMaintenancePolicy` spec field in BIOSVersion.
5. BIOSVersion issues the BIOS upgrade using the Redfish `SimpleUpgrade` API and monitors the upgrade task created by the API.
6. BIOSVersion moves to Failed state:
    - If `SimpleUpgrade` is issued but unable to get the task to monitor the progress of the BIOS upgrade.
    - If the upgrade task created by `SimpleUpgrade` fails and does not reach completed state.
    - If the BIOS version requested is lower than the current BIOS version.
7. BIOSVersion moves to reboot the server once the upgrade task has been completed.
8. BIOSVersion verifies the BIOS version post reboot, removes the ServerMaintenance resource if created by itself, and transitions to Completed state.
9. Any further update to the BIOSVersion spec will restart the process.

## Example

```yaml
apiVersion: metal.ironcore.dev/v1alpha1
kind: BIOSVersion
metadata:
  name: biosversion-sample
spec:
  version: "U59 v2.34 (10/04/2024)"
  image:
    URI: "https://foo-2.34_10_04_2024.signed.flash"
    transferProtocol: "HTTPS"
  serverRef:
    name: endpoint-sample-system-0
  serverMaintenancePolicy: OwnerApproval
```

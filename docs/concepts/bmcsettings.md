# BMCSettings

BMCSettings represents a BMC setting update operation for a physical server's BMC (compute system). It updates the BMC settings on a physical server's BMC.

## Key Points

- BMCSettings maps a BMC version and settings as a map for a given server.
- Only one BMCSettings can be active per BMC at a time.
- BMCSettings related changes are applied once the BMC version matches the physical server's BMC version.
- BMCSettings handles reboots of the BMC.
- BMCSettings requests Maintenance; `serverMaintenancePolicy` determines the maintenance type.
- Once BMCSettings moves to Failed state, it stays in this state unless manually moved out.

## Workflow

1. A separate operator (e.g., BMCSettingsSet) or user creates a BMCSettings resource referencing a specific BMC.
2. The provided settings are checked against the current BMC settings.
3. If settings are the same as on the server, the state is moved to Applied (even if the version does not match).
4. If the settings need an update, BMCSettings checks the version of the BMC. If the required version does not match, it waits for the BMC version to reach the spec version.
5. If ServerMaintenance is not provided already, it requests one per Server managed by the BMC and waits for all Servers to enter Maintenance state.
6. The setting update process is started and the physical server's BMC is rebooted if required.
7. BMCSettings verifies the settings have been applied and transitions the state to Applied. It removes all ServerMaintenance resources if created by itself.
8. Any further update to the BMCSettings spec will restart the process.
9. If BMCSettings fails to apply the BMC settings, it moves to Failed state until manually moved out.

## Example

```yaml
apiVersion: metal.ironcore.dev/v1alpha1
kind: BMCSettings
metadata:
  name: bmcsettings-sample
spec:
  bmcRef:
    name: sample-bmc
  version: 2.10.3
  settings:
    otherSettings: "123"
    someother: Disabled
  serverMaintenancePolicy: OwnerApproval
```

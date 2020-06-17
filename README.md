# PowerTrust
A Go service which exposes an API to sign PowerShell scripts.

## Usage

### Service Daemon

```
powertrust service
```
Launches a service that listens on 7974 for PowerShell script uploads. Requires a signing certificate be loaded into `Cert:\LocalMachine\My`. 

### CLI Upload

```
powertrust sign "http://localhost:7974" "C:\Users\me\myscript.ps1"
```
Will upload, sign, download, and then promptly delete the script from the remote server in one command.

# CSI Certification Suite
## [Background](https://docs.google.com/document/d/1XzPogq3TFUUhWGNvW33UNJM0CeKo51EKp-WhY4D9gOA)
## Roadmap
- (P0) Basic sanity check BeforeEach test - make sure needed components are available
- (P0) As a storage vendor I want to run a certification test against my CSI driver that validates it against the CSI SPEC.
- (P0/P1) As a storage vendor & Kube admin I want to run a certification tests against my CSI driver that validates it against a CO
- (P2) As a kube admin I want to validate that the CSI driver I’m using is the same one that was tested (ie. hash, version etc..)

## (P0) Basic CSI API validation Suite
- The [kubernetes-csi/csi-test](https://github.com/kubernetes-csi/csi-test) test suite does the basic CSI API validation
  - This is a sanity check that simply checks if the CSI driver conforms to the [CSI Spec](https://github.com/container-storage-interface/spec) 
  - You can refer to this [spreadsheet](https://docs.google.com/spreadsheets/d/1cyGLU_zEyq-i6D5FJpDu-jM2oTynPupbO1KrGCrrDVw/edit?usp=sharing) to see all the test cases that is covered by the sanity test

#### Running the CSI API validation on a Driver
- Clone the csi-test [repo](https://github.com/kubernetes-csi/csi-test)
- Build the csi-sanity tool: `cd go/src/github.com/kubernetes-csi/csi-test/cmd/csi-sanity/` and run `make all`
- Run a CSI Driver (In this example the [ebs csi driver](https://github.com/bertinatto/ebs-csi) is used) 
  - Launch an AWS EC2 instance, connect to it, and clone the driver repository
  - `cd go/src/github.com/bertinatto/ebs-csi/`
  - Install Dependencies: `dep ensure`
  - Build: `make ebs-csi-driver`
  - Run: `bin/ebs-csi-driver -endpoint tcp://127.0.0.1:10000 -logtostderr -v 5`
- Run the sanity test on the ebs driver
  - `cd go/src/github.com/kubernetes-csi/csi-test/cmd/csi-sanity/`
  - `./csi-sanity -csi.endpoint 127.0.0.1:10000`
- The results of the test run will be printed

## (P1) As a storage vendor & Kube admin I want to run a certification tests against my CSI driver that validates it against a CO 
#### Gap Analysis (What needs to be added on top of the sanity tests for P1)
- Add API Validation tests for Topology and Quota
- Requires Functional tests (Ensure that they actually work in kubernetes) for the following:
  - Provision
  - Delete
  - Attach
  - Detach
  - File write / read validation on CSI PV
  - Block volume read/write validation on CSI PV
  - Resize
  - Quota
  - Topology
  - Create Snapshot
  - Delete Snapshot 
  - Test against non dynamically provisioned volume 
### [Design & Current Implementation Status](https://github.com/mathu97/CSI-Certification-Suite/blob/master/Design.md)

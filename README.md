# CSI Certification Suite
## [Background](https://docs.google.com/document/d/1XzPogq3TFUUhWGNvW33UNJM0CeKo51EKp-WhY4D9gOA)
## Basic CSI API validation Suite
- The [kubernetes-csi/csi-test](https://github.com/kubernetes-csi/csi-test) repository is the Basic CSI API validation Suite
  - This is a sanity check that simply checks if the CSI driver conforms to the [CSI Spec](https://github.com/container-storage-interface/spec) 
  - You can refer to this [spreadsheet](https://docs.google.com/spreadsheets/d/1cyGLU_zEyq-i6D5FJpDu-jM2oTynPupbO1KrGCrrDVw/edit?usp=sharing) of all the test cases that the sanity test covers.

### Running the CSI API validation Suite
- Clone the csi-test [repo](https://github.com/kubernetes-csi/csi-test)
- Build the csi-sanity tool: `cd cmd/csi-sanity/` and run `make all`
- Run a CSI Driver (In this example we use the [ebs csi driver](https://github.com/bertinatto/ebs-csi)) 
  - Build: `make ebs-csi-driver`
  - Run: `bin/ebs-csi-driver -endpoint tcp://127.0.0.1:10000 -logtostderr -v 5`
- Run sanity test on the ebs driver
  - `cd csi-test/cmd/csi-sanity/`
  - `./csi-sanity -csi.endpoint 127.0.0.1:10000`
- You will a list of the failed test cases, as well as the number of test cases that passed or were skipped

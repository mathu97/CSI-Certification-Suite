# Design and Current Status
## Summary 
	A Certification Framework that shows if any given CSI plugin is kubernetes certified or not.
## Motivation 
A CSI plugin compliant with the CSI spec should "just work" with kubernetes. Currently, csi-sanity exists to help test compliance with the spec. But e2e testing of plugins is needed as well to provide plugin authors and users stronger guarantees that their plugin "just works." In kubernetes upstream, e2e tests are being run for select CSI and in-tree plugins. There should be an easy way for other CSI plugins to run the same tests for their plugin and to write their own tests, and for users of these plugins to tell that they've passed the tests.	
### Goals 
Provide a script that will certify a CSI plugin by testing it’s conformance to the Spec (csi-test suite) as well as it’s functionality in Kubernetes (e2e-tests)
### Non-Goals 
## Proposal 
### User Stories 
Story 1  
  As a storage vendor I want to run a certification test against my CSI plugin that validates it against the CSI SPEC.  
Story 2  
  As a storage vendor & Kube admin I want to run certification tests against my CSI plugin that validates it against a CO.  
Story 3  
As a kube admin I want to validate that the CSI plugin I’m using is the same one that was tested (ie. hash, version etc..) 

### Implementation Details 
### Risks and Mitigations 
## Graduation Criteria 
## Implementation History 

Issues around making kubernetes/kubernetes e2e storage tests pluggable:  
- https://github.com/kubernetes/kubernetes/issues/69819 
  - Discussion for CSI E2E tests 
- https://github.com/kubernetes/kubernetes/issues/72288 
- https://github.com/kubernetes/kubernetes/issues/71237 
- https://github.com/kubernetes/kubernetes/issues/72242 
- https://github.com/kubernetes/kubernetes/issues/70258 

Important PRs:  
- https://github.com/kubernetes/kubernetes/pull/68025 (merged) 
  - Enables in-tree drivers and CSI drivers to share the same tests 
  - This takes care a lot of the tests-case work for CSI-Certification, as now CSI Drivers will be tested against the same test cases as In-Tree Drivers 
- https://github.com/kubernetes/kubernetes/pull/68483 (merged) 
- https://github.com/kubernetes/kubernetes/pull/70992
  - Clean up work for E2E Storage Tests
- https://github.com/kubernetes/kubernetes/pull/72434
  - More refactoring/clean up
- https://github.com/kubernetes/kubernetes/pull/72002 
  - Adding more types of tests that would be useful for CSI Drivers

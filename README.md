# CR Manager

This code uses k8s client-go to create a custom resource of a particular kind (Keevakind), populating each of the fields with sample (random) values.

## Prerequisites

- A cluster (e.g. minikube) is up and running and available (e.g. current cluster in  .kube config)

## Instructions
To run it, just type:`go build`
Then run the execuatable `cr-manager`.

If the CRD is not already on the cluster, apply as follows:

`kubectl apply -f ./keevakind-crd.yaml`

Currently it only creates but I will be adding other functions soon.

The Keevakind CRD is also included in this project as an example - it's source of truth is in:
 
   `github.com/speak2jc/k-op `
This project is imported here in order to use the types.

## Using the pattern 
The assumption in a real working environment is that the cluster bootstrap would create the CRD (using k-op project YAML files)

If changes were to occur in the k-op project which imposed breaking changes on the structure of the CRD, this would be represented by a semver version update, allowing the refernece in the go.mod of this project to be updated in a branch and merged back when tested.
Both the old and new versions of the CRD would need to co-exist on the cluster (if possible). 
If not feasible, the change would need to be otherwise managed (e.g. try to avoid breaking changes and coordinate carefully if unavoidable).
 

## Further work
- Add functionality to update, delete, list and get Keevakind objects
- Find out how to make fields mandatory in CRD spec - currently everything is optional
- This project also contains code to generate the CRD on the cluster - this is incomplete but lo-prio as the CRD would not typically be created from this project.
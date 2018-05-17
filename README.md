# cloud-launcher
Utility for launching cloud instances for using as OpenShift nodes

## Quick Start
```
$ cat ~/.aws/credentials 
[default]
aws_access_key_id = mykeyid
aws_secret_access_key = mysecretaccesskey

make
./cloud-launcher start --cluster-name=mycluster --token=mytoken
cd ../../openshift/openshift-ansible/
ansible-playbook -i ~/mycluster.inventory playbooks/prerequisites.yml
ansible-playbook -i ~/mycluster.inventory vi playbooks/deploy_cluster.yml

```

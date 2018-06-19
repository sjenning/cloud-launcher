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
git fetch
git checkout origin/master
ansible-playbook -i ~/mycluster.inventory playbooks/prerequisites.yml
ansible-playbook -i ~/mycluster.inventory playbooks/deploy_cluster.yml

```

## Using Your Cluster
To use the cluster, you will need to ssh to the master node and run commands from there, for example `oc login -u system:admin`.

## Creating Your Credentials
If you're familiar with the process of logging in to AWS, you can skip this part.  If you're new, or need a refresher, please read this.

The process of getting an aws account may be found here: https://mojo.redhat.com/docs/DOC-1081313#jive_content_id_Amazon_AWS.  This will create ~/.aws/credentials as above.

Get your token by logging in to https://console.reg-aws.openshift.com/console/ using Google login.  After you log in, click your name at the top right of the page, and use `Copy Login Command` to copy the login command to the clipboard.  The token may be found at the end of the login command thus provided.  Your token must be renewed every 30 days.  If your token has expired, you will succeed in launching the cluster, but when you try to run the prerequisite playbook, it will fail with the following error:

```
FAILED - RETRYING: Create credentials for docker cli registry auth (1 retries left).
```

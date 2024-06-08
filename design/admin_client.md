# Admin Client CLI

## Purpose

This client is responsible for managing the RAFT cluster itself.
It's mainly responsible for:

- creating a hole new cluster with 5 nodes (standard from protocol)

- changing the configuration of the cluster 

- retrieving the current configuration of the cluster

- listing the nodes in the cluster with relative ips (private and public)

- delete the entire cluster with all the vms

After every operation the client will receive a return value of this form 
if it's not specified a specific return value type

    {

        STATUS: [SUCCESS, FAILURE]

        description: string

    }

## Useful definition 

- **cluster configuration** : an ADT with the following structure:

    {

         nodes : []string

         leaderIp : string //optional

    }
    > Be sure that your list of nodes is of the form ["ip1","ip2",...]

- **change cluster configuration request** : an ADT with the following structure:

    {

         OP_TYPE : [CHANGE, NEW]

         cluster configuration

    }

    The configuration must be absolute not relative

- **Establish connection with the leader**: the cluster must already exist

    1. connect to a random node in the cluster (which one is irrelevant)

    2. wait response of node of type {leaderIP: string}

    3. check if the ip in the answer is the same of the ip of the node you are talking to

        - if it's different:

            31. close the connection with this node 

            32. connect to the other 

                323. if the answer is blank, wait a bit and restart from point 1

            33. goto point 2

        - if it's the same:

            31. connection with leader established
         
- **info request** : an ADT with the following structure:

    {

        REQ_TYPE: [config]

        timestamp: string

    }

- **info response** : an ADT with the following structure:
The timestamp is of the corresponding info request

    {

        REQ_TYPE: [config]

        original_timestamp: string

        payload: []byte

    }

## How it operates

In this section is described how each functionality has to be implemented:

### Creating a hole new cluster with 5 nodes (standard from protocol)

This functionality is usable only if it does ***not already exist*** a cluster in the network.
The procedure consists of 5 points:

1. Create 5 vm 
2. Save their ips
3. Create a new change cluster configuration, of type NEW, with the 5 ips you saved
4. Send the created change cluster configuration to one of the 5 nodes (which one is irrelevant)

### Changing the configuration of the cluster 

This functionality is usable only if it does ***already exist*** a cluster in the network.
There two possible changes you can apply to the cluster, each of them has its own procedure:

1. ADD a set of existing nodes to the cluster:

    11. Create a change cluster configuration request with the ips of the node you want to add

    12. Establish connection with the leader 

    13. Send to the leader the change cluster configuration request

2. ADD a number (i) of non-existing nodes to the cluster:

    21. Create i new nodes in the network

    22. create a set S with the ips of the new nodes

    23. goto point 1 with the set S
    
3. REM a set of existing nodes from the cluster:

    31. Create a change cluster configuration request with the ips of the node you want to remove

    32. Establish connection with the leader 

    33. Send to the leader the change cluster configuration request

    34. When the node to remove are shut off delete them from the network

### Retrieving the current configuration of the cluster

This functionality is usable only if it does ***already exist*** a cluster in the network.

1. Establish connection with the leader 

2. Create info request of type config

3. Send info request to the leader

4. wait response

### listing the nodes in the cluster with relative ips (private and public)

Use libvirt utility to obtain info about the nodes: 
With this command you obtain the list of the nodes
virsh --connect=qemu:///system list | grep VM_RAFT | cut -d ' ' -f 6
With this command you obtain the info about a node ips included
virsh --connect=qemu:///system domifaddr --source arp --domain <NODE_NUM>

### delete the entire cluster with all the vms

Use libvirt utility to obtain info about the nodes: 
With this command you obtain the list of the nodes
virsh --connect=qemu:///system list | grep VM_RAFT | cut -d ' ' -f 6
With this command you obtain the list of the volumes  
volumes=$(virsh --connect=qemu:///system vol-list --pool default| grep disk_raft_node)


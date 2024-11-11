# Mutual Exclusion

## System requirements

This system uses the Ricart & Agrawala algorithm for mutual exclusion.

When the nodes start, they start off by sending a request to all other nodes that tells that the node want to access the critical section (R1). The nodes will all reply wether or not the server can access the critical section at this time. Only one node is granted access to the critical section at any time (R2: Safety). If a node tries to access to critical section while another nodes is accessing it, the node will be added to a queue and told that access is not granted for now. When the node is done using the critical section, the next node in the queue is notified ensuring that all nodes that want to access the critical section, will eventually get access (R3: Liveliness).

## Discussion



## GitHub repository

You can find the repository on my GitHub: <https://github.com/Kanerix/mutual-exclusion>

## System logs

System logs can also be found in the repo: <https://github.com/Kanerix/mutual-exclusion>

# Mutual Exclusion

## System requirements

This system uses the Ricart & Agrawala algorithm for mutual exclusion.

When the nodes start, they start off by sending a request to all other nodes that tells that the node want to access the critical section (R1). The nodes will all reply wether or not the server can access the critical section at this time. Only one node is granted access to the critical section at any time (R2: Safety). If a node tries to access to critical section while another nodes is accessing it, the node will be added to a queue and told that access is not granted for now. When the node is done using the critical section, the next node in the queue is notified ensuring that all nodes that want to access the critical section, will eventually get access (R3: Liveliness).

## Discussion

Node 3001's request is granted by Node 3002: Node 3001 request granted by localhost:3002.
Node 3001 receives and grants an access request from Node 3002: [3001] Access Requested by 3002 followed by [3001] Access Granted to 3002.
Node 3001 processes Node 3003's request: [3001] Access Requested by 3003.
Node 3001 finishes its work and releases access: Node 3001 released access.
Node 3002:
Node 3002 handles access requests in this sequence:
Grants access to Node 3001: [3002] Access Granted to 3001.
Requests and gets granted access by Node 3001: [3002] Access Requested by 3001.
Node 3002 requests access to Node 3003: [3002] Access Requested by 3003.
Enters the critical section after Node 3001 releases: Node 3002 released access.
Node 3003:
Initial requests from Node 3003 are denied: Node 3003 request denied by localhost:3001 and Node 3003 request denied by localhost:3002. This denial indicates that both nodes 3001 and 3002 were busy or had pending higher-priority requests (from a Lamport clock perspective).
When nodes 3001 and 3002 release access, Node 3003 can then proceed: [3003] Access Released by 3001, [3003] Access Released by 3002.
Finally, Node 3003 finishes its work in the critical section: Node 3003 released access.

## GitHub repository

You can find the repository on my GitHub: <https://github.com/Kanerix/mutual-exclusion>

## System logs

System logs can also be found in the repo: <https://github.com/Kanerix/mutual-exclusion>

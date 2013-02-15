#glowing-nemesis

##Introduction

The objective of this lab part is to implement a Failure Detection and a Leader Election Module
that will later be used in your Paxos system. We assume only a partially synchronous system,
with a crash-stop process abstraction.

Failure Detector: A Failure Detector can be implemented either using ping messages or
using active heartbeat messages.
A Failure Detector, pinging the other replicas is shown as Algorithm 2.7 in the book. It
sends a HeartbeatRequest to all other nodes, and if the request is not answered with a
HeartbeatReply within a certain time, it suspects the silent process.
A Failure Detector using active heartbeats is divided into a sending and a receiving process.
The receiving process is the same as the Failure Detector described in Algorithm 2.7, with
the one crucial diﬀerence, that it does not send HeartbeatRequest messages. Instead, the
sending process is a simple loop, that, upon timeout, sends a HeartbeatReply to all other
processes. Thus this Failure Detector uses two timeouts, one for receiving and one for sending
heartbeat messages.

Leader Election and Failure Detection Module: In brief, in case of leader failure, a new
leader should take over the leadership and all non-faulty nodes should accept him as the leader.
You are free to design your own way of leader election, as long as it satisﬁes the assumptions
and requirements. We recommend that you draw inspiration from the pseudocode found in
Sections 2.6.4 and 2.6.5 in the textbook.

These are the requirements expected from your implementation (max 90% score):

• At start-up, the module shall set up connections between several nodes (e.g. 5 nodes),
and there should be a pre-selected leader.

• You can make the module start-up statically. That means that you can choose a few
nodes in the lab, and use their DNS name or IP address and port number statically to set
up your group of machines. These can either be hardcoded in your source code, or read
from a conﬁguration ﬁle.

• Failure Detection should be implemented using one of the two possibilities above. Each
node should monitor the other nodes and forward detected failures and recoveries to the
Leader Election module.

• Leader Election should be implemented according to Algorithm 2.8, or an equivalent one,
reacting to input from the Failure Detector and creating notiﬁcations on leader change.

• The two implementations should be modular, and composable so that they can easily be
reused.

• You should test your code using several of the machines in the lab.

Note that even though you have implemented all of the above, we may decide not to reward a
100% score, e.g. if your code is not as modular as expected, not documented well enough or has
other major deﬁciencies. It is expected that your code be as simple as possible, but
no simpler.

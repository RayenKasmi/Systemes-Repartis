# Some useful info:

## Exchanges and Exchange Types
Exchanges are AMQP 0-9-1 entities where messages are sent to. Exchanges take a message and route it into zero or more queues. The routing algorithm used depends on the exchange type and rules called bindings. AMQP 0-9-1 brokers provide four exchange types:

Exchange type	Default pre-declared names
Direct exchange	(Empty string) and amq.direct
Fanout exchange	amq.fanout
Topic exchange	amq.topic
Headers exchange	amq.match (and amq.headers in RabbitMQ)
Besides the exchange type, exchanges are declared with a number of attributes, the most important of which are:

Name
Durability (exchanges survive broker restart)
Auto-delete (exchange is deleted when last queue is unbound from it)
Arguments (optional, used by plugins and broker-specific features)

### Default Exchange
The default exchange is a direct exchange with no name (empty string) pre-declared by the broker. It has one special property that makes it very useful for simple applications: every queue that is created is automatically bound to it with a routing key which is the same as the queue name.

### Direct Exchange
A direct exchange delivers messages to queues based on the message routing key. A direct exchange is ideal for the unicast routing of messages. They can be used for multicast routing as well.

Here is how it works:

- A queue binds to the exchange with a routing key K
- When a new message with routing key R arrives at the direct exchange, the exchange routes it to the queue if K = R
- If multiple queues are bound to a direct exchange with the same routing key K, the exchange will route the message to all queues for which K = R

[see the full documention here](https://www.rabbitmq.com/tutorials/amqp-concepts)

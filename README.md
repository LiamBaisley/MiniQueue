# MiniQ

Just trying to write an http message queue in less than 1000 lines. 

Inspired by minikeyvalue by @geohot

# Features
### Security
MiniQ takes a secret string as an argument on start up. This secret should be provided in Authorization request headers when interacting with the queue to ensure access.

### Read guarantee
MiniQ requires users to hit the /confirm endpoint once a message has been consumed so that is can be removed from the Queue. This is to ensure all messages are successfully consumed.

### Fast and lightweight.
Because MiniQ is less than 1000 lines and written in Go, it is fast ann incredibly lightweight. We take advantage of the features of LevelDB to support this. 

### Persistent
MiniQ uses LevelDB as its datastore. This means that even if the queue crashes or the host machine unexpectedly shuts down, all messages in the queue are preserved and can be picked up once the error has been resolved. 

### Minimal complexity
This project was written to be easy to use and understand. Removing a lot of the complexity and unnecessary features of larger message queues and giving users only what they need for asynchronous messaging.

## Todo
  - Create topics? Or some kind of queue separation

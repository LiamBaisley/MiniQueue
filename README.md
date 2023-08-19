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

# Usage
## Docker
### Using docker commands
**First, open the dockerfile:** <br />
![image](https://github.com/LiamBaisley/MiniQueue/assets/50359625/76819004-1c1c-4c12-9763-d8d13def784e)
<br />
**Replace the secret highlighted in the screenshot above with a strong, secure string of your choice.**
**If you would like to use the queue without consumption guarantees (not recommended) then add -c to the CMD line in the dockerfile like this:**
![image](https://github.com/LiamBaisley/MiniQueue/assets/50359625/e83f20b9-a21b-4e35-be74-9da98f6db5b0)


In a shell of your preference: 
* Run: Docker build -t miniq .
* Run: Docker run -p <port of your choice>:8080 -d miniq
* **Example**
  * ![image](https://user-images.githubusercontent.com/50359625/210170503-f239c90b-ee0b-429f-bcef-81e9547f8d44.png)
  * ![image](https://user-images.githubusercontent.com/50359625/210170528-81c78e29-8cba-4fa5-b44f-d82341d55403.png)
 
## Docker Compose
**First, open the dockerfile:** <br />
![image](https://github.com/LiamBaisley/MiniQueue/assets/50359625/76819004-1c1c-4c12-9763-d8d13def784e)
<br />
**Replace the secret highlighted in the screenshot above with a strong, secure string of your choice.**
**If you would like to use the queue without consumption guarantees (not recommended) then add -c to the CMD line in the dockerfile like this:**
![image](https://github.com/LiamBaisley/MiniQueue/assets/50359625/e83f20b9-a21b-4e35-be74-9da98f6db5b0)
**Run ```docker-compose up``` in your terminal from the project root.** 
**If you would like to run it detatched (no active terminal output, best when running unsupervised or on a server) run ```docker-compose up -d```**
**The queue will be available on http://localhost:9000, when running on a server route requests here using something like NGINX or Apache Tomcat** 

### Once MiniQ is running here are the steps for consumption and publishing of messages:
**1: To push a message to the Queue**
  * Make a Post request to "queueurl:port/message" with your message in the request body and the Authorization header set to the secret which you configured.<br />
  ![image](https://user-images.githubusercontent.com/50359625/210170836-159c0221-013f-41b4-8d0a-137cc7080515.png)
  
**2: To read a message**
  * Make a get request to "queueurl:port/message" with the Authorization header set to the secret which you configured.<br />
  ![image](https://user-images.githubusercontent.com/50359625/210170759-8f789b97-d3cd-4b16-a7c4-10b120318ba8.png)
  
**3: After consuming a message**
  * Make a post request to "queueurl:port/confirm" with the Authorization header set to the secret which you configured and send the message key (returned in the request to read the message) to confirm that the message has been consumed and can be removed from the queue. <br />
  ![image](https://user-images.githubusercontent.com/50359625/210170831-a3482f22-8cf9-4bcc-8952-103d778efbb9.png)
  
## Todo
  - Create topics? Or some kind of queue separation
  - ~~Timestamp queue messages~~
  - Distribute and replicate queue for higher fault tolerance
  - add functionality to read entire queue at once. 

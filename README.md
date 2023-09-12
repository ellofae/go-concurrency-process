# go-concurrency-process
____

Implementation of service that manages users and their available privileges. Users are associated with specific privileges for further testing of the system.

**Example of privileges:**

* `READONLY_ACCESS`
* `WRITEONLY_ACCESS`
* `CONFIG_ADMISSION`
* `ADMIN_ACCESS`

**Example of results**:


| User  | Privileges user is attached to |
| ------------- | ------------- |
| 1000  | `READONLY_ACCESS`, `ADMIN_ACCESS`, `CONFIG_ADMISSION`  |
| 2000  | `READONLY_ACCESS`, `ADMIN_ACCESS`, ... |
| 3000  | None  |

## HTTP API Structure
___
Public endpoints:

    GET /priv - get privileges by title
    GET /priv/user - get users' privileges list

    POST /priv - create a privilege
    POST /priv/user/add - add privileges to user
    POST /priv/user/remove - remove privileges of user

    DELETE /priv/{id:[0-9]+} - delete a specific privelege
    DELETE /priv/user/{id:[0-9]+} - delete all user's priveleges

## Requests (Linux)

#### Privilege creation:

`curl -XPOST http://localhost:8000/priv -d '{"privilege_title": "READONLY_ACCESS"}'`

#### Attach privileges to specific user:

`curl -XPOST http://localhost:8000/priv/user/add -d '{"user_id":1, "add_privilege": ["READONLY_ACCESS", "ADMIN_ACCESS", "CONFIG_ADMISSION"]}'`

#### Remove privileges from specific user:

`curl -XPOST http://localhost:8000/priv/user/remove -d '{"user_id":1, "add_privilege": ["READONLY_ACCESS", "ADMIN_ACCESS"]}'`

#### Get users' privileges:

`curl -XGET http://localhost:8000/priv/user`

#### Get existing privileges:

`curl -XGET http://localhost:8000/priv`


## Architecture
____
The project follows clean architecture principles by organizing the code into several independent layers. This approach ensures encapsulation and separation of concerns between components, as well as flexibility and code portability. 

In accordance with the principles of clean architecture, the project facilitates communication between layers using interfaces. This allows for loose coupling, dependency inversion, and modularity.

## Patterns
___
The project utilizes the CQRS (Command Query Responsibility Segregation) and SRP (Single Responsibility Principle) patterns.



## Concurrency part
___
The concurrency part is an implementation of working with channels and goroutines. Specifically, it includes a pipeline implementation to control race conditions during the execution of multiple parallel requests. The aim is to ensure thread safety and prevent data races.

The implementation has been thoroughly tested using ApacheBench to verify its effectiveness.

#### Benchmarking

![main](https://i.imgur.com/vWCVFfo.png)

#### Connection time

![main](https://i.imgur.com/WKFmIYZ.png)
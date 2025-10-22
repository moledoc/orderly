# orderly

`orderly` is project management system/approach alternative to Jira.

## TODOs

- [] initial UI
    - use errr-s
    - when creating an order, delegated task id parent ids should point to the new order we create. Need to change some validations etc and how orders are created.
- [] split Write method to specific methods
- [] add necessary new endpoints
    - [] getting orders a user is accountable for
    - [] get user by email
- [] fix testcases
    - task.accountable went from id to user obj
    - task.state went from str to *str
- [] MAYBE: login
- [] authorization
- [] integrate with postgres
- [] run in docker
- [] documentation:
    - [] README
    - [] swagger
- [] link user objects in orders and add validations
    - by validations I mean like unique emails, etc
- [] MAYBE: proper logging and separate spans from logs 

### Smaller TODOs

* accept correct Content-Type
* MAYBE: TODO: pagination
* implement changelog/diff between obj versions
* MAYBE: soft-delete obj

## Author

Meelis Utt
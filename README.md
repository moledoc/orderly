# orderly

`orderly` is project management system/approach alternative to Jira.

## TODOs

- [] order UI
    - [] check [order.html](./order.html) for first idea draft
    - [] idea to use htmx, so there might be some UI related endpoints
        - POST/GET/PATCH/DELETE /ui/<endpoint>
- [] user UI
    - [] idea to use htmx, so there might be some UI related endpoints
        - POST/GET/PATCH/DELETE /ui/<endpoint>
- [] MAYBE: login
- [] authorization
- [x] performance testing
    - performance improvements on local storage: getting all and subordinates/suborders
- [] integrate with postgres
- [] run in docker
- [] documentation:
    - [] README
    - [] swagger
- [] link user objects in orders and add validations
- [] MAYBE: proper logging and separate spans from logs 

### Smaller TODOs

* accept correct Content-Type
* MAYBE: TODO: pagination
* implement changelog/diff between obj versions
* MAYBE: soft-delete obj

## Author

Meelis Utt
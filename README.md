# orderly

`orderly` is project management system/approach alternative to Jira.

## TODOs

- [x] user management
    * REST API: POST, GET (by id and listing), PATCH, DELETE
        * POST /user
        * GET /user/{id}
        * GET /users
        * GET /user/{id}/subordinates
        * PATCH /user
        * DELETE /user/{id}\[?type=hard\]
- [x] order managment
    * REST API: POST, GET (by id and listing), PATCH, DELETE
        * POST /order
        * GET /order/{id}
        * GET /orders
        * GET /order/{id}/suborders
        * GET /order/{id}/delegated_tasks
        * PATCH /order/{id}
        * DELETE /user/{id}\[?type=hard\]
* order UI
    * check [order.html](./order.html) for first idea draft
    * idea to use htmx, so there might be some UI related endpoints
        * POST/GET/PATCH/DELETE /ui/<endpoint>
* user UI
    * idea to use htmx, so there might be some UI related endpoints
        * POST/GET/PATCH/DELETE /ui/<endpoint>
- [x] input validation
* accept correct Content-Type
* login(?)
* testing (functional, performance)?
    * refactor the structure
* deploying (docker)
* documentation (README, swagger)
* MAYBE: TODO: pagination
* logging
* implement changelog
* MAYBE: soft-delete

## Author

Meelis Utt
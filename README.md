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
* order managment
    * REST API: POST, GET (by id and listing), PATCH, DELETE
        * POST /order
        * GET /order/{id}
        * GET /orders
        * GET /order/{id}/suborders
        * PATCH /order/{id}
        * DELETE /order/{id}
* order UI
    * check [order.html](./order.html) for first idea draft
    * idea to use htmx, so there might be some UI related endpoints
        * POST/GET/PATCH/DELETE /ui/<endpoint>
* user UI
    * idea to use htmx, so there might be some UI related endpoints
        * POST/GET/PATCH/DELETE /ui/<endpoint>
* login(?)
* optimize versioning
* testing (functional, performance)?
* deploying (docker)
* documentation (README)

## Author

Meelis Utt
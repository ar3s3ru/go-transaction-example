# go-transaction-example

An example application to show independent `*sql.Tx` propagation over
different decoration layers, for separation of concerns.

- [go-transaction-example](#go-transaction-example)
  - [Context](#context)
  - [Problem](#problem)
  - [More problems](#more-problems)
  - [(One possible) Solution](#one-possible-solution)
    - [`context.Context` to the rescue!](#contextcontext-to-the-rescue)
    - [Transactional decorator](#transactional-decorator)
      - [Pros](#pros)
      - [Cons](#cons)
  - [Project](#project)
  - [License](#license)

## Context

I've come across cases where, for a single Use-case invocation, the Service
had to **update different parts** of the datasource **in a transaction**.

An example might be:
* Use-case: _Create new User_
* Database structure:
  * **Users** table: contains all the Users of the application
  * **Users History** table: append-only table, containing historical events happened for **Users**
* Flow:
  1. _POST /users?name=\<name>&age=\<age>_
  2. Parse input parameters
  3. Call use-case interactor
     * This is the piece of code that contains the use-case business logic
  4. Call the repository
  5. Return the result

## Problem

From the example above, let's talk about point `4.`

Each time some _Users_ get created or updated (or deleted), a new event **has to be
appended** to the _Users History_.

The simplest solution would be to have a Repository interface, like this:
```go
// I'll only consider User creation for simplicity.
type Adder interface {
    Add(context.Context, User) (User, error)
}
```
and a Repository implementation that will update **both tables**:
```go
package postgres

type UsersRepository struct {
    *sql.DB
}

func (ur UsersRepository) Add(ctx context.Context, user User) (User, error) {
    // Start transaction
    tx, err := ur.BeginTx(ctx)

    // Insert into Users table (and get the committed row)
    tx.QueryContext(ctx,
        "INSERT INTO users (...) VALUES (...) RETURNING *;",
        ...
    )

    // Insert into Users History table
    tx.ExecContext(ctx,
        "INSERT INTO users_history (...) VALUES (...);",
        ...
    )
}
```

The Repository implementation above is deliberately simple, but in a real-world
scenario it might probably _be more complicated_ (using a query DSL, sanitizing or
mapping domain values to database types, etc...).

Hence, this strategy, albeit _simpler_, shows the following drawbacks:
1. This Repository implementation is doing _too many things_
2. The Repository code can get long and tedious
3. This implementation is harder to test, since it's touching different parts of the database

## More problems

Let's add another requirement to our service.

We want to expose an endpoint to list all the historical events happened for _Users_,
akin to an **audit log**.

Now, the _Users History_ becomes a **first-class citizen of the service domain**.
As such, we most likely need a _new Repository interface_ to model such an audit log:
```go
type Entry struct {
    // Historical data
}

// It might use some temporal boundaries to get a slice of the all historical log.
type Logger interface {
    Log(ctx context.Context, from, to time.Time) ([]Entry, error)
}
```
and, a new implementation:
```go
type UsersHistoryRepository struct {
    *sql.DB
}

func (r UsersHistoryRepository) Log(ctx context.Context, from, to time.Time) ([]Entry, error) {
    // Implementation here...
}
```

The picture looks like this now:
1. `UsersRepository.Add` inserts both into `users` and `users_history` tables
2. `UsersHistoryRepository.Log` reads from `users_history` table

The _Users History_ is now all over the place, with no component having
the single responsibility of interacting with it, but many different ones instead.

This can be a huge pain to deal with if something has to change with
the history (e.g. table migration, etc.), because it will require _touching all the component of our application_ that are accessing the history.

## (One possible) Solution

Let's turn our faith to the **_Five Only Truths_** (yes, I'm talking about [SOLID]),
and take the [Most Important One].

We could design our Repository implementations like so:
* `UsersRepository`
  * Interacts **only** with `users` table
  * Implements `Adder` interface
* `UsersHistoryRepository`
  * Interacts **only** with `users_history` table
  * Implements `Logger` interface
  * **Decorates** and implements `Adder` interface

The key aspect here is the **decoration** that `UsersHistoryRepository` will do
over the `Adder` interface.

Let's see how the code would look like:
```go
type UsersHistoryRepository struct {
    *sql.DB
    Adder // Embeds the interface, this will be the decorated instance
}

// UsersHistoryRepository.Log as above, won't repeat it here

func (ur UsersRepository) Add(ctx context.Context, user User) (User, error) {
    user, err := ur.Adder.Add(ctx, user)
    if err != nil {
        // Executing decorated instance failed
        return User{}, err
    }

    // Insert into Users History table
    ur.ExecContext(ctx,
        "INSERT INTO users_history (...) VALUES (...);",
        ...
    )

    // ...
}
```
and on our `main.go`/entrypoint/configuration layer, we will create the instances
like so:
```go
usersRepository := UsersRepository{
    DB: db
}

// We will use this Repository interface as Adder in the use-case!
usersHistoryRepository := UsersHistoryRepository{
    DB:    db,
    Adder: usersRepository, // Decorating UsersRepository
}
```

This approach allows us to:
1. Make each Repository implementation **easier**
2. Make maintanability **easier**
3. Make testing **easier**
4. Keep concerns **separated**

However, with this approach we lose the ability of _using transactions_, since
all these modifications happen in a chain of decorators.

But we **need** to use a database transaction to execute side-effects...

<div align="center">
    <img width="240" src="https://media.giphy.com/media/TjS7u7yoMC2KubI5wE/giphy.gif">
</div>

### `context.Context` to the rescue!

From the [Go Blog](https://blog.golang.org/context):

>At Google, we developed a context package that makes it easy to **pass request-scoped values**, cancelation signals, and deadlines **across API boundaries** [..].

So, we can create a `*sql.Tx` to put in the request `context.Context` and
make it available to both Repositories.

However, the `*sql.Tx` lifecycle (`Commit()` and `Rollback()`) must be handled by some component.

### Transactional decorator

We could create an additional decorator that deals with `*sql.Tx` lifecycle
management for a single request, and pass the transaction down the decorators chain
by leveraging `context.Context`.

Later, each component in the chain can access the transaction by using a
_context accessor_ function with a closure, like this:
```go
err = WithTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
    // Use the transaction here
})
```

Taking the previous example, we can design the solution as:
```
use-case
 - calls the repository

 --> transactional-decorator
      - creates a new *sql.Tx
      - puts it in the context
      - calls the decorated interface

    --> users history repository
         - calls the decorated interface

        --> users repository
         - uses WithTransaction to access *sql.Tx
         - inserts into users table

    <-- users history repository
         - uses WithTransaction to access *sql.Tx
         - inserts into users_history table

 <-- transactional-decorator
      - if failed, rollback the transaction
      - if not failed, commit the transaction
      - returns decorated interface result
```

#### Pros

* Each component fulfills a _single responsibility_
* Transaction lifecycle is easier to handle and test
* _Transactional decorator implements all the interfaces that require transactional access_
  * The transactional aspect of a repository access strategy is made more explicit

#### Cons

* Given the current limitations of the Go language, this requires some code duplication on the Transactional decorator
* _Transactional decorator implements all the interfaces that require transactional access_
  * More code to write

## Project

In this repository you can find a project showcasing the approach described above.

To run the application, use:
```
docker-compose up
```

The application exposes two endpoints:
1. `POST /users?name=<name>&age=<age>`
   * `201 Created`: created new User successfully
   * `400 Bad Request`: invalid or missing required query parameters
   * `409 Conflict`: an User with the same name already exists
   * `500 Internal Server Error`: unhandled errors
2.  `GET /users/history?from=<from>&to=<to>`
    * _Note: `from` is required_
    * `200 OK`: lists the
    * `400 Bad Request`: invalid or missing required parameters
    * `500 Internal Server Error`: unhandled errors

The solution described above can be found in `internal/platform/postgres`:
1. `transactional.go` contains the [Transactional decorator](#transactional-decorator)
2. `user_history.go` contains the _Users History_ Repository implementation
3. `user_repository.go` contains the _Users_ Repository implementation

## License

This project is licensed under the [MIT license](LICENSE).

[SOLID]: https://en.wikipedia.org/wiki/SOLID
[Most Important One]: https://en.wikipedia.org/wiki/Single-responsibility_principle

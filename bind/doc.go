// Package bind offers functionality for dependency injection.
//
// Bindings are defined within a context. They follow the same
// semantics as other context APIs. It's possible to override
// bindings inherited form a parent context.
//
// Example
//
//  // DB is a generic database interface
//  type DB interface {
//    Open() (*sql.Conn, err error)
//  }
//
//  // DBImpl implements the database interface
//  type DBImpl struct {
//    Username string `bind:"username"`
//    Password string `bind:"password"`
//  }
//
//  func (db *DBImpl) Open() (*sql.Conn, err error) {
//    cfg := mysql.Config{
//      User:   db.Username,
//      Passwd: db.Password,
//      /* ... */
//    }
//    return sql.Open("mysql", cfg.FormatDSN())
//  }
//
//  // Configure bindings in the context
//  //
//  // We bind a single instance of *DBImpl for DB.
//  // The *DBImpl struct itself will receive the configuration
//  // parameters via the binding context.
//  ctx, _ = bind.Configure(ctx,
//    bind.Instance[string]("admin").For("username"),
//    bind.Instance[string]("s3cr3t").For("password"),
//    bind.Implementation[DB, *DBImpl](),
//    bind.Once[*DBImpl]())
//
//  // Get a reference to the database interface
//  db := bind.Get[DB](ctx)
//  db.Open()
//
//  type UserRepository struct {
//    Database *DB `bind:"-"`
//  }
//
//  // UserRepository was never setup via bind, but we can
//  // still create it and inject all values.
//  repo := bind.New[*UserRepository](ctx)
//  fmt.Println(repo.Database == db) // true
//
// Initialization
//
// For certain types it is important to run additional code after
// injection has happened. For instance, you may want to establish
// the database connection eagerly or do other sanity checks.
//
// If a type implements the bind.Initializer interface the InitAfter
// method will be executed immediately after all values have been
// injected.
//
//  type UserRepository struct {
//    Database *DB `bind:"-"`
//
//    cache map[string]*User
//  }
//
//  func (u *UserRepository) InitAfter() (err error) {
//    // This method is guaranteed to be called after initialization
//    // hence u.Database will not be nil. However, we may want to
//    // initialize other (private) properties.
//
//    cache = make(map[string]*User)
//  }
//
//  repo := bind.New[*UserRepository](ctx)
//  fmt.Println(repo.cache["foo"]) // won't panic
//
package bind

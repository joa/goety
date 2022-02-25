### goety
General purpose Go utilities.

### Package `bind`
Mostly type-safe IoC container.

- `context.Context` based
- `bind.*` methods to setup and resolve bindings
- Thread-safe

#### Example
```go
// first we define some types
// - Database is our database interface
// - sqlDBImpl implements the interfaces
// - dao uses the database
type Database interface { /* ... */ }

type sqlDBImpl struct {
    Username string   `bind:"username"` // get the username
    Password string   `bind:"password"` // get the password
    Host     string   `bind:"host"`     // get the host
    Port     int      `bind:"port"`     // get the port
}

type dao struct {
	DB       Database `bind:"-"`        // get the database 
}

ctx, _ = bind.Configure(ctx,
    bind.String("admin").For("username"),
    bind.String("s3cr3t").For("password"),
    bind.String("localhost").For("host"),
    bind.Port(5432).For("port"),
    bind.ImplementationOnce[Database, *sqlDBImpl](), // bind Database to one instance of *sqlDBImpl
)

// get Database which will be the *sqlDBImpl instance
// this is useful if you'd manually fetch dependencies
// or assign private fields
db := bind.Get[Database](ctx)

// this returns a new instance of *dao with all fields
// that contain the bind tag populated
dao := bind.New[*dao](ctx)
```

#### API
The API allows for various setups. Note that a type is considered a leaf if there is no
other mapping for that type. Hence if `bind.Implementation[X, Y]` and both `bind.Instance[Y]` have
been configured the instance binding is the result of the injection of `X`.

- `bind.Configure(ctx, bindings...)`: configure bindings in a context; can overwrite existing bindings of the parent context
- `bind.Implementation[X, Y]()`: bind `Y` for `X`, return instances of `Y` if `Y` is a leaf
- `bind.Once[X]()`: bind `X` for exactly one instance
- `bind.ImplementationOnce[X, Y]()`: bind exactly one instance of `Y` for `X`
- `bind.Instance[X](inst X)`: bind `X` to `inst`
- `bind.Many[X]()`: bind `X` and return instances of `X`
- `bind.Provider[X](f func() (X, error))`: bind `X` to invocations of `f`
- `bind.New[X](ctx)`: resolve `X` or create a new instance of `X` (X doesn't need to be bound)
- `bind.Get[X](ctx)`: resolve `X`
- `bind.For[X](ctx, scope)`: resolve `X` for `scope`
- `bind.MaybeNew[X](ctx)`: resolve `X` or create a new instance of `X`; return error instead of panic
- `bind.MaybeGet[X](ctx)`: resolve `X`; return error instead of panic
- `bind.MaybeFor[X](ctx, scope)`: resolve `X` for `scope`; return error instead of panic
- `bind.Initializer`: When implemented, calls `InitAfter` after a type was initialized

#### Type-Safety
`bind.Implementation[Iface, Impl]()` can't guarantee `Impl` is assignable to `Iface` at compile time and panics at runtime.
Internally there are several instances of `any` and reflection is still used given the nature of how Go generics
work. 

Note though that we can use generics to ensure that certain types exist and the use of `any` in the
public API is non-existent.

### Package `channel`
Utilities to work with channels.

- `Safe*` methods to perform common actions on channels that won't panic (read: either you don't care or it's a code smell)
- `Maybe*` methods to perform common patterns with less ceremony

### Package `slice`
Utilities to work with slices.

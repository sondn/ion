// Copyright (c) 2017 The Ion Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Ion nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
Package ion provides a beautifully expressive and easy to use foundation for your next website, API, or distributed app.

Source code and other details for the project are available at GitHub:

   https://github.com/get-ion/ion

Installation

The only requirement is the Go Programming Language, at least version 1.8.x

    $ go get -u github.com/get-ion/ion


Example code:


    package main

    import (
        "github.com/get-ion/ion"
        "github.com/get-ion/ion/context"
        "github.com/get-ion/ion/view"
    )

    // User is just a bindable object structure.
    type User struct {
        Username  string `json:"username"`
        Firstname string `json:"firstname"`
        Lastname  string `json:"lastname"`
        City      string `json:"city"`
        Age       int    `json:"age"`
    }

    func main() {
        app := ion.New()

        // Define templates using the std html/template engine.
        // Parse and load all files inside "./views" folder with ".html" file extension.
        // Reload the templates on each request (development mode).
        app.RegisterView(view.HTML("./views", ".html").Reload(true))

        // Regster custom handler for specific http errors.
        app.OnErrorCode(ion.StatusInternalServerError, func(ctx context.Context) {
            // .Values are used to communicate between handlers, middleware.
            errMessage := ctx.Values().GetString("error")
            if errMessage != "" {
                ctx.Writef("Internal server error: %s", errMessage)
                return
            }

            ctx.Writef("(Unexpected) internal server error")
        })

        app.Use(func(ctx context.Context) {
            ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
            ctx.Next()
        })

        // app.Done(func(ctx context.Context) {})

        // Method POST: http://localhost:8080/decode
        app.Post("/decode", func(ctx context.Context) {
            var user User
            ctx.ReadJSON(&user)
            ctx.Writef("%s %s is %d years old and comes from %s", user.Firstname, user.Lastname, user.Age, user.City)
        })

        // Method GET: http://localhost:8080/encode
        app.Get("/encode", func(ctx context.Context) {
            doe := User{
                Username:  "Johndoe",
                Firstname: "John",
                Lastname:  "Doe",
                City:      "Neither FBI knows!!!",
                Age:       25,
            }

            ctx.JSON(doe)
        })

        // Method GET: http://localhost:8080/profile/anytypeofstring
        app.Get("/profile/{username:string}", profileByUsername)

        usersRoutes := app.Party("/users", logThisMiddleware)
        {
            // Method GET: http://localhost:8080/users/42
            usersRoutes.Get("/{id:int min(1)}", getUserByID)
            // Method POST: http://localhost:8080/users/create
            usersRoutes.Post("/create", createUser)
        }

        // Listen for incoming HTTP/1.x & HTTP/2 clients on localhost port 8080.
        app.Run(ion.Addr(":8080"), ion.WithCharset("UTF-8"))
    }

    func logThisMiddleware(ctx context.Context) {
        ctx.Application().Logger().Infof("Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())

        // .Next is required to move forward to the chain of handlers,
        // if missing then it stops the execution at this handler.
        ctx.Next()
    }

    func profileByUsername(ctx context.Context) {
        // .Params are used to get dynamic path parameters.
        username := ctx.Params().Get("username")
        ctx.ViewData("Username", username)
        // renders "./views/users/profile.html"
        // with {{ .Username }} equals to the username dynamic path parameter.
        ctx.View("users/profile.html")
    }

    func getUserByID(ctx context.Context) {
        userID := ctx.Params().Get("id") // Or convert directly using: .Values().GetInt/GetInt64 etc...
        // your own db fetch here instead of user :=...
        user := User{Username: "username" + userID}

        ctx.XML(user)
    }

    func createUser(ctx context.Context) {
        var user User
        err := ctx.ReadForm(&user)
        if err != nil {
            ctx.Values().Set("error", "creating user, read and parse form failed. "+err.Error())
            ctx.StatusCode(ion.StatusInternalServerError)
            return
        }
        // renders "./views/users/create_verification.html"
        // with {{ . }} equals to the User object, i.e {{ .Username }} , {{ .Firstname}} etc...
        ctx.ViewData("", user)
        ctx.View("users/create_verification.html")
    }


Routing

All HTTP methods are supported, developers can also register handlers for same paths for different methods.
The first parameter is the HTTP Method,
second parameter is the request path of the route,
third variadic parameter should contains one or more context.Handler executed
by the registered order when a user requests for that specific resouce path from the server.

Example code:


    app := ion.New()

    app.Handle("GET", "/contact", func(ctx context.Context){
        ctx.HTML("<h1> Hello from /contact </h1>")
    })


In order to make things easier for the user, ion provides functions for all HTTP Methods.
The first parameter is the request path of the route,
second variadic parameter should contains one or more context.Handler executed
by the registered order when a user requests for that specific resouce path from the server.

Example code:


    app := ion.New()

    // Method: "GET"
    app.Get("/", handler)

    // Method: "POST"
    app.Post("/", handler)

    // Method: "PUT"
    app.Put("/", handler)

    // Method: "DELETE"
    app.Delete("/", handler)

    // Method: "OPTIONS"
    app.Options("/", handler)

    // Method: "TRACE"
    app.Trace("/", handler)

    // Method: "CONNECT"
    app.Connect("/", handler)

    // Method: "HEAD"
    app.Head("/", handler)

    // Method: "PATCH"
    app.Patch("/", handler)

    // register the route for all HTTP Methods
    app.Any("/", handler)

    func handler(ctx context.Context){
    ctx.Writef("Hello from method: %s and path: %s", ctx.Method(), ctx.Path())
    }



Grouping Routes


A set of routes that are being groupped by path prefix can (optionally) share the same middleware handlers and template layout.
A group can have a nested group too.

`.Party` is being used to group routes, developers can declare an unlimited number of (nested) groups.


Example code:


    users:= app.Party("/users", myAuthHandler)

    // http://myhost.com/users/42/profile
    users.Get("/{userid:int}/profile", userProfileHandler)
    // http://myhost.com/users/messages/1
    users.Get("/inbox/{messageid:int}", userMessageHandler)

    app.Run(ion.Addr("myhost.com:80"))



Custom HTTP Errors


ion developers are able to register their own handlers for http statuses like 404 not found, 500 internal server error and so on.

Example code:


    // when 404 then render the template $templatedir/errors/404.html
    app.OnErrorCode(ion.StatusNotFound, func(ctx context.Context){
        ctx.View("errors/404.html")
    })

    app.OnErrorCode(500, func(ctx context.Context){
        // ...
    })

Basic HTTP API

With the help of ion's expressionist router you can build any form of API you desire, with
safety.

Example code:


    package main

    import (
        "github.com/get-ion/ion"
        "github.com/get-ion/ion/context"
    )

    func main() {
        app := ion.New()

        // registers a custom handler for 404 not found http (error) status code,
        // fires when route not found or manually by ctx.StatusCode(ion.StatusNotFound).
        app.OnErrorCode(ion.StatusNotFound, notFoundHandler)

        // GET -> HTTP Method
        // / -> Path
        // func(ctx context.Context) -> The route's handler.
        //
        // Third receiver should contains the route's handler(s), they are executed by order.
        app.Handle("GET", "/", func(ctx context.Context) {
            // navigate to the middle of $GOPATH/src/github.com/get-ion/ion/context/context.go
            // to overview all context's method (there a lot of them, read that and you will learn how ion works too)
            ctx.HTML("Hello from " + ctx.Path()) // Hello from /
        })

        app.Get("/home", func(ctx context.Context) {
            ctx.Writef(`Same as app.Handle("GET", "/", [...])`)
        })

        app.Get("/donate", donateHandler, donateFinishHandler)

        // Pssst, don't forget dynamic-path example for more "magic"!
        app.Get("/api/users/{userid:int min(1)}", func(ctx context.Context) {
            userID, err := ctx.Params().GetInt("userid")

            if err != nil {
                ctx.Writef("error while trying to parse userid parameter," +
                    "this will never happen if :int is being used because if it's not integer it will fire Not Found automatically.")
                ctx.StatusCode(ion.StatusBadRequest)
                return
            }

            ctx.JSON(map[string]interface{}{
                // you can pass any custom structured go value of course.
                "user_id": userID,
            })
        })
        // app.Post("/", func(ctx context.Context){}) -> for POST http method.
        // app.Put("/", func(ctx context.Context){})-> for "PUT" http method.
        // app.Delete("/", func(ctx context.Context){})-> for "DELETE" http method.
        // app.Options("/", func(ctx context.Context){})-> for "OPTIONS" http method.
        // app.Trace("/", func(ctx context.Context){})-> for "TRACE" http method.
        // app.Head("/", func(ctx context.Context){})-> for "HEAD" http method.
        // app.Connect("/", func(ctx context.Context){})-> for "CONNECT" http method.
        // app.Patch("/", func(ctx context.Context){})-> for "PATCH" http method.
        // app.Any("/", func(ctx context.Context){}) for all http methods.

        // More than one route can contain the same path with a different http mapped method.
        // You can catch any route creation errors with:
        // route, err := app.Get(...)
        // set a name to a route: route.Name = "myroute"

        // You can also group routes by path prefix, sharing middleware(s) and done handlers.

        adminRoutes := app.Party("/admin", adminMiddleware)

        adminRoutes.Done(func(ctx context.Context) { // executes always last if ctx.Next()
            ctx.Application().Logger().Infof("response sent to " + ctx.Path())
        })
        // adminRoutes.Layout("/views/layouts/admin.html") // set a view layout for these routes, see more at view examples.

        // GET: http://localhost:8080/admin
        adminRoutes.Get("/", func(ctx context.Context) {
            // [...]
            ctx.StatusCode(ion.StatusOK) // default is 200 == ion.StatusOK
            ctx.HTML("<h1>Hello from admin/</h1>")

            ctx.Next() // in order to execute the party's "Done" Handler(s)
        })

        // GET: http://localhost:8080/admin/login
        adminRoutes.Get("/login", func(ctx context.Context) {
            // [...]
        })
        // POST: http://localhost:8080/admin/login
        adminRoutes.Post("/login", func(ctx context.Context) {
            // [...]
        })

        // subdomains, easier than ever, should add localhost or 127.0.0.1 into your hosts file,
        // etc/hosts on unix or C:/windows/system32/drivers/etc/hosts on windows.
        v1 := app.Party("v1.")
        { // braces are optional, it's just type of style, to group the routes visually.

            // http://v1.localhost:8080
            v1.Get("/", func(ctx context.Context) {
                ctx.HTML("Version 1 API. go to <a href='" + ctx.Path() + "/api" + "'>/api/users</a>")
            })

            usersAPI := v1.Party("/api/users")
            {
                // http://v1.localhost:8080/api/users
                usersAPI.Get("/", func(ctx context.Context) {
                    ctx.Writef("All users")
                })
                // http://v1.localhost:8080/api/users/42
                usersAPI.Get("/{userid:int}", func(ctx context.Context) {
                    ctx.Writef("user with id: %s", ctx.Params().Get("userid"))
                })
            }
        }

        // wildcard subdomains.
        wildcardSubdomain := app.Party("*.")
        {
            wildcardSubdomain.Get("/", func(ctx context.Context) {
                ctx.Writef("Subdomain can be anything, now you're here from: %s", ctx.Subdomain())
            })
        }

        // http://localhost:8080
        // http://localhost:8080/home
        // http://localhost:8080/donate
        // http://localhost:8080/api/users/42
        // http://localhost:8080/admin
        // http://localhost:8080/admin/login
        //
        // http://localhost:8080/api/users/0
        // http://localhost:8080/api/users/blabla
        // http://localhost:8080/wontfound
        //
        // if hosts edited:
        //  http://v1.localhost:8080
        //  http://v1.localhost:8080/api/users
        //  http://v1.localhost:8080/api/users/42
        //  http://anything.localhost:8080
        app.Run(ion.Addr(":8080"))
    }

    func adminMiddleware(ctx context.Context) {
        // [...]
        ctx.Next() // to move to the next handler, or don't that if you have any auth logic.
    }

    func donateHandler(ctx context.Context) {
        ctx.Writef("Just like an inline handler, but it can be " +
            "used by other package, anywhere in your project.")

        // let's pass a value to the next handler
        // Values is the way handlers(or middleware) are communicating between each other.
        ctx.Values().Set("donate_url", "https://github.com/get-ion/ion#buy-me-a-cup-of-coffee")
        ctx.Next() // in order to execute the next handler in the chain, look donate route.
    }

    func donateFinishHandler(ctx context.Context) {
        // values can be any type of object so we could cast the value to a string
        // but ion provides an easy to do that, if donate_url is not defined, then it returns an empty string instead.
        donateURL := ctx.Values().GetString("donate_url")
        ctx.Application().Logger().Infof("donate_url value was: " + donateURL)
        ctx.Writef("\n\nDonate sent(?).")
    }

    func notFoundHandler(ctx context.Context) {
        ctx.HTML("Custom route for 404 not found http code, here you can render a view, html, json <b>any valid response</b>.")
    }


Parameterized Path

At the previous example,
we've seen static routes, group of routes, subdomains, wildcard subdomains, a small example of parameterized path
with a single known paramete and custom http errors, now it's time to see wildcard parameters and macros.

ion, like net/http std package registers route's handlers
by a Handler, the ion' type of handler is just a func(ctx context.Context)
where context comes from github.com/get-ion/ion/context.
Until go 1.9 you will have to import that package too, after go 1.9 this will be not be necessary.

ion has the easiest and the most powerful routing process you have ever meet.

At the same time,
ion has its own interpeter(yes like a programming language)
for route's path syntax and their dynamic path parameters parsing and evaluation,
I am calling them "macros" for shortcut.
How? It calculates its needs and if not any special regexp needed then it just
registers the route with the low-level path syntax,
otherwise it pre-compiles the regexp and adds the necessary middleware(s).

Standard macro types for parameters:

    +------------------------+
    | {param:string}         |
    +------------------------+
    string type
    anything

    +------------------------+
    | {param:int}            |
    +------------------------+
    int type
    only numbers (0-9)

    +------------------------+
    | {param:alphabetical}   |
    +------------------------+
    alphabetical/letter type
    letters only (upper or lowercase)

    +------------------------+
    | {param:file}           |
    +------------------------+
    file type
    letters (upper or lowercase)
    numbers (0-9)
    underscore (_)
    dash (-)
    point (.)
    no spaces ! or other character

    +------------------------+
    | {param:path}           |
    +------------------------+
    path type
    anything, should be the last part, more than one path segment,
    i.e: /path1/path2/path3 , ctx.Params().GetString("param") == "/path1/path2/path3"

if type is missing then parameter's type is defaulted to string, so
{param} == {param:string}.

If a function not found on that type then the "string"'s types functions are being used.
i.e:


    {param:int min(3)}


Besides the fact that ion provides the basic types and some default "macro funcs"
you are able to register your own too!.

Register a named path parameter function:


    app.Macros().Int.RegisterFunc("min", func(argument int) func(paramValue string) bool {
        [...]
        return true/false -> true means valid.
    })

at the func(argument ...) you can have any standard type, it will be validated before the server starts
so don't care about performance here, the only thing it runs at serve time is the returning func(paramValue string) bool.

    {param:string equal(ion)} , "ion" will be the argument here:
    app.Macros().String.RegisterFunc("equal", func(argument string) func(paramValue string) bool {
        return func(paramValue string){ return argument == paramValue }
    })


Example code:


	// you can use the "string" type which is valid for a single path parameter that can be anything.
	app.Get("/username/{name}", func(ctx context.Context) {
		ctx.Writef("Hello %s", ctx.Params().Get("name"))
	}) // type is missing = {name:string}

	// Let's register our first macro attached to int macro type.
	// "min" = the function
	// "minValue" = the argument of the function
	// func(string) bool = the macro's path parameter evaluator, this executes in serve time when
	// a user requests a path which contains the :int macro type with the min(...) macro parameter function.
	app.Macros().Int.RegisterFunc("min", func(minValue int) func(string) bool {
		// do anything before serve here [...]
		// at this case we don't need to do anything
		return func(paramValue string) bool {
			n, err := strconv.Atoi(paramValue)
			if err != nil {
				return false
			}
			return n >= minValue
		}
	})

	// http://localhost:8080/profile/id>=1
	// this will throw 404 even if it's found as route on : /profile/0, /profile/blabla, /profile/-1
	// macro parameter functions are optional of course.
	app.Get("/profile/{id:int min(1)}", func(ctx context.Context) {
		// second parameter is the error but it will always nil because we use macros,
		// the validaton already happened.
		id, _ := ctx.Params().GetInt("id")
		ctx.Writef("Hello id: %d", id)
	})

	// to change the error code per route's macro evaluator:
	app.Get("/profile/{id:int min(1)}/friends/{friendid:int min(1) else 504}", func(ctx context.Context) {
		id, _ := ctx.Params().GetInt("id")
		friendid, _ := ctx.Params().GetInt("friendid")
		ctx.Writef("Hello id: %d looking for friend id: ", id, friendid)
	}) // this will throw e 504 error code instead of 404 if all route's macros not passed.

	// http://localhost:8080/game/a-zA-Z/level/0-9
	// remember, alphabetical is lowercase or uppercase letters only.
	app.Get("/game/{name:alphabetical}/level/{level:int}", func(ctx context.Context) {
		ctx.Writef("name: %s | level: %s", ctx.Params().Get("name"), ctx.Params().Get("level"))
	})

	// let's use a trivial custom regexp that validates a single path parameter
	// which its value is only lowercase letters.

	// http://localhost:8080/lowercase/anylowercase
	app.Get("/lowercase/{name:string regexp(^[a-z]+)}", func(ctx context.Context) {
		ctx.Writef("name should be only lowercase, otherwise this handler will never executed: %s", ctx.Params().Get("name"))
	})

	// http://localhost:8080/single_file/app.js
	app.Get("/single_file/{myfile:file}", func(ctx context.Context) {
		ctx.Writef("file type validates if the parameter value has a form of a file name, got: %s", ctx.Params().Get("myfile"))
	})

	// http://localhost:8080/myfiles/any/directory/here/
	// this is the only macro type that accepts any number of path segments.
	app.Get("/myfiles/{directory:path}", func(ctx context.Context) {
		ctx.Writef("path type accepts any number of path segments, path after /myfiles/ is: %s", ctx.Params().Get("directory"))
	})

    // for wildcard path (any number of path segments) without validation you can use:
	// /myfiles/*directory
	// "{param}"'s performance is exactly the same of ":param"'s.

	// alternatives -> ":param" for single path parameter and "*paramPath" for wildcard path parameter
	// acquire them by ctx.Params().Get as always.

	if err := app.Run(ion.Addr(":8080")); err != nil {
		panic(err)
	}
}



A path parameter name should contain only alphabetical letters, symbols, containing '_' and numbers are NOT allowed.
If route failed to be registered, the app will panic without any warnings
if you didn't catch the second return value(error) on .Handle/.Get....

Last, do not confuse ctx.Values() with ctx.Params().
Path parameter's values goes to ctx.Params() and context's local storage
that can be used to communicate between handlers and middleware(s) goes to
ctx.Values(), path parameters and the rest of any custom values are separated for your own good.

Run

  $ go run main.go



Static Files


    // StaticServe serves a directory as web resource
    // it's the simpliest form of the Static* functions
    // Almost same usage as StaticWeb
    // accepts only one required parameter which is the systemPath,
    // the same path will be used to register the GET and HEAD method routes.
    // If second parameter is empty, otherwise the requestPath is the second parameter
    // it uses gzip compression (compression on each request, no file cache).
    //
    // Returns the GET *Route.
    StaticServe(systemPath string, requestPath ...string) (*Route, error)

    // StaticContent registers a GET and HEAD method routes to the requestPath
    // that are ready to serve raw static bytes, memory cached.
    //
    // Returns the GET *Route.
    StaticContent(reqPath string, cType string, content []byte) (*Route, error)

    // StaticEmbedded  used when files are distributed inside the app executable, using go-bindata mostly
    // First parameter is the request path, the path which the files in the vdir will be served to, for example "/static"
    // Second parameter is the (virtual) directory path, for example "./assets"
    // Third parameter is the Asset function
    // Forth parameter is the AssetNames function.
    //
    // Returns the GET *Route.
    //
    // Example: https://github.com/get-ion/ion/tree/master/_examples/file-server/embedding-files-into-app
    StaticEmbedded(requestPath string, vdir string, assetFn func(name string) ([]byte, error), namesFn func() []string) (*Route, error)

    // Favicon serves static favicon
    // accepts 2 parameters, second is optional
    // favPath (string), declare the system directory path of the __.ico
    // requestPath (string), it's the route's path, by default this is the "/favicon.ico" because some browsers tries to get this by default first,
    // you can declare your own path if you have more than one favicon (desktop, mobile and so on)
    //
    // this func will add a route for you which will static serve the /yuorpath/yourfile.ico to the /yourfile.ico
    // (nothing special that you can't handle by yourself).
    // Note that you have to call it on every favicon you have to serve automatically (desktop, mobile and so on).
    //
    // Returns the GET *Route.
    Favicon(favPath string, requestPath ...string) (*Route, error)

    // StaticWeb returns a handler that serves HTTP requests
    // with the contents of the file system rooted at directory.
    //
    // first parameter: the route path
    // second parameter: the system directory
    // third OPTIONAL parameter: the exception routes
    //      (= give priority to these routes instead of the static handler)
    // for more options look app.StaticHandler.
    //
    //     app.StaticWeb("/static", "./static")
    //
    // As a special case, the returned file server redirects any request
    // ending in "/index.html" to the same path, without the final
    // "index.html".
    //
    // StaticWeb calls the StaticHandler(systemPath, listingDirectories: false, gzip: false ).
    //
    // Returns the GET *Route.
    StaticWeb(requestPath string, systemPath string, exceptRoutes ...*Route) (*Route, error)


Example code:


    package main

    import (
        "github.com/get-ion/ion"
        "github.com/get-ion/ion/context"
    )

    func main() {
        app := ion.New()

        // This will serve the ./static/favicons/ion_32_32.ico to: localhost:8080/favicon.ico
        app.Favicon("./static/favicons/ion_32_32.ico")

        // app.Favicon("./static/favicons/ion_32_32.ico", "/favicon_48_48.ico")
        // This will serve the ./static/favicons/ion_32_32.ico to: localhost:8080/favicon_48_48.ico

        app.Get("/", func(ctx context.Context) {
            ctx.HTML(`<a href="/favicon.ico"> press here to see the favicon.ico</a>.
            At some browsers like chrome, it should be visible at the top-left side of the browser's window,
            because some browsers make requests to the /favicon.ico automatically,
            so ion serves your favicon in that path too (you can change it).`)
        }) // if favicon doesn't show to you, try to clear your browser's cache.

        app.Run(ion.Addr(":8080"))
    }

More examples can be found here: https://github.com/get-ion/ion/tree/master/_examples/beginner/file-server


Middleware Ecosystem

Middleware is just a concept of ordered chain of handlers.
Middleware can be registered globally, per-party, per-subdomain and per-route.


Example code:

      // globally
      // before any routes, appends the middleware to all routes
      app.Use(func(ctx context.Context){
         // ... any code here

         ctx.Next() // in order to continue to the next handler,
         // if that is missing then the next in chain handlers will be not executed,
         // useful for authentication middleware
      })

      // globally
      // after or before any routes, prepends the middleware to all routes
      app.UseGlobal(handler1, handler2, handler3)

      // per-route
      app.Post("/login", authenticationHandler, loginPageHandler)

      // per-party(group of routes)
      users := app.Party("/users", usersMiddleware)
      users.Get("/", usersIndex)

      // per-subdomain
      mysubdomain := app.Party("mysubdomain.", firstMiddleware)
      mysubdomain.Use(secondMiddleware)
      mysubdomain.Get("/", mysubdomainIndex)

      // per wildcard, dynamic subdomain
      dynamicSub := app.Party(".*", firstMiddleware, secondMiddleware)
      dynamicSub.Get("/", func(ctx context.Context){
        ctx.Writef("Hello from subdomain: "+ ctx.Subdomain())
      })


ion is able to wrap and convert any external, third-party Handler you used to use to your web application.
Let's convert the https://github.com/rs/cors net/http external middleware which returns a `next form` handler.


Example code:

    package main

    import (
        "github.com/rs/cors"

        "github.com/get-ion/ion"
        "github.com/get-ion/ion/context"
    )

    func main() {

        app := ion.New()
        corsOptions := cors.Options{
            AllowedOrigins:   []string{"*"},
            AllowCredentials: true,
        }

        corsWrapper := cors.New(corsOptions).ServeHTTP

        app.WrapRouter(corsWrapper)

        v1 := app.Party("/api/v1")
        {
            v1.Get("/", h)
            v1.Put("/put", h)
            v1.Post("/post", h)
        }

        app.Run(ion.Addr(":8080"))
    }

    func h(ctx context.Context) {
        ctx.Application().Logger().Infof(ctx.Path())
        ctx.Writef("Hello from %s", ctx.Path())
    }


View Engine


ion supports 5 template engines out-of-the-box, developers can still use any external golang template engine,
as `context.ResponseWriter()` is an `io.Writer`.

All of these five template engines have common features with common API,
like Layout, Template Funcs, Party-specific layout, partial rendering and more.

      The standard html,
      its template parser is the golang.org/pkg/html/template/.

      Django,
      its template parser is the github.com/flosch/pongo2

      Pug(Jade),
      its template parser is the github.com/Joker/jade

      Handlebars,
      its template parser is the github.com/aymerick/raymond

      Amber,
      its template parser is the github.com/eknkc/amber


Example code:

    package main

    import (
        "github.com/get-ion/ion"
        "github.com/get-ion/ion/context"
        "github.com/get-ion/ion/view"
    )

    func main() {
        app := ion.New() // defaults to these

        // - standard html  | view.HTML(...)
        // - django         | view.Django(...)
        // - pug(jade)      | view.Pug(...)
        // - handlebars     | view.Handlebars(...)
        // - amber          | view.Amber(...)

        tmpl := view.HTML("./templates", ".html")
        tmpl.Reload(true) // reload templates on each request (development mode)
        // default template funcs are:
        //
        // - {{ urlpath "mynamedroute" "pathParameter_ifneeded" }}
        // - {{ render "header.html" }}
        // - {{ render_r "header.html" }} // partial relative path to current page
        // - {{ yield }}
        // - {{ current }}
        tmpl.AddFunc("greet", func(s string) string {
            return "Greetings " + s + "!"
        })
        app.RegisterView(tmpl)

        app.Get("/", hi)

        // http://localhost:8080
        app.Run(ion.Addr(":8080"), ion.WithCharset("UTF-8")) // defaults to that but you can change it.
    }

    func hi(ctx context.Context) {
        ctx.ViewData("Title", "Hi Page")
        ctx.ViewData("Name", "ion") // {{.Name}} will render: ion
        // ctx.ViewData("", myCcustomStruct{})
        ctx.View("hi.html")
    }



View engine supports bundled(https://github.com/jteeuwen/go-bindata) template files too.
go-bindata gives you two functions, asset and assetNames,
these can be setted to each of the template engines using the `.Binary` func.

Example code:

    package main

    import (
        "github.com/get-ion/ion"
        "github.com/get-ion/ion/context"
        "github.com/get-ion/ion/view"
    )

    func main() {
        app := ion.New()
        // $ go get -u github.com/jteeuwen/go-bindata/...
        // $ go-bindata ./templates/...
        // $ go build
        // $ ./embedding-templates-into-app
        // html files are not used, you can delete the folder and run the example
        app.RegisterView(view.HTML("./templates", ".html").Binary(Asset, AssetNames))
        app.Get("/", hi)

        // http://localhost:8080
        app.Run(ion.Addr(":8080"))
    }

    type page struct {
        Title, Name string
    }

    func hi(ctx context.Context) {
        ctx.ViewData("", page{Title: "Hi Page", Name: "ion"})
        ctx.View("hi.html")
    }


A real example can be found here: https://github.com/get-ion/ion/tree/master/_examples/view/embedding-templates-into-app.

Enable auto-reloading of templates on each request. Useful while developers are in dev mode
as they no neeed to restart their app on every template edit.

Example code:


    pugEngine := view.Pug("./templates", ".jade")
    pugEngine.Reload(true) // <--- set to true to re-build the templates on each request.
    app.RegisterView(pugEngine)


Each one of these template engines has different options located here: https://github.com/get-ion/ion/tree/master/view .


Sessions


This example will show how to store and access data from a session.

You don’t need any third-party library,
but If you want you can use any session manager compatible or not.

In this example we will only allow authenticated users to view our secret message on the /secret page.
To get access to it, the will first have to visit /login to get a valid session cookie,
which logs him in. Additionally he can visit /logout to revoke his access to our secret message.


Example code:


    // main.go
    package main

    import (
        "github.com/get-ion/ion"
        "github.com/get-ion/ion/context"

        "github.com/get-ion/sessions"
    )

    var (
        cookieNameForSessionID = "mycookiesessionnameid"
        sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
    )

    func secret(ctx context.Context) {

        // Check if user is authenticated
        if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
            ctx.StatusCode(ion.StatusForbidden)
            return
        }

        // Print secret message
        ctx.WriteString("The cake is a lie!")
    }

    func login(ctx context.Context) {
        session := sess.Start(ctx)

        // Authentication goes here
        // ...

        // Set user as authenticated
        session.Set("authenticated", true)
    }

    func logout(ctx context.Context) {
        session := sess.Start(ctx)

        // Revoke users authentication
        session.Set("authenticated", false)
    }

    func main() {
        app := ion.New()

        app.Get("/secret", secret)
        app.Get("/login", login)
        app.Get("/logout", logout)

        app.Run(ion.Addr(":8080"))
    }


Running the example:


    $ go get github.com/get-ion/sessions
    $ go run main.go

    $ curl -s http://localhost:8080/secret
    Forbidden

    $ curl -s -I http://localhost:8080/login
    Set-Cookie: mycookiesessionnameid=MTQ4NzE5Mz...

    $ curl -s --cookie "mycookiesessionnameid=MTQ4NzE5Mz..." http://localhost:8080/secret
    The cake is a lie!

More examples:

    https://github.com/get-ion/sessions


Websockets

In this example we will create a small chat between web sockets via browser.

Example Server Code:

    // main.go
    package main

    import (
        "fmt"

        "github.com/get-ion/ion"
        "github.com/get-ion/ion/context"

        "github.com/get-ion/websocket"
    )

    func main() {
        app := ion.New()

        app.Get("/", func(ctx context.Context) {
            ctx.ServeFile("websockets.html", false) // second parameter: enable gzip?
        })

        setupWebsocket(app)

        // x2
        // http://localhost:8080
        // http://localhost:8080
        // write something, press submit, see the result.
        app.Run(ion.Addr(":8080"))
    }

    func setupWebsocket(app *ion.Application) {
        // create our echo websocket server
        ws := websocket.New(websocket.Config{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        })
        ws.OnConnection(handleConnection)

        // register the server on an endpoint.
        // see the inline javascript code i the websockets.html, this endpoint is used to connect to the server.
        app.Get("/echo", ws.Handler())

        // serve the javascript built'n client-side library,
        // see weboskcets.html script tags, this path is used.
        app.Any("/ion-ws.js", func(ctx context.Context) {
            ctx.Write(websocket.ClientSource)
        })
    }

    func handleConnection(c websocket.Connection) {
        // Read events from browser
        c.On("chat", func(msg string) {
            // Print the message to the console, c.Context() is the ion's http context.
            fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
            // Write message back to the client message owner:
            // c.Emit("chat", msg)
            c.To(websocket.Broadcast).Emit("chat", msg)
        })
    }

Example Client(javascript) Code:

    <!-- websockets.html -->
    <input id="input" type="text" />
    <button onclick="send()">Send</button>
    <pre id="output"></pre>
    <script src="/ion-ws.js"></script>
    <script>
        var input = document.getElementById("input");
        var output = document.getElementById("output");

        // Ws comes from the auto-served '/ion-ws.js'
        var socket = new Ws("ws://localhost:8080/echo");
        socket.OnConnect(function () {
            output.innerHTML += "Status: Connected\n";
        });

        socket.OnDisconnect(function () {
            output.innerHTML += "Status: Disconnected\n";
        });

        // read events from the server
        socket.On("chat", function (msg) {
            addMessage(msg)
        });

        function send() {
            addMessage("Me: "+input.value) // write ourselves
            socket.Emit("chat", input.value);// send chat event data to the websocket server
            input.value = ""; // clear the input
        }

        function addMessage(msg) {
            output.innerHTML += msg + "\n";
        }
    </script>


Running the example:


    $ go get github.com/get-ion/websocket
    $ go run main.go
    $ start http://localhost:8080


That's the basics

But you should have a basic idea of the framework by now, we just scratched the surface.
If you enjoy what you just saw and want to learn more, please follow the below links:

Examples:

    https://github.com/get-ion/ion/tree/master/_examples

Built'n Middleware:

    https://github.com/get-ion/ion/tree/master/middleware

Home Page:

    http://github.com/get-ion/ion

*/
package ion

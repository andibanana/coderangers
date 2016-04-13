0. Install [Git](https://git-scm.com/).

1. Install the [Go compiler and related tools](https://golang.org).

2. Set your `GOPATH` environment variable to an empty folder where you want 
   your Go code and libraries to be stored.

3. Clone the CodRangers repo into the src folder of your `GOPATH`.
     - `git clone https://github.com/andibanana/coderangers`

4. Get the dependencies.
     - `go get "github.com/gorilla/sessions"`
     - `go get "github.com/go-sql-driver/mysql"`
     - `go get "github.com/mattn/go-sqlite3"`

5. Install [Node.js](https://nodejs.org/).

6. Clone UVA-NODE into the root directory (C:\\).
     - `git clone https://github.com/lucastan/uva-node`

7. Get the dependencies of UVA-NODE.
     - `npm install` (inside the uva-node folder)

8. To run the server, run `go build server.go && server` in the coderangers 
   folder. (In the future, this will become `go install coderangers`)


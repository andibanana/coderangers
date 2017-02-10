0. Install [Git](https://git-scm.com/).

1. Install the [Go compiler and related tools](https://golang.org).

2. Set your `GOPATH` environment variable to an empty folder where you want 
   your Go code and libraries to be stored.

3. Clone the CodeRangers repo into the src folder of your `GOPATH`.
     - `git clone https://github.com/andibanana/coderangers`

4. Get the dependencies.
     - `go get "github.com/gorilla/sessions"`
     - `go get "github.com/go-sql-driver/mysql"`
     - `go get "github.com/mattn/go-sqlite3"`

5. Install [Isolate](https://github.com/ioi/isolate) for C support.

6. Install [JDK](http://www.oracle.com/technetwork/java/javase/downloads/index.html) for Java support.

7. Install GCC also for C support.

8. Json files in the root folder are config files and should be updated as needed.

9. To run the server, run `go build server.go && server` in the coderangers 
   folder. (In the future, this will become `go install coderangers`)


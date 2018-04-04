1. Install [Git](https://git-scm.com/).
2. Install the [Go compiler and related tools](https://golang.org).
3. Set your `GOPATH` environment variable to an empty folder where you want 
   your Go code and libraries to be stored.
4. Clone the CodeRangers repo into the src folder of your `GOPATH`.
     - `git clone https://github.com/andibanana/coderangers`
5. Get the dependencies.
     - `go get "github.com/gorilla/sessions"`
     - `go get "github.com/go-sql-driver/mysql"`
     - `go get "github.com/mattn/go-sqlite3"`
	 - `go get "golang.org/x/crypto/acme/autocert"`
6. Install [Isolate](https://github.com/ioi/isolate) for C support.
7. Install [JDK](http://www.oracle.com/technetwork/java/javase/downloads/index.html) for Java support.
7. Install GCC also for C support.
8. Json files in the root folder are config files and should be updated as needed.
9. To run the server, run `go build server.go && server` in the coderangers 
   folder. (In the future, this will become `go install coderangers`)


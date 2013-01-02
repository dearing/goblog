blog
====
**still in development**

blogging with go, markdown and redis

about
-----

Simply this little gem takes some templates and some articles written in markdown
and parses and then pushes them to an active redis db where they are served.  A simple
filewatcher then pushes modified documents as they are created or modified.


try it out
---
**Templates at this time have hardcoded names in the source so it will fail until I provide some basic ones for testing.**
```
go get github.com/dearing/blog
blog --help

Usage of blog:
  -host=":8080": host to bind to
  -rdb=-1: redis db index
  -rh="localhost:6379": redis host
  -root="wwwroot": webserver document root folder
  -rp="": redis password
  -verbose=false: log common operations and not just errors
```

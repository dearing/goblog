blog
====
**still in development**

blogging with go, markdown and redis

about
-----
Wanting a simple blogging engine that kept things simple but performant I got to writing one for myself.  Markdown presented itself as a nice format to store posts in and better yet, write posts in on a desktop or phone.  I also wanted this blogging engine to keep metrics and update information on itself that can be presented to an authenticated reader (or no) internally.  Finally I wanted all this to be adminstrated from a OpenID authentication scheme so that other services can reliably handle credentials and what not.  From all this was this project born.


incoming android
----
In the future, I will jazz up a basic android app that can handle all the administration relevant to the blog remotely so that I can post on the go.  This is important to me because I am otherwise too lazy to share my whims, but with an app in hand this may change...


try it out
---
```
go get github.com/dearing/blog
go install

cd example
blog --help

Usage of blog:
  -articles="articles": markdown posts
  -redis-db=-1: redis db index
  -redis-host="localhost:6379": redis host
  -redis-pass="": redis password
  -suffix=".md": filtered extension
  -templates="templates": templates posts
  -verbose=false: log common operations and not just errors
  -wwwhost=":8080": host to bind to
  -wwwroot="wwwroot": webserver document root folder

blog -verbose=true
```

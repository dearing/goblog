blog
====
**still in development**

blogging with go, markdown and redis

whatisit?
-----
Simply this little gem takes some templates and some articles written in markdown, parses tme and then pushes them to an active redis db where they are served for the web requests.  
A simple filewatcher then pushes modified documents as they are created or modified.


why use redis?
-----
The idea was to play around some with an engine and redis to store blog stats, changesets and timestamps as a complete solution for play.


user administration
----
I'm working on implementing Openid for administration needs and whatnot after I consider how this whole suite fits together.


and my file uploads?
----
I'm wresting with the idea of going with a RESTful service or using something cool and fancy like HTML5's file api or Websockets.  Practical concerns with proxy (read Nginx) and browser support with the technologies are the angles here.

an android app too?
----
Yea, I know overkill right?  This is all for fun and my personal practical use and I feel having a simplistic app on my phone will allow me to post more often.

extras
----
systemd units
nginx sample proxy configs
and android app

try it out
---
```
go get github.com/dearing/blog
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
```

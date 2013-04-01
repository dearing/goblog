messing with oauth2 (github)
====

of OpenID
----
Sometime yesterday I got around to pushing myself into the next phase of working on this blog, authentication.  
First I was thinking that [OpenID](http://en.wikipedia.org/wiki/OpenID) was the way to go since it is designed 
for authentication alone but as I dove deeper into figuring out implementing it I started to get the willies.
It would appear that OpenID is dying because I had trouble finding libraries to use and the links from the 
[official website](http://openid.net/) itself were dated. -See how [YADIS](http://yadis.org/wiki/Yadis_1.0_%28HTML%29) isn't the website you expect, granted YADIS not entirely an OpenID thing.-


on to OAuth2
----
So back to the drawing board I resurrected my interest in [OAuth](http://oauth.net/), specifically OAuth2.  Now OAuth2 
does support authentication but it also is about fetching data of that user.  So a user like *dearing* at Github has 
information the Github stores.  With an OAuth2 request (given a scope) we authorize an APP (like this blog) to fetch 
data on behalf of the client (well you can set data too, think posting tweets on Twitter for example) after an 
authentication -dance- as it is commonly called on the interwebs.

There was recent post on [Hacker News](http://news.ycombinator.com/) that lead me to a post 
[Thoughts on Go after writing 3 websites](http://blog.kowalczyk.info/article/uvw2/Thoughts-on-Go-after-writing-3-websites.html) 
some time back.  This post actually rekindled an interest I had in designing this whole thing but he talks about using Gary Burd's 
OAuth1 library to authenticate with twitter.  I wanted OAuth2 and prepared for some long nights like back when I 
decided to write a personal websocket library from the spec alone.  Luckily I found [goOauth2](https://code.google.com/p/goauth2/) 
for Go already and was quickly able to make use of it.

So the final bit of the puzzle was correctly identifing single valid admin for the blog itself.  I mean, can't let anyone with 
github account edit my posts on a whim now so I needed to include [securecookie](http://www.gorillatoolkit.org/pkg/securecookie) to store relevant information and a final 
login value to compare against information we pull with credentials from the [Github API v3](http://developer.github.com/v3/).

So here's the angle.  Anyone can pull [myself for example](https://api.github.com/users/dearing) without a token but the url for [authenticated user](http://developer.github.com/v3/users/#get-the-authenticated-user) only returns information for that [user](https://api.github.com/user).  So this suffices as authentication.  At least in my book.

Converting the existing code to handle other OAuth2 providers is trivial, one just needs to register an APP supply the credentials and compare values gleaned from whatever service for validation.
# An unsystematic list of the sources and tutorials I looked at while developing this project

## cookies
* https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies  
* https://stackoverflow.com/questions/8889679/how-to-create-a-cookie-and-add-to-http-response-from-inside-my-service-layer  
* https://dev.to/emmykolic/cookies-based-authentication-vs-session-based-authentication-1f6  
* https://stackoverflow.com/questions/17769011/how-does-cookie-based-authentication-work (stack exchange description of  process)
* https://www.golangprograms.com/how-do-you-set-cookies-in-an-http-request-with-an-http-client-in-go.html  
* https://www.calhoun.io/securing-cookies-in-go/  (secure cookies)
* https://www.alexedwards.net/blog/working-with-cookies-in-go  (complete guide)

## Some additional sources (see also capfront.03) which is currently a private archive:  
* https://www.jetbrains.com/guide/go/tutorials/rest_api_series/  
This frontend works in conjunction with the backend that provides an API through which this app retrieves data and to which it sends instructions.  

### jwt
* see https://www.sohamkamani.com/golang/jwt-authentication/  

### basic and token
* https://www.jetbrains.com/guide/go/tutorials/authentication-for-go-apps/auth/  

### sending POST request to remote server
* https://www.golangprograms.com/how-do-you-send-an-http-post-request-with-an-http-client-in-go.html  

### Servers, ssl, etc
* https://caddyserver.com/

### gin tutorials
* https://docs.studygolang.com/doc/tutorial/web-service-gin The docs
* https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/ gin binding (json - xml = structs) basic intro
* https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin  more detailed  
* https://dev.to/jacobsngoodwin/16-create-gin-middleware-to-extract-authorized-user-1jom more detailed stil
* https://www.squash.io/exploring-advanced-features-of-gin-in-golang/  
* https://chenyitian.gitbooks.io/gin-web-framework/content/docs/29.html custom function examples  
* https://developer.hashicorp.com/nomad/tutorials/templates/go-template-syntax  template syntax


# Notes on redirection

* https://gin-gonic.com/docs/examples/redirects/  
* https://stackoverflow.com/questions/68836237/how-to-get-full-server-url-from-any-endpoint-handler-in-gin  

# Handling text pages
I wanted a simple 'about' page to explain the project. This can be made more sophisticated as time goes on.

* https://go.dev/doc/articles/wiki/ very simple basic

# Initialising array structures
We use quite a few structures in which arrays are sub-nested. Space has to be reserved for these nested arrays explicitly, so they have to be initialised.  

* https://www.delftstack.com/howto/go/create-an-empty-slice-in-go/  

# Gin, middleware, and route handling
* https://www.squash.io/exploring-advanced-features-of-gin-in-golang/
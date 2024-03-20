# Frontend for capitalism simulation 'capfront.05'
* Interface with remote server via its API.  
* I wrote this in GO both to learn GO and to separate the UI side of the project from the logic.  
* So much of it is overcomplicated and not very clean. That's because it is a learning project.  
* The result seems a bit non-standard in that this 'frontend' is also a server.  
* I'm hoping it will eventually become a WASM project so that it's just a replacement for a normal js frontend   
  See https://golangbot.com/webassembly-using-go/  
  Meantime, it is server-resident.  
  
* The backend (which does the heavy lifting) is at http://github.com/axfreeman/capsim-12.4    
    
  Authentication is carried out in the backend. This frontend stores and uses the JWT token returned by the server


# Unsystematic list of required changes and fixes.  

* Display numbers nicely formatted.
* Move Action processing to sidebar  
* Get 'static' working  
* Browsers don't like the domain field in the cookie  
* Need a more efficient method for accessing links between objects (eg maps, or a local database)
* Limit number of simulations you can create 
  
## Admin
* Should be able to delete users 
* Should be able to delete user simulations 
* Should be able to view user simulations
* Menu in admin Dashboard needs to be styled properly
* complete admin dashboard  
  
## Authentication    
* Detecting logout attempt by user who is not logged in 
* check user didn't select non-existent individual item via the browser eg choose non-existent commodity 
* Switch gin to production mode
* Logout when browser closes (and clear cookie)  
* Clear cookie when logging out  
* Users should have a facility to resynchronise with the server. This is needed if, for example, the server has restarted (or reset) 
since their last login. Alternatively, the user should force resynchronise with each login (at present, we merely reconnect with the front end)  

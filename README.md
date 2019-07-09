# server 

Server is a package used to make a general design for all sort of server types.
This can be wrappers for a database, api or any kind of connection.

The use case is for Programs that needs to access a number of different apis, dbs or ssh connections.
Trying to make the design for each of those server wrappers as common as possible for an easier overview and maintaining state.

Server is actually just an interfave with a set of methods to have.
The rest of this package is for making it easier to follow the server standard, such as the cmd tool to generate a template for 
the kind of wrapper wanted, etc db.

There is a ServerContainer which is a helper struct that can store these generated servers
and making sure Connections are maintained or reconnected etc.


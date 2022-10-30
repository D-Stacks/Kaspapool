high priority 

- Implement command client
- Implement and handle keep alive
- Custom block and hashing - don't rely on external api as much
	+ same for kaspa-addresses, but less important. 
- Downsize Share struct to most important values. 

mid priority

- Handle config 
- Handle flag parsing

low priority

- don't rely on kaspad grpc client, build own grpc connection.
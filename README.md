# Service

A utility to standardize and simplify managing go routines in a go program. It's intent is to handle channels and closing all serviceGroups specified if one of the contained services fails.

Services could be long running or short running routines.

## Interface

Responsibilities of a Service

1. Log all errors
2. Run until the completion of the service
3. Only return fatal errors to the ServiceGroup
4. Return `nil` on successful execution or graceful exit

`func Start() error` : The entry point of the service, should not return until the service is done.

`func Stop()`: The graceful stop function that should take a thread safe operation on the implementation of `Service` to cause the `Start` function to exit.

## ServiceGroup

A list of services to start, watch for errors and properly exit all services in the same serviceGroup. There should be some dependence between services in a serviceGroup since a fatal error in one of the services causes the entire ServiceGroup to shut down.

`New` : creates a new ServiceGroup
`Add` : adds a service to the ServiceGroup
`Wait`: ensures the main thread will block until all services in the ServiceGroup are done running
`Kill`: force close everything in a ServiceGroup
`Start`: starts all the services in the ServiceGroup


## Other Notes

Catching SIGNALS is left to the user of `Service`, and Call `Kill` to force all services in the group to stop.

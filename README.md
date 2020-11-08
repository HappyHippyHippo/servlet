# Servlet

This library intent to provide a flexible and opinionated structure to build
applications in a way that removes the problems of coding the same
functionalities over and over regarding system components as configuration
or logging.

Is the opinion of the library that the application scaffolding should revolve
around the idea of a service container, seen in other types of frameworks.
This container is responsible for storing factories of the application
services that should be called for service instantiation, per need. This
facilitates not only the composition of the application, but also directs the
application architecture for a more dependency injection way of structuring
services. By following this decoupled way of writing the application. will
enable for more scalable and maintainable testing and production code.

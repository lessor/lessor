# Lessor

Lessor is a set of tools for deploying single-tenant applications in a secure, distributed fashion. This allows you to proxy to and independently scale each tenant, with network and data isolation by default.

Lessor is a permissively licensed open source project that was created by the SRE team at [Kolide](https://kolide.com) to deploy, manage, observe, and secure the Kolide Cloud product.

## Motivation

Companies that create products for other companies or teams often have to reason about how to deal with the tenancy of each team. There are generally two paths:

- Deploy one monolithic application that handles multi-tenant data isolation via application logic
- Deploy and proxy to many instances of smaller, more isolated single-tenant applications

When faced with these two options, most companies choose to build the multi-tenant monolith. While the second path results in simpler, more secure software, many single-tenant applications are much more difficult to operate and observe. Large multi-tenant monoliths, however, have a habit of becoming difficult to operate and observe as well though.

Lessor aims to make it easier to choose to deploy and proxy to many instances of a single-tenant application by providing tools, services, and libraries that are purpose-built for this kind of deployment strategy.

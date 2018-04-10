# Frequently Asked Questions

<p align="center">
  <img src="./images/gophers/learn.png" width="500">
</p>

### Labels in the metadata of a Tenant

Most tenant resources start off like this:

```
apiVersion: lessor.io/v1
kind: Tenant
metadata:
  name: acme-labs
  labels:
    name: acme-labs
```

The `name: acme-labs` may seem redundant, but the reason that it's there is because the Kubernetes API only allows for certain operations (like search) based on labels, not fields. So to update tenants by name programatically, it is useful (if not necessary) to include the name as label.

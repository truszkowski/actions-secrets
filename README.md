# GitHub Actions Secrets Batch Management Tool

If you want to add many secrets at once, you can use this tool.
Also, if you want, delete them. 

## Token

You need to define `TOKEN` env-var with your GitHub token.

## List all secrets

```
./actions-secrets -owner <your-namespace> -repo <your-repository> -list-all
```

## List secrets defined in file

```
./actions-secrets -owner <your-namespace> -repo <your-repository> -list <env-file>
```

## Apply secrets

```
./actions-secrets -owner <your-namespace> -repo <your-repository> -apply <env-file>
```

## Delete secrets

```
./actions-secrets -owner <your-namespace> -repo <your-repository> -delete <env-file>
```

## Other options

- `-override` with `-apply` overrides env-var (otherwise we only check if it exists)
- `-verbose` a little bit more informations



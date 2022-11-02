# About

[oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy) as an GitHub authenticating reverse proxy example.

# Usage

Add `127.0.0.1 example.test` to your `hosts` file.

Register a [new GitHub OAuth Application](https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/oauth_provider#github-auth-provider) with:

| Setting                      | Value                                      |
|------------------------------|--------------------------------------------|
| `Homepage URL`               | `http://example.test:4180`                 |
| `Authorization callback URL` | `http://example.test:4180/oauth2/callback` |

Generate a new client secret, and export the created application OAuth credentials as environemnt variables:

```bash
export OAUTH2_PROXY_PROVIDER='github'
export OAUTH2_PROXY_SCOPE='user:email'
#export OAUTH2_PROXY_GITHUB_USER='rgl'
export OAUTH2_PROXY_CLIENT_ID='YOUR_OAUTH2_PROXY_GITHUB_APP_CLIENT_ID'
export OAUTH2_PROXY_CLIENT_SECRET='YOUR_OAUTH2_PROXY_GITHUB_APP_CLIENT_SECRET'
```

Download `oauth2-proxy`:

```bash
wget https://github.com/oauth2-proxy/oauth2-proxy/releases/download/v7.4.0/oauth2-proxy-v7.4.0.windows-amd64.tar.gz
tar xf oauth2-proxy-v7.4.0.windows-amd64.tar.gz --strip-components 1
```

Execute `oauth2-proxy`:

```bash
export OAUTH2_PROXY_COOKIE_SECRET="$(openssl rand -hex 16)"
./oauth2-proxy \
    --email-domain=* \
    --http-address=:4180 \
    --redirect-url=http://example.test:4180/oauth2/callback \
    --cookie-secure=false \
    --upstream="file:///$(cygpath --windows "$PWD" | tr \\\\ /)/#/"
```

Access the root endpoint:

1. http://example.test:4180

Access some of the [endpoints](https://oauth2-proxy.github.io/oauth2-proxy/docs/features/endpoints):

1. http://example.test:4180/oauth2/userinfo
1. http://example.test:4180/oauth2/sign_out

# Alternatives

* [caddy-security](https://github.com/greenpau/caddy-security)
* [ory](https://github.com/ory)

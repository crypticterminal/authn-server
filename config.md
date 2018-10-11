---
title: Server Configuration
---

# Server Configuration

* Core Settings: [`AUTHN_URL`](#authn_url) • [`APP_DOMAINS`](#app_domains) • [`HTTP_AUTH_USERNAME`](#http_auth_username) • [`HTTP_AUTH_PASSWORD`](#http_auth_password) • [`SECRET_KEY_BASE`](#secret_key_base)
* Databases: [`DATABASE_URL`](#database_url) • [`REDIS_URL`](#redis_url)
* Sessions:
[`ACCESS_TOKEN_TTL`](#access_token_ttl) • [`REFRESH_TOKEN_TTL`](#refresh_token_ttl) • [`SESSION_KEY_SALT`](#session_key_salt) • [`DB_ENCRYPTION_KEY_SALT`](#db_encryption_key_salt) • [`RSA_PRIVATE_KEY`](#rsa_private_key)
* OAuth Clients: [`FACEBOOK_OAUTH_CREDENTIALS`](#facebook_oauth_credentials) • [`GITHUB_OAUTH_CREDENTIALS`](#github_oauth_credentials) • [`GOOGLE_OAUTH_CREDENTIALS`](#google_oauth_credentials)
* Username Policy: [`USERNAME_IS_EMAIL`](#username_is_email) • [`EMAIL_USERNAME_DOMAINS`](#email_username_domains)
* Password Policy: [`PASSWORD_POLICY_SCORE`](#password_policy_score) • [`BCRYPT_COST`](#bcrypt_cost)
* Password Resets: [`APP_PASSWORD_RESET_URL`](#app_password_reset_url) • [`PASSWORD_RESET_TOKEN_TTL`](#password_reset_token_ttl) • [`APP_PASSWORD_CHANGED_URL`](#app_password_changed_url)
* Stats: [`TIME_ZONE`](#time_zone) • [`DAILY_ACTIVES_RETENTION`](#daily_actives_retention) • [`WEEKLY_ACTIVES_RETENTION`](#weekly_actives_retention)
* Operations: [`PORT`](#port) • [`PUBLIC_PORT`](#public_port) • [`PROXIED`](#proxied) • [`SENTRY_DSN`](#sentry_dsn) • [`AIRBRAKE_CREDENTIALS`](#airbrake_credentials)

## Core Settings

### `AUTHN_URL`

|           |    |
| --------- | --- |
| Required? | Yes |
| Value | URL |

This specifies the base URL of the AuthN service. It will be embedded in all issued JWTs as the `iss`. Clients will depend on this information to find and fetch the service's public key when verifying JWTs.

### `APP_DOMAINS`

|           |    |
| --------- | --- |
| Required? | Yes |
| Value | comma-delimited list of domains (scheme & host, no path) |

Any domain listed in this variable will be trusted for three things:

1. Requests sent from these domains (as determined by the Origin header) will satisfy CSRF requirements.
2. Access tokens generated by requests sent from these domains (as determined by the Origin header) will specify the domain as their intended `aud` (audience).
3. Any endpoints that accept redirects will only allow the redirect if it uses one of these domains.

### `HTTP_AUTH_USERNAME`

|           |    |
| --------- | --- |
| Required? | Yes |
| Value | string |

Any access to private AuthN endpoints must use HTTP Basic Auth, with this username.

### `HTTP_AUTH_PASSWORD`

|           |    |
| --------- | --- |
| Required? | Yes |
| Value | string |

Any access to private AuthN endpoints must use HTTP Basic Auth, with this password.

### `SECRET_KEY_BASE`

|           |    |
| --------- | --- |
| Required? | Yes |
| Value | string |

Any HMAC keys used by AuthN will be derived from this base value. Currently this only includes the key used to securely sign sessions maintained with the AuthN service.

This value is commonly a 64-byte string, and can be generated with [`SecureRandom.hex(64)`](http://ruby-doc.org/stdlib-2.3.3/libdoc/securerandom/rdoc/Random/Formatter.html#method-i-hex) or `bin/rake secret`. Some deployment systems (e.g. Heroku) can provision it automatically.

## Databases

### `DATABASE_URL`

|           |    |
| --------- | --- |
| Required? | Yes |
| Value | string |

The database URL specifies the driver, host, port, database name, and connection credentials.

Formats:

* `sqlite3://local/db/authn` (note: SQLite3 ignores the host name and connects by path)
* `mysql://username:password@host:port/database_name`
* `postgres://username:password@host:port/database_name`

### `REDIS_URL`

|           |    |
| --------- | --- |
| Required? | No |
| Value | string |

Redis is the preferred database for session refresh tokens, encrypted blobs, and active user stats.
Currently the SQLite3 database is able to manage those functions in a limited deployment, but Redis
is required when SQLite3 is not configured.

Format:

* `redis://username:password@host:port/database_number`

## Sessions

### `ACCESS_TOKEN_TTL`

|           |    |
| --------- | --- |
| Required? | No |
| Value | seconds |
| Default | `3600` (1 hour) |

This setting controls how long the access tokens (aka application sessions) will live. This is an important precaution because it allows the AuthN server to revoke sessions (e.g. on logout) with confidence that any related access tokens will expire soon and have limited damage potential.

Worried about short sessions? Applications can and should implement a periodic refresh process to keep the effective session alive much longer than the expiry listed here. The [keratin/authn-js](https://github.com/keratin/authn-js) client library implements a half-life maintenance strategy when you configure it to manage sessions. This strategy will attempt to refresh the session when it has half-expired, or earlier if there's reason to severely distrust the client's clock. If a user closes their client and doesn't return before the access token expires, the refresh logic will restore their session on the first page load.

### `REFRESH_TOKEN_TTL`

|           |    |
| --------- | --- |
| Required? | No |
| Value | seconds |
| Default | `2592000` (30 days) |

This setting controls how frequently a refresh token must be used to keep a session alive. Changing this setting will not apply retroactively to previous tokens.

### `SESSION_KEY_SALT`

|           |    |
| --------- | --- |
| Required? | No |
| Value | string |
| Default | `session-key-salt` |

This salt is added to [`SECRET_KEY_BASE`](#secret_key_base) and used to derive the session key. Customizing this value can provide extra defense against brute-force attacks on AuthN's HMAC signatures, but is not required because the work factor involved in a brute-force attack already involves 20k rounds of SHA-256 per guess.

### `DB_ENCRYPTION_KEY_SALT`

|           |    |
| --------- | --- |
| Required? | No |
| Value | string |
| Default | `db-encryption-key-salt` |

This salt is added to [`SECRET_KEY_BASE`](#secret_key_base) and used to derive the encryption key for objects stored in a database. Customizing this value can provide extra defense against brute-force attacks on stolen or leaked data, but is not required because the work factor involved in a brute-force attack already involves 20k rounds of SHA-256 per guess.

### `RSA_PRIVATE_KEY`

|           |    |
| --------- | --- |
| Required? | No |
| Value | PEM |
| Default | none |

The private key must be in PEM format, with no passphrase. If you've run `ssh-keygen -N '' -f keratin-authn-rsa`, then you can get the PEM private key by copying the entire output of `cat keratin-authn-rsa`.

Some systems (e.g. Heroku) make it easy to add multi-line environment variables. If your system does not, you may collapse the public key into a single line by replacing all line breaks with `\n` characters.

Note that specifying a `RSA_PRIVATE_KEY` will prevent AuthN from automatically rotating keys. If you wish to implement your own key rotation, remember to restart the process to pick up changes.

## OAuth Clients

When configuring OAuth you will need to know your AuthN server's return URL. You may determine this by joining the AuthN server's base URL with the path `/oauth/:providerName/return`. For example, for Google you might enter:

* `https://authn.example.com/oauth/google/return`

or

* `https://www.example.com/authn/oauth/google/return`

### `FACEBOOK_OAUTH_CREDENTIALS`

|           |    |
| --------- | --- |
| Required? | No |
| Value | AppID:AppSecret |
| Default | nil |

Create a Facebook app at https://developers.facebook.com and enable the Facebook Login product. In the Quickstart, enter [AuthN's OAuth Return](api.md#oauth-return) as the Site URL. Then switch over to Settings and find the App ID and Secret. Join those together with a `:` and provide them to AuthN as a single variable.

### `GITHUB_OAUTH_CREDENTIALS`

|           |    |
| --------- | --- |
| Required? | No |
| Value | ClientID:ClientSecret |
| Default | nil |

Sign up for GitHub OAuth 2.0 credentials with the instructions here: https://developer.github.com/apps/building-oauth-apps. Your client's ID and secret must be joined together with a `:` and provided to AuthN as a single variable.

### `GOOGLE_OAUTH_CREDENTIALS`

|           |    |
| --------- | --- |
| Required? | No |
| Value | ClientID:ClientSecret |
| Default | nil |

Sign up for Google OAuth 2.0 credentials with the instructions here: https://developers.google.com/identity/protocols/OpenIDConnect. Your client's ID and secret must be joined together with a `:` and provided to AuthN as a single variable.

## Username Policy

### `USERNAME_IS_EMAIL`

|           |    |
| --------- | --- |
| Required? | No |
| Value | boolean (`/^t|true|yes$/i`) |
| Default | `false` |

If you ask users to sign up with an email address, enable this so that AuthN can validate properly.

### `EMAIL_USERNAME_DOMAINS`

|           |    |
| --------- | --- |
| Required? | No |
| Value | comma-delimited list of domains |
| Default | nil |

If you need to restrict account creation to specific email domains, declare the domains here. Note that your application is still responsible for verifying email ownership.

## Password Policy

### `PASSWORD_POLICY_SCORE`

|           |    |
| --------- | --- |
| Required? | No |
| Value | 0 - 5 |
| Default | 2 |

* 0 - too guessable
* 1 - very guessable
* 2 - somewhat guessable
* 3 - safely unguessable
* 4 - very unguessable

Password complexity is calculated by estimating how many guesses it would take a smart attacker armed with a dictionary, simple transformations like L337, and spatial walks across the QWERTY keyboard. The specific algorithm used is [zxcvbn](https://blogs.dropbox.com/tech/2012/04/zxcvbn-realistic-password-strength-estimation/), which has a JavaScript implementation if you'd like to provide real-time user feedback on password fields.

### `BCRYPT_COST`

|           |    |
| --------- | --- |
| Required? | No |
| Value | 10+ |
| Default | `11` |

BCrypt costs describe how many times a password should be hashed. Costs are exponential, and may be increased later without waiting for a user to return and log in.

The ideal cost is the slowest one that can be performed without _feeling_ slow and without creating CPU bottlenecks or easy DDOS attacks on your AuthN server. There's no reason to go below 10, and 12 starts to become noticeable, so 11 is the default.

Please run your own benchmarks, but consider that on my Macbook Pro:

| Cost | Iterations | Time |
| ---- | ---------- | ---- |
| 10   | 1024       | ~0.067s |
| 11   | 2048       | ~0.136s |
| 12   | 4096       | ~0.276s |

## Password Resets

### `APP_PASSWORD_RESET_URL`

|           |    |
| --------- | --- |
| Required? | No |
| Value | URL |
| Default | nil |

Must be provided to enable password resets. This URL must respond to `POST`, should expect to receive `account_id` and `token` params, and is expected to deliver the `token` to the specified `account_id`.

### `PASSWORD_RESET_TOKEN_TTL`

|           |    |
| --------- | --- |
| Required? | No |
| Value | seconds |
| Default | 1800 (30.minutes) |

Specifies the amount of time a user has to complete a password reset process. After this period of time, the reset token will no longer be accepted. (Note that a reset token will also be invalidated if the password changes before this TTL.)

### `APP_PASSWORD_CHANGED_URL`

|           |    |
| --------- | --- |
| Required? | No |
| Value | URL |
| Default | nil |

Must be provided to enable notifications of password changes. This URL must respond to `POST`, should expect to receive an `account_id` param, and is expected to deliver an email confirmation.

## Stats

### `TIME_ZONE`

|           |    |
| --------- | --- |
| Required? | No |
| Value | Time Zone descriptor |
| Default | `UTC` |

Specifies the time zone for tracking and reporting stats. Note that changing this time zone is applied on write, not on read. This means that if you change the time zone later your historical stats (e.g. daily active accounts) will not update.

### `DAILY_ACTIVES_RETENTION`

|           |    |
| --------- | --- |
| Required? | No |
| Value | days |
| Default | `365` (~1 year) |

Stats on daily actives will be set to expire after this many days. No mechanism is provided for changing this TTL retroactively.

### `WEEKLY_ACTIVES_RETENTION`

|           |    |
| --------- | --- |
| Required? | No |
| Value | years |
| Default | `104` (~2 years) |

Stats on weekly actives will be set to expire after this many weeks. No mechanism is provided for changing this TTL retroactively.

## Operations

### `PORT`

|           |    |
| --------- | --- |
| Required? | No |
| Value | integer |
| Default | from AUTHN_URL |

The PORT specifies where the AuthN server should bind. This may be different from the AUTHN_URL in scenarios with port mapping, as with load balancers and Docker containers.

### `PUBLIC_PORT`

|           |    |
| --------- | --- |
| Required? | No |
| Value | integer |
| Default | nil |

Specifying PUBLIC_PORT instructs AuthN to bind on a second port with only public routes. This supports network configurations with separate public and private routing. The public load balancer can route to the public port without needing to create and maintain path- & method-based lists of allowed endpoints.

### `PROXIED`

|           |    |
| --------- | --- |
| Required? | No |
| Value | boolean (`/^t|true|yes$/i`) |
| Default | `false` |

Specifying PROXIED allows AuthN to safely read common proxy headers like X-FORWARDED-FOR to determine the true client's IP address. This is currently useful for logging.

### `SENTRY_DSN`

|           |     |
| --------- | --- |
| Required? | No |
| Value | string |
| Default | nil |

Configures AuthN to report panics and unhandled errors to a Sentry backend.

### `AIRBRAKE_CREDENTIALS`

|           |     |
| --------- | --- |
| Required? | No |
| Value | string |
| Default | nil |

Configures AuthN to report panics and unhandled errors to an Airbrake backend. The format is `projectID:projectKey`.
# encrypt

The `encrypt` package provides AES-256-GCM encryption for snapshot data stored
on disk, protecting sensitive configuration content at rest.

## Components

### `Encryptor`

Derives a 32-byte key from a passphrase using SHA-256 and exposes `Encrypt` /
`Decrypt` methods.

```go
e := encrypt.New("my-passphrase")
ct, err := e.Encrypt([]byte("password=secret"))
pt, err := e.Decrypt(ct)
```

- Each `Encrypt` call generates a random nonce, so identical plaintexts produce
  different ciphertexts.
- `Decrypt` returns `ErrInvalidCiphertext` when the data has been tampered with
  or the wrong key is used.

### `Store`

A thin file-system layer that transparently encrypts data on write and decrypts
on read.

```go
s, err := encrypt.NewStore("/var/lib/confsnap/enc", "my-passphrase")
_ = s.Write("web01:/etc/nginx/nginx.conf", plaintext)
data, err := s.Read("web01:/etc/nginx/nginx.conf")
```

Files are stored with mode `0600`; the store directory is created with `0700`.
Key strings containing `/` or `\` are sanitised to `_` before use as file names.

## Security notes

- The passphrase should be sourced from an environment variable or a secrets
  manager rather than hard-coded.
- AES-256-GCM provides both confidentiality and integrity; any modification to
  the ciphertext will cause `Decrypt` to return an error.

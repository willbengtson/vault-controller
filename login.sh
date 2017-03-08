curl $VAULT_ADDR/v1/auth/userpass/login/testuser \
    -d '{ "password": "testpassword" }' > .vault-response

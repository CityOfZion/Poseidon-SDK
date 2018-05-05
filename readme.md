#ReadMe


In each coin's implementation, there is a coin struct which will be used to liberally add things I may need in the future.
All methods have been attached to this struct, so that the API may call the methods like so: Neo.Coin.GeneratePrivateKey()


Need to implement Signing and Verifying feature, a checkWIF and decode function

Todos:

- Test Signing function.
- Implement verify in signing.go
- implement checkWIF function
- Implement DecodeWif function
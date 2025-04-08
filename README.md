# Uniswap price checker with go 

Small Go-based HTTP API that return prcies of Uniswap V3 pool

## Running the API

1. **Clone the repository**:
   ```bash
   git clone https://github.com/gokuseii/go-uniswap-price-fetcher
   cd go-uniswap-price-fetcher
   ```

2. Install dependencies:

   ```
   go mod tidy
   ```

3. Run server
   ```
   go run main.go
   ```

Test:
   ```
   curl "http://localhost:8080/price?pool=0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8"
   ```

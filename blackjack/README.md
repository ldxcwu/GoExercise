# **BlackJack Game**


# **Gains**
### **1. Import local module**
   ```go
   go mod init blackjack
   go mod edit -replace deck=../deck
   go mod tidy
   ```
   参考：[官网文档--Call Your Code From Another Module](https://go.dev/doc/tutorial/call-module-code)

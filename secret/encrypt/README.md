# encrypt pkg
This package is build to provide aes encrypt and decrypt functions.   
1. Use md5 to hash a provided key  
   
   why use md5 first is because that the key's length we provided is uncertain.  
   The encrypt function may not work well if the key is too long.   
   But md5 can generate a hex with specific length. 
    ```go
    hasher := md5.New()
    //hash.Hash implemented io.Writer interface
    //or fmt.Fprintf(hash, key)
    //or io.WriteString(hash, key)
    hasher.Write([]byte(key))
    newKey := hasher.Sum(nil)
    ```
2. Use aes to build an implentation of cipher.Block
3. Build a cfb with the Block and other args.
4. Encrypt and Decrypt.
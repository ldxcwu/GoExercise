# 问题描述
``` 读取指定目录，并修改内部文件的文件名 ```
# 解题步骤
1. 读取目录文件  
1.a ioutil.ReadDir(path) 该方法不能递归查找子目录，需要手动编写  
1.b filepath.Walk(path, func) 推荐，该方法自动递归查找，并对途中每个文件调用func方法
2. func的第一个参数是文件的路径，第二个参数fs.Fileinfo对象，其 .name 就只是文件名，不包含路径
3. 自定义文件名修改规则
4. 拼接文件名和路径，也就是上述 path 和 info.name  
4.a 可以使用fmt.Sprintf()进行格式化拼接 
4.b 可以使用filepath.join(s,s) 拼接，串中间自动加斜杠
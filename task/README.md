# Task
Task is a simple CLI app that you can record tasks need to be completed.  
Task CLI use cobra to build cli and use bolt to store tasks info.
## Usage
1.  download Task CLI
2.  ```export PATH=$PATH:$GOPATH/bin```
3.  ```go install .```
4.  ```task``` to get help info
5.  ```task list``` to show tasks you have now
6.  ```task add taskName``` to add a task in task list
7.  ```task do id``` to mark task as completed
##  Gains
1. cobra is a powerful tool to help us build CLI tools quickly.
2. cobra needs a directory named cmd in which to write all command.go file
3. bolt is a lightly key-value database.
4. bolt provides ```db.Update(func())``` to start a read-write transaction.
5. bolt provides ```db.View(func())``` to start a read-only transaction.
6. go install command will put executable file after build in $GOPATH/bin
7. export $GOPATH/bin to make our tool run anywhere.
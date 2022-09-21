// Package killdren do these two things:
//
//  - Do best to ensure children get killed when panic
//  - Stop/kill all "grandchild process" when stop/kill child process
//
// "grandchild process": "child process" started by child process started in go
package killdren

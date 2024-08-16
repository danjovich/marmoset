package arm

import "fmt"

func MakeSyscall(nr int, args string) string {
	return fmt.Sprintf(`	%s 
	mov r7, #%d 
	svc #0`, args, nr)
}

package main

import (
	"log"
	"sync"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Sample Messages (who|what) or (who|what|where)
	body := [10]string{
		"c68b1eb3-75e5-446e-b1fe-7f26fc5588e4|public.post.create",
		"9da1c4e6-b674-4536-9b04-7a7e781c26cf|public.post.create",
		"c248cb3f-0d7b-487c-8ed9-f90bdbe6d19a|public.post.create",
		"f9b91a8e-2280-4d1e-9987-7a1918a3eba7|public.post.create",
		"821e60f0-aeda-4533-bfd1-e853deb8edc3|public.post.create",
		"ef70a148-d046-4a88-815a-d8142f15839b|public.post.create",
		"c248cb3f-0d7b-487c-8ed9-f90bdbe6d19a|public.post.delete",
		"f9b91a8e-2280-4d1e-9987-7a1918a3eba7|public.post.delete",
		"821e60f0-aeda-4533-bfd1-e853deb8edc3|public.post.delete",
		"ef70a148-d046-4a88-815a-d8142f15839b|public.post.delete",
	}

	perms := makePermissionsManager("172.17.0.2", 5672, "user", "pass")
	defer perms.Close()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			log.Printf("[%d]", i)
			perms.check(body[i*13%10])
			// perms.check(body[i])
			wg.Done()
		}(i)
	}

	wg.Wait()

	// perms.check(body[1])

	// for i := 0; i < 10; i++ {
	// 	perms.check(body[i])
	// }

	// var wg sync.WaitGroup

	// for i := 0; i < 5; i++ {
	// 	wg.Add(1)

	// 	go func() {
	// 		defer wg.Done()
	// 		sendMessages(50000)
	// 	}()
	// }

	// wg.Wait()

	// sendMessages(50000)
}

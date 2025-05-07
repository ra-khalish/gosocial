package db

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/ra-khalish/gosocial/internal/store"
)

var usernames = []string{
	"Alice", "Bob", "Charlie", "David", "Emma", "Frank", "Grace", "Hank", "Ivy", "Jack",
	"Karen", "Leo", "Mia", "Nathan", "Olivia", "Peter", "Quinn", "Rachel", "Sam", "Tina",
	"Uma", "Victor", "Wendy", "Xander", "Yvonne", "Zane", "Aaron", "Bella", "Caleb", "Diana",
	"Ethan", "Fiona", "George", "Holly", "Isaac", "Jasmine", "Kevin", "Luna", "Mason", "Nina",
	"Oscar", "Penny", "Quincy", "Riley", "Sophia", "Travis", "Ursula", "Vincent", "Willow", "Xavier",
}

var titles = []string{
	"Mastering Go Routines",
	"Building APIs with Go",
	"Understanding Pointers in Go",
	"Effective Error Handling in Go",
	"Using PostgreSQL with Go",
	"Introduction to Generics in Go",
	"Writing Clean Code in Go",
	"Go Concurrency Patterns",
	"Structs vs Interfaces in Go",
	"Unit Testing in Go",
	"Optimizing Go Performance",
	"Dependency Injection in Go",
	"Working with JSON in Go",
	"Building a CLI with Go",
	"Understanding Context in Go",
	"Using WebSockets in Go",
	"Logging Best Practices in Go",
	"Microservices with Go and gRPC",
	"Go Memory Management",
	"Secure Coding in Go",
}

var contents = []string{
	"Go routines allow for lightweight concurrency, making it easy to handle multiple tasks in parallel.",
	"Building RESTful APIs in Go is efficient with frameworks like Gin and Echo, offering high performance.",
	"Pointers in Go enable direct memory access, which can optimize performance when handling large data structures.",
	"Effective error handling in Go involves using custom error types and wrapping errors for better debugging.",
	"Connecting Go to PostgreSQL is simple with the pgx package, which offers a performant and feature-rich database driver.",
	"Generics in Go reduce code duplication by allowing functions and data structures to operate on multiple types.",
	"Writing clean Go code involves following idioms like short variable names, meaningful comments, and consistent formatting.",
	"Go’s concurrency model is based on goroutines and channels, making it easy to write highly scalable applications.",
	"Structs and interfaces are key components of Go’s type system, allowing flexible and efficient data modeling.",
	"Unit testing in Go is streamlined with the testing package, enabling easy test automation and coverage tracking.",
	"Optimizing Go performance involves profiling CPU/memory usage and reducing unnecessary allocations.",
	"Dependency injection in Go enhances testability and maintainability by reducing direct dependencies.",
	"JSON encoding/decoding in Go is handled efficiently with the encoding/json package for API communication.",
	"Building CLI tools in Go is straightforward with the Cobra package, which provides powerful command parsing.",
	"The context package in Go helps manage request lifecycles, preventing memory leaks and handling timeouts.",
	"WebSockets in Go enable real-time communication between clients and servers using the gorilla/websocket package.",
	"Logging in Go can be improved using structured logging libraries like Zerolog and Logrus for better observability.",
	"Microservices in Go benefit from gRPC’s efficient binary serialization and automatic code generation.",
	"Go’s garbage collector efficiently manages memory, reducing the risk of memory leaks in long-running applications.",
	"Secure coding in Go requires input validation, proper authentication, and safe use of cryptographic functions.",
}

var tags = []string{
	"golang", "programming", "web-development", "api", "concurrency",
	"database", "postgresql", "microservices", "performance", "testing",
	"docker", "kubernetes", "security", "json", "cli",
	"websockets", "logging", "generics", "grpc", "error-handling",
}

var comments = []string{
	"Great article! Really helped me understand goroutines.",
	"Could you provide an example with error handling?",
	"This was exactly what I was looking for, thanks!",
	"I think you could expand on the database section.",
	"Awesome post! Keep up the good work.",
	"I tried this approach, but I ran into an issue with memory usage.",
	"Does this work with the latest Go version?",
	"Thanks for sharing! Very well explained.",
	"Could you compare this with another framework like Fiber?",
	"This cleared up so much confusion for me. Thanks!",
	"I appreciate the practical examples. They made it easy to follow.",
	"How would you implement this in a microservices architecture?",
	"Nice breakdown! Would love to see a performance comparison.",
	"Any recommendations for handling large-scale deployments?",
	"Can you write a follow-up on best practices?",
	"This was super helpful. I’ll try it in my next project.",
	"Your explanation of interfaces was spot on!",
	"I never knew about this technique. Thanks for sharing!",
	"Great post! Any thoughts on handling authentication?",
	"Would love to see a video tutorial on this.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	comments := generateComment(500, users, posts)
	for _, comment := range comments {
		if err := store.Comment.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.IntN(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.IntN(len(titles))],
			Content: contents[rand.IntN(len(contents))],
			Tags: []string{
				tags[rand.IntN(len(tags))],
				tags[rand.IntN(len(tags))],
			},
		}
	}

	return posts
}

func generateComment(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
			Content: comments[rand.IntN(len(comments))],
		}
	}
	return cms
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/codepnw/gopher-social/internal/database"
	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/codepnw/gopher-social/internal/repository"
	"github.com/codepnw/gopher-social/internal/utils/env"
	"github.com/joho/godotenv"
)

var usernames = []string{
	"alice", "bob", "dave", "eve", "frank", "grace", "heidi", "ivan", "judy", "mallory",
	"nancy", "oliver", "pat", "quincy", "rachel", "steve", "trent", "ursula", "victor", "wendy",
}

var titles = []string{
	"Mastering Go: Tips & Tricks",
	"Clean Architecture in Go",
	"Building APIs with Gin",
	"PostgreSQL Best Practices",
	"Understanding JWT in Go",
	"Microservices: The Essentials",
	"Docker for Go Developers",
	"Writing Efficient Go Code",
	"Debugging Go Applications",
	"Scaling Your Go Services",
	"Concurrency in Go: A Guide",
	"gRPC vs REST in Go",
	"Using WebSockets in Go",
	"Logging Best Practices in Go",
	"Testing Go Applications",
	"Optimizing SQL Queries in Go",
	"Implementing OAuth2 in Go",
	"Go Generics: How & Why",
	"Using Redis with Go",
	"Event-Driven Architecture in Go",
}

var contents = []string{
	"Go is a powerful language for building scalable applications. This guide covers essential tips and tricks.",
	"Clean Architecture helps structure your Go applications efficiently. Learn how to implement it step by step.",
	"Building REST APIs with Gin can be fast and easy. This tutorial walks through the key concepts.",
	"Optimizing PostgreSQL queries can improve performance. Discover the best practices for Go developers.",
	"JWT authentication is essential for secure APIs. This guide explains how to implement it in Go.",
	"Microservices architecture allows scalable development. Learn how to build and deploy microservices in Go.",
	"Docker simplifies Go application deployment. This post covers how to containerize your Go projects.",
	"Writing efficient Go code is key to performance. Learn optimization techniques and best practices.",
	"Debugging Go applications can be challenging. Discover useful tools and techniques to find and fix issues.",
	"Scaling Go services requires proper architecture and tools. Learn the best strategies for handling high traffic.",
	"Go concurrency allows efficient multitasking. This guide explains goroutines, channels, and best practices.",
	"gRPC provides efficient communication in Go. Learn how it compares to REST and when to use it.",
	"WebSockets enable real-time communication in Go applications. This post covers implementation and use cases.",
	"Logging is crucial for debugging and monitoring. Explore logging best practices for Go applications.",
	"Unit and integration testing ensure code reliability. Learn how to write effective tests in Go.",
	"Slow SQL queries impact performance. This guide covers how to optimize queries in Go applications.",
	"OAuth2 provides secure authentication. Learn how to integrate OAuth2 in Go applications.",
	"Go Generics simplify code reuse. Understand how and when to use them in real-world projects.",
	"Redis is a powerful caching solution. Learn how to use Redis with Go for improved performance.",
	"Event-driven architecture helps scale applications. Explore how to implement it in Go.",
}

var tags = []string{
	"go", "golang", "web-development", "api", "postgresql",
	"microservices", "docker", "jwt", "clean-architecture", "grpc",
}

var comments = []string{
	"Great article! Really helped me understand Go better.",
	"Can you provide more examples for beginners?",
	"Nice explanation! Would love to see a deep dive into concurrency.",
	"This saved me hours of debugging. Thanks!",
	"Could you compare this approach with another framework?",
	"Well written! Do you have any GitHub repos for reference?",
	"How would you implement this in a real-world project?",
	"This was exactly what I was looking for!",
	"One of the best tutorials I’ve read on this topic.",
	"Can you explain error handling in more detail?",
	"Amazing post! Helped me optimize my API performance.",
	"Would you recommend using this approach for a large-scale project?",
	"Do you have any tips for debugging issues in production?",
	"I followed your guide and it worked perfectly!",
	"Could you add a section on best security practices?",
	"This is very useful! Looking forward to your next post.",
	"Your examples are clear and easy to follow. Thanks!",
	"How does this compare to other programming languages?",
	"I ran into an issue while implementing this. Can you help?",
	"This was super informative. Keep up the great work!",
}

func main() {
	godotenv.Load("dev.env")

	addr := env.GetString("DB_ADDR", "")
	conn, err := database.NewDatabase(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	seed(conn)
}

func seed(db *sql.DB) {
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	ctx := context.Background()
	tx, _ := db.BeginTx(ctx, nil)

	users := generateUsers(20)
	for _, user := range users {
		if err := userRepo.Create(ctx, tx, user); err != nil {
			tx.Rollback()
			log.Println("error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(50, users)
	for _, post := range posts {
		if err := postRepo.Create(ctx, post); err != nil {
			log.Println("error creating post:", err)
			return
		}
	}

	comments := generateComments(100, users, posts)
	for _, comment := range comments {
		if err := commentRepo.Create(ctx, comment); err != nil {
			log.Println("error creating comment:", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateUsers(num int) []*entity.User {
	users := make([]*entity.User, num)

	for i := 0; i < num; i++ {
		users[i] = &entity.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*entity.User) []*entity.Post {
	posts := make([]*entity.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.IntN(len(users))]

		posts[i] = &entity.Post{
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

func generateComments(num int, users []*entity.User, posts []*entity.Post) []*entity.Comment {
	cms := make([]*entity.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &entity.Comment{
			PostID:  posts[rand.IntN(len(posts))].ID,
			UserID:  users[rand.IntN(len(users))].ID,
			Content: comments[rand.IntN(len(comments))],
		}
	}

	return cms
}

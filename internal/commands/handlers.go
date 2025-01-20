package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/znataniel/gator/internal/database"
	"github.com/znataniel/gator/internal/rss"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	if _, err := s.Db.GetUserByName(context.Background(),
		sql.NullString{
			String: cmd.Args[0],
			Valid:  true,
		}); err == sql.ErrNoRows {
		return fmt.Errorf("error: user is not registered, register and try again")
	}

	err := s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	println("user", cmd.Args[0], "has logged in")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	if _, err := s.Db.GetUserByName(context.Background(),
		sql.NullString{
			String: cmd.Args[0],
			Valid:  true,
		}); err != sql.ErrNoRows {
		return fmt.Errorf("error: user already exists")
	}

	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      sql.NullString{String: cmd.Args[0], Valid: true},
	}
	s.Db.CreateUser(context.Background(), createUserParams)

	err := s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User was created")
	fmt.Println("id:", createUserParams.ID)
	fmt.Println("created_at:", createUserParams.CreatedAt)
	fmt.Println("updated_at:", createUserParams.UpdatedAt)
	fmt.Println("name:", cmd.Args[0])
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())
	return err
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("user could not be retrieved: %s", err)
	}

	for _, u := range users {
		if s.Cfg.CurrentUser == u.Name.String {
			fmt.Println("\t*", u.Name.String, "(current)")
			continue
		}
		fmt.Println("\t*", u.Name.String)
	}
	return nil
}

func scrapeFeeds(s *State) error {
	ctx := context.Background()

	nextFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("could not retrieve next feed to fetch: %s", err)
	}

	err = s.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: nextFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %s", err)
	}

	feed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %s", err)
	}

	if len(feed.Channel.Item) == 0 {
		fmt.Println("no posts available in current feed")
		return nil
	}

	for _, item := range feed.Channel.Item {
		fmt.Println("*", item.Title)
	}
	return nil

}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not parse time string: %s", err)
	}

	fmt.Println("collecting feeds every", cmd.Args[0])

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		fmt.Println()
		fmt.Println("available posts:")
		scrapeFeeds(s)
	}
}

func HandlerAddFeed(s *State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	ctx := context.Background()

	feed, err := s.Db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %s", err)
	}

	fmt.Println("New feed created!")
	fmt.Println(feed)

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: currentUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %s", err)
	}

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feedData, err := s.Db.GetFeedsToPrint(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feedData {
		fmt.Println("*", f.Name)
		fmt.Println("\turl:", f.Url)
		fmt.Println("\tadded by:", f.Name_2.String)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	feedData, err := s.Db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not retrieve provided rss feed data: %s", err)
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: currentUser.ID,
		FeedID: feedData.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %s", err)
	}

	fmt.Println("the feed:", feedFollow.FeedName)
	fmt.Println("has been followed by current user:", feedFollow.UserName.String)
	fmt.Println()
	return nil
}

func HandlerFollowing(s *State, cmd Command, currentUser database.User) error {
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("could not retrieve current user's followed feeds: %s", err)
	}

	fmt.Println("followed feeds")
	for _, f := range feeds {
		fmt.Println("\t*", f.FeedName, "added by", f.UserName.String)
	}

	return nil
}

func HandlerUnfollow(s *State, cmd Command, currentUser database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	feedData, err := s.Db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not retrieve provided rss feed data: %s", err)
	}

	err = s.Db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: currentUser.ID,
		FeedID: feedData.ID,
	})
	if err != nil {
		return fmt.Errorf("could not delete follow record: %s", err)
	}

	fmt.Println("you have unfollowed", feedData.Name)
	return nil
}

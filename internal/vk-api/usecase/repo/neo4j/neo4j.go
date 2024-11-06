package neo4j

import (
	"context"
	"github.com/Nimartemoff/vk-api/internal/vk-api/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserNeo4jRepo struct {
	session neo4j.SessionWithContext
}

func NewUserNeo4jRepo(session neo4j.SessionWithContext) *UserNeo4jRepo {
	return &UserNeo4jRepo{session: session}
}

func (r *UserNeo4jRepo) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.session.Run(ctx,
		"CREATE (u:User {id: $id, screen_name: $screen_name, name: $name, sex: $sex})",
		map[string]interface{}{
			"id":          user.ID,
			"screen_name": user.ScreenName,
			"name":        user.FirstName + " " + user.LastName,
			"sex":         user.Sex,
			//"city.title":  user.City.Title,
		},
	)
	return err
}

func (r *UserNeo4jRepo) CreateGroup(ctx context.Context, group models.Group) error {
	_, err := r.session.Run(ctx,
		"CREATE (g:Group {id: $id, name: $name, screen_name: $screen_name})",
		map[string]interface{}{
			"id":          group.ID,
			"name":        group.Name,
			"screen_name": group.ScreenName,
		},
	)
	return err
}

func (r *UserNeo4jRepo) CreateFollowRelationship(ctx context.Context, follower models.User, followee models.User) error {
	_, err := r.session.Run(ctx,
		"MATCH (f:User {id: $followerId}), (e:User {id: $followeeId}) CREATE (f)-[:Follow]->(e)",
		map[string]interface{}{
			"followerId": follower.ID,
			"followeeId": followee.ID,
		},
	)
	return err
}

func (r *UserNeo4jRepo) CreateSubscribeUserUserRelationship(ctx context.Context, subscriber models.User, subscribed models.User) error {
	_, err := r.session.Run(ctx,
		"MATCH (s:User {id: $subscriberId}), (u:User {id: $subscribedId}) CREATE (s)-[:Subscribe]->(u)",
		map[string]interface{}{
			"subscriberId": subscriber.ID,
			"subscribedId": subscribed.ID,
		},
	)
	return err
}

func (r *UserNeo4jRepo) CreateSubscribeUserGroupRelationship(ctx context.Context, user models.User, group models.Group) error {
	_, err := r.session.Run(ctx,
		"MATCH (u:User {id: $userId}), (g:Group {id: $groupId}) CREATE (u)-[:Subscribe]->(g)",
		map[string]interface{}{
			"userId":  user.ID,
			"groupId": group.ID,
		},
	)
	return err
}

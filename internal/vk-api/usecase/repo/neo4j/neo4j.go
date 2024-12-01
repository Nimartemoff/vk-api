package neo4j

import (
	"context"
	"fmt"
	"github.com/Nimartemoff/vk-api/internal/vk-api/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
	"strings"
)

const (
	srcTypeUnprocessed = 0
	srcTypeUser        = 1
	srcTypeGroup       = 2
)

type UserNeo4jRepo struct {
	session neo4j.SessionWithContext
}

func NewUserNeo4jRepo(session neo4j.SessionWithContext) *UserNeo4jRepo {
	return &UserNeo4jRepo{session: session}
}

func (r *UserNeo4jRepo) CreateIndexes(ctx context.Context) error {
	_, err := r.session.Run(ctx, "CREATE CONSTRAINT FOR (u:User) REQUIRE u.id IS UNIQUE", map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = r.session.Run(ctx, "CREATE CONSTRAINT FOR (u:Group) REQUIRE u.id IS UNIQUE", map[string]interface{}{})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserNeo4jRepo) CreateUser(ctx context.Context, user models.User) error {
	log.Debug().Msgf("Создание пользователя %s", user.FirstName+" "+user.LastName)
	_, err := r.session.Run(ctx,
		"MERGE (u:User {id: $id}) "+
			"SET u.screen_name = $screen_name, u.name = $name, u.sex = $sex, u.city = $city",
		map[string]interface{}{
			"id":          user.ID,
			"screen_name": user.ScreenName,
			"name":        user.FirstName + " " + user.LastName,
			"sex":         user.Sex,
			"city":        user.City.Title,
		},
	)
	return err
}

func (r *UserNeo4jRepo) CreateGroup(ctx context.Context, group models.Group) error {
	log.Debug().Msgf("Создание группы %+v", group.Name)
	_, err := r.session.Run(ctx,
		"MERGE (g:Group {id: $id}) "+
			"SET g.name = $name, g.screen_name = $screen_name ",
		map[string]interface{}{
			"id":          group.ID,
			"name":        group.Name,
			"screen_name": group.ScreenName,
		},
	)
	return err
}

func (r *UserNeo4jRepo) CreateFollowRelationship(ctx context.Context, follower models.User, followee models.User) error {
	log.Debug().Msgf("Создание фоллов связи follower: %+v - followee: %+v", follower.FirstName+" "+follower.LastName, followee.FirstName+" "+followee.LastName)
	_, err := r.session.Run(ctx,
		"MATCH (f:User {id: $followerId}), (e:User {id: $followeeId}) "+
			"MERGE (f)-[:Follow]->(e)",
		map[string]interface{}{
			"followerId": follower.ID,
			"followeeId": followee.ID,
		},
	)
	return err
}

func (r *UserNeo4jRepo) CreateSubscribeUserUserRelationship(ctx context.Context, subscriber models.User, subscribed models.User) error {
	log.Debug().Msgf("Создание subscribe связи subscriber: %+v - subscribed: %+v", subscriber.FirstName+" "+subscriber.LastName, subscribed.FirstName+" "+subscribed.LastName)
	_, err := r.session.Run(ctx,
		"MATCH (s:User {id: $subscriberId}), (u:User {id: $subscribedId}) "+
			"MERGE (s)-[:Subscribe]->(u)",
		map[string]interface{}{
			"subscriberId": subscriber.ID,
			"subscribedId": subscribed.ID,
		},
	)
	return err
}

func (r *UserNeo4jRepo) CreateSubscribeUserGroupRelationship(ctx context.Context, user models.User, group models.Group) error {
	log.Debug().Msgf("Создание связи user: %+v - group: %+v", user.FirstName+" "+user.LastName, group.Name)
	_, err := r.session.Run(ctx,
		"MATCH (u:User {id: $userId}), (g:Group {id: $groupId}) "+
			"MERGE (u)-[:Subscribe]->(g)",
		map[string]interface{}{
			"userId":  user.ID,
			"groupId": group.ID,
		},
	)
	return err
}

func (r *UserNeo4jRepo) DeleteNode(ctx context.Context, id uint64) error {
	query := `
		MATCH (n)
		WHERE id(n) = $nodeId
		DETACH DELETE n
	`
	_, err := r.session.Run(ctx, query, map[string]interface{}{"nodeId": id})
	return err
}

func (r *UserNeo4jRepo) GetUsersCount(ctx context.Context) (int, error) {
	result, err := r.session.Run(ctx, "MATCH (u:User) RETURN COUNT(u) AS count", nil)
	if err != nil {
		return 0, err
	}

	if result.Next(ctx) {
		count, ok := result.Record().Get("count")
		if !ok {
			return 0, fmt.Errorf("failed to get users count")
		}

		if countInt, ok := count.(int64); ok {
			return int(countInt), nil
		}
		return 0, fmt.Errorf("failed to get users count")
	}

	return 0, nil
}

func (r *UserNeo4jRepo) GetGroupsCount(ctx context.Context) (int, error) {
	result, err := r.session.Run(ctx, "MATCH (u:Group) RETURN COUNT(u) AS count", nil)
	if err != nil {
		return 0, err
	}

	if result.Next(ctx) {
		count, ok := result.Record().Get("count")
		if !ok {
			return 0, fmt.Errorf("failed to get users count")
		}

		if countInt, ok := count.(int64); ok {
			return int(countInt), nil
		}
		return 0, fmt.Errorf("failed to get users count")
	}

	return 0, nil
}

func (r *UserNeo4jRepo) GetTopUsersByFollowersCount(ctx context.Context, limit int) ([]models.User, error) {
	query := `
		MATCH (u:User)<-[:Follow]-(f:User)
		RETURN u.id AS id, u.screen_name AS screen_name, u.name AS name, u.sex AS sex, u.city AS city, COUNT(f) AS followersCount 
		ORDER BY followersCount DESC
		LIMIT $limit
	`
	result, err := r.session.Run(ctx, query, map[string]interface{}{"limit": limit})
	if err != nil {
		return nil, err
	}

	var users []models.User
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		screenName, _ := record.Get("screen_name")
		name, _ := record.Get("name")
		sex, _ := record.Get("sex")
		city, _ := record.Get("city")

		nameParts := strings.Split(name.(string), " ")
		var firstName, lastName string
		if len(nameParts) > 0 {
			firstName = nameParts[0]
		}
		if len(nameParts) > 1 {
			lastName = nameParts[1]
		}

		users = append(users, models.User{
			ID:         uint64(id.(int64)),
			ScreenName: screenName.(string),
			FirstName:  firstName,
			LastName:   lastName,
			Sex:        byte(sex.(int64)),
			City: models.City{
				Title: city.(string),
			},
		})
	}

	return users, nil
}

func (r *UserNeo4jRepo) GetTopGroupsBySubscribersCount(ctx context.Context, limit int) ([]models.Group, error) {
	query := `
		MATCH (g:Group)<-[:Subscribe]-(u:User)
		RETURN g.id AS group_id, g.name AS name, g.screen_name AS screen_name, COUNT(u) AS subscribersCount 
		ORDER BY subscribersCount DESC
		LIMIT $limit
	`
	result, err := r.session.Run(ctx, query, map[string]interface{}{"limit": limit})
	if err != nil {
		return nil, err
	}

	var groups []models.Group
	for result.Next(ctx) {
		record := result.Record()
		groupID, _ := record.Get("group_id")
		name, _ := record.Get("name")
		screenName, _ := record.Get("screen_name")

		groups = append(groups, models.Group{
			ID:         uint64(groupID.(int64)),
			Name:       name.(string),
			ScreenName: screenName.(string),
		})
	}

	return groups, nil
}

func (r *UserNeo4jRepo) GetUsersWithDifferentGroups(ctx context.Context, limit int) ([]models.User, error) {
	query := `
		MATCH (g:Group)<-[:Subscribe]-(u:User)
		WITH u
		MATCH (g)<-[:Subscribe]-(other:User)
		WHERE u <> other
		RETURN DISTINCT u.id AS id, u.screen_name AS screen_name, u.name AS name, u.sex AS sex, u.city AS city
		LIMIT $limit;
	`
	result, err := r.session.Run(ctx, query, map[string]interface{}{"limit": limit})
	if err != nil {
		return nil, err
	}

	var users []models.User
	for result.Next(ctx) {
		record := result.Record()
		id, _ := record.Get("id")
		screenName, _ := record.Get("screen_name")
		name, _ := record.Get("name")
		sex, _ := record.Get("sex")
		city, _ := record.Get("city")

		nameParts := strings.Split(name.(string), " ")
		var firstName, lastName string
		if len(nameParts) > 0 {
			firstName = nameParts[0]
		}
		if len(nameParts) > 1 {
			lastName = nameParts[1]
		}

		users = append(users, models.User{
			ID:         uint64(id.(int64)),
			ScreenName: screenName.(string),
			FirstName:  firstName,
			LastName:   lastName,
			Sex:        byte(sex.(int64)),
			City: models.City{
				Title: city.(string),
			},
		})
	}

	return users, nil
}

func (r *UserNeo4jRepo) GetAllNodes(ctx context.Context) ([]models.Node, error) {
	query := `
		MATCH (n)
		RETURN id(n) AS id, labels(n) AS labels
	`
	result, err := r.session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var nodes []models.Node
	for result.Next(ctx) {
		record := result.Record()
		nodeID, _ := record.Get("id")
		nodeLabels, _ := record.Get("labels")

		var labels []string
		for _, label := range nodeLabels.([]interface{}) {
			labels = append(labels, label.(string))
		}

		nodes = append(nodes, models.Node{
			ID:     nodeID.(int64),
			Labels: labels,
		})
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *UserNeo4jRepo) GetNodeWithRelationships(ctx context.Context, id uint64) (interface{}, error) {
	query := `
		OPTIONAL MATCH (n)-[r]->(m) WHERE id(n) = $nodeID 
		WITH n, r, m 
		WHERE n IS NOT NULL  
		RETURN n, r, m, null as r2 
		UNION
		OPTIONAL MATCH (n)<-[r2]-(m) WHERE id(n) = $nodeID 
		WITH n, r2, m 
		WHERE n IS NOT NULL  
		RETURN n, null as r, m, r2 
	`
	result, err := r.session.Run(ctx, query, map[string]interface{}{"nodeID": id})
	if err != nil {
		return nil, err
	}

	var srcType byte // 0 - источник не обработан 1 - юзер 2 - группа.
	var group models.GroupWithSubscribers
	var user models.User

	for result.Next(ctx) {
		record := result.Record()

		node, _ := record.Get("n")
		n, ok := node.(neo4j.Node)
		if !ok {
			return nil, fmt.Errorf("cant assert node %+v (type %T) to neo4j.Node", node, node)
		}

		if len(n.Labels) == 0 {
			return nil, fmt.Errorf("can not find labels for node %+v", node)
		}

		if srcType == srcTypeUnprocessed {
			switch n.Labels[0] {
			case "User":
				user = processUserNode(n)
				srcType = srcTypeUser
			case "Group":
				group = processGroupWithSubscribersNode(n)
				srcType = srcTypeGroup
			}
		}

		connectedNode, _ := record.Get("m")
		m, ok := connectedNode.(neo4j.Node)
		if !ok {
			return nil, fmt.Errorf("cant assert node %+v (type %T) to neo4j.Node", connectedNode, connectedNode)
		}

		switch srcType {
		case srcTypeUser:
			relationship, _ := record.Get("r")
			r, ok := relationship.(neo4j.Relationship)
			if ok {
				processUserRelation(&user, m, r)
			}

			relationship, _ = record.Get("r2")
			r2, ok := relationship.(neo4j.Relationship)
			if ok {
				processUserRelation(&user, m, r2)
			}
		case srcTypeGroup:
			group.Subscribers = append(group.Subscribers, processUserNode(m))
		}
	}

	switch srcType {
	case srcTypeUser:
		return user, result.Err()
	case srcTypeGroup:
		return group, result.Err()
	}

	return nil, nil
}

func processUserRelation(user *models.User, m neo4j.Node, r neo4j.Relationship) {
	switch r.Type {
	case "Follow":
		user.Followers = append(user.Followers, processUserNode(m))
	case "Subscribe":
		switch m.Labels[0] {
		case "Group":
			user.Subscriptions.Groups = append(user.Subscriptions.Groups, processGroupNode(m))
		case "User":
			user.Subscriptions.Users = append(user.Subscriptions.Users, processUserNode(m))
		}
	}
}

func processUserNode(n neo4j.Node) models.User {
	var user models.User
	props := n.Props

	if id, ok := props["id"].(int64); ok {
		user.ID = uint64(id)
	}
	if screenName, ok := props["screen_name"].(string); ok {
		user.ScreenName = screenName
	}
	if sex, ok := props["sex"].(int64); ok {
		user.Sex = byte(sex)
	}
	if city, ok := props["city"].(string); ok {
		user.City = models.City{Title: city}
	}
	if name, ok := props["name"].(string); ok {
		user.FirstName, user.LastName = splitFullName(name)
	}

	return user
}

func processGroupNode(n neo4j.Node) models.Group {
	var group models.Group
	props := n.Props

	if id, ok := props["id"].(int64); ok {
		group.ID = uint64(id)
	}
	if name, ok := props["name"].(string); ok {
		group.Name = name
	}
	if screenName, ok := props["screen_name"].(string); ok {
		group.ScreenName = screenName
	}

	return group
}

func processGroupWithSubscribersNode(n neo4j.Node) models.GroupWithSubscribers {
	var group models.GroupWithSubscribers
	props := n.Props

	if id, ok := props["id"].(int64); ok {
		group.ID = uint64(id)
	}
	if name, ok := props["name"].(string); ok {
		group.Name = name
	}
	if screenName, ok := props["screen_name"].(string); ok {
		group.ScreenName = screenName
	}

	return group
}

func splitFullName(fullName string) (string, string) {
	nameParts := strings.Fields(fullName)
	var firstName, lastName string
	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
	}
	return firstName, lastName
}

package app

import (
	"context"
	"github.com/Nimartemoff/vk-api/cmd/vk-api/config"
	v1 "github.com/Nimartemoff/vk-api/internal/vk-api/controller/http/v1"
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase"
	neo4jRepo "github.com/Nimartemoff/vk-api/internal/vk-api/usecase/repo/neo4j"
	"github.com/Nimartemoff/vk-api/internal/vk-api/usecase/rest"
	"github.com/Nimartemoff/vk-api/pkg/httpserver"
	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/zerolog/log"
)

const (
	userID   = 183170347
	fileName = "Копырин.json"
)

func Run(cfg *config.Config) {
	c := rest.NewVKClient(cfg.VKAPI.URLs, cfg.VKAPI.Token)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ContextTimeout)
	defer cancel()

	driver, err := neo4j.NewDriverWithContext(cfg.Neo4j.URL, neo4j.NoAuth())
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	defer driver.Close(ctx)

	session := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
	})
	defer session.Close(ctx)

	userUsecase := usecase.NewUserUsecase(c, neo4jRepo.NewUserNeo4jRepo(session))

	r := chi.NewRouter()
	v1.NewRouter(cfg, r, userUsecase)

	if err := httpserver.New(r, cfg.API.Port); err != nil {
		log.Error().Err(err).Msgf("could not start router http server at port: %s ", cfg.API.Port)
	}

	//if err := userUsecase.CreateIndexes(ctx); err != nil {
	//	log.Error().Err(err).Send()
	//}
	//
	//user, err := userUsecase.GetUsersWithDepth(userID, 3)
	//if err != nil {
	//	log.Error().Err(err).Send()
	//	return
	//}
	//
	//if err := userUsecase.SaveUser(ctx, user); err != nil {
	//	log.Error().Err(err).Send()
	//	return
	//}
	//
	//queryType := flag.String("query_type", "0", "Выбор запроса из параметров командной строки")
	//flag.Parse()
	//
	//switch *queryType {
	//case "0":
	//	count, err := userUsecase.GetUsersCount(ctx)
	//	if err != nil {
	//		log.Error().Err(err).Send()
	//		return
	//	}
	//
	//	log.Info().Msgf("Количество пользователей: %d", count)
	//case "1":
	//	count, err := userUsecase.GetGroupsCount(ctx)
	//	if err != nil {
	//		log.Error().Err(err).Send()
	//		return
	//	}
	//
	//	log.Info().Msgf("Количество групп: %d", count)
	//case "2":
	//	users, err := userUsecase.GetTopUsersByFollowersCount(ctx, 5)
	//	if err != nil {
	//		log.Error().Err(err).Send()
	//		return
	//	}
	//
	//	log.Info().Msgf("Топ 5 пользователей по числу фолловеров: %+v", users)
	//case "3":
	//	groups, err := userUsecase.GetTopGroupsBySubscribersCount(ctx, 5)
	//	if err != nil {
	//		log.Error().Err(err).Send()
	//		return
	//	}
	//
	//	log.Info().Msgf("Топ 5 групп по числу подписок: %+v", groups)
	//case "4":
	//	users, err := userUsecase.GetUsersWithDifferentGroups(ctx, 5)
	//	if err != nil {
	//		log.Error().Err(err).Send()
	//		return
	//	}
	//
	//	log.Info().Msgf("Пользователи с непересекающимися группами в подписаках: %+v", users)
	//}
	//
	//file, err := os.Create(fileName)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer file.Close()
	//
	//encoder := json.NewEncoder(file)
	//encoder.SetIndent("", "    ")
	//if err := encoder.Encode(user); err != nil {
	//	log.Error().Err(err).Send()
	//	return
	//}
	//
	//log.Info().Msgf("Завершение работы программы")
}

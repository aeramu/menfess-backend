package cmd

import (
	"context"
	"github.com/aeramu/menfess-backend/infra/graphql"
	"github.com/aeramu/menfess-backend/modules/auth"
	logModule "github.com/aeramu/menfess-backend/modules/log"
	"github.com/aeramu/menfess-backend/modules/notification"
	"github.com/aeramu/menfess-backend/modules/post"
	"github.com/aeramu/menfess-backend/modules/user"
	"github.com/aeramu/menfess-backend/service"
	"github.com/aeramu/menfess-backend/utils/playground"
	"github.com/aeramu/mongolib"
	"log"
	"net/http"
)

func Run() {
	client, err := mongolib.NewSingletonClient(context.Background(), "mongodb+srv://admin:admin@qiup-wrbox.mongodb.net")
	if err != nil {
		log.Fatalln("[Init DB Client]", err)
	}
	db := mongolib.NewDatabase(client, "menfess-staging")
	adapter := service.Adapter{
		UserModule:         user.NewUserModule(db),
		PostModule:         post.NewPostModule(db),
		AuthModule:         auth.NewAuthModule(),
		NotificationModule: notification.NewNotificationModule(db),
		LogModule:          logModule.NewLogModule(),
	}
	svc := service.NewService(adapter)
	srv, err := graphql.NewServer(svc)
	if err != nil {
		log.Fatalln("[Init Server]", err)
	}
	pg := playground.Handler("Playground", "/")
	http.Handle("/", srv)
	http.Handle("/playground", pg)
	log.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

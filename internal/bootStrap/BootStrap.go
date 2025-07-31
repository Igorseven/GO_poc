package bootStrap

import (
	config "PocGo/internal/configuration"
	constant "PocGo/internal/domain/constants"
	notify "PocGo/internal/domain/notification"
	repository "PocGo/internal/repositories"
	applicationServer "PocGo/internal/server"
	service "PocGo/internal/services"
	dataBaseConnection "PocGo/pkg/database"
	"context"
	dbProvider "database/sql"
	logger "log"
	"sync"
	"time"
)

type Application struct {
	Server        *applicationServer.ApplicationServer
	Configuration *config.Config
	dataBase      *dbProvider.DB
	schedulerDone chan struct{}
	schedulerWg   sync.WaitGroup
}

func NewApplication() *Application {
	configuration := config.LoadConfig("development")
	dbInstance, err := dataBaseConnection.NewConnection(configuration)

	if err != nil {
		logger.Fatalf(notify.ErrorDbFatal, err)
	}

	defer func(dbInstance *dataBaseConnection.Database) {
		err := dbInstance.Close()
		if err != nil {

		}
	}(dbInstance)

	application := &Application{
		Configuration: configuration,
		dataBase:      dbInstance.GetConnection(),
		schedulerDone: make(chan struct{}),
	}

	dataBase, _ := application.setupDatabase()
	repositories, _ := application.setupRepositories(dataBase)
	services := application.setupServices(repositories)
	server := application.setupServer(services)

	application.Server = server
	return application
}

func (app *Application) setupDatabase() (*dbProvider.DB, error) {
	dataBase, err := dataBaseConnection.NewConnection(app.Configuration)
	if err != nil {
		logger.Fatalf(notify.ErrorDbFatal, err)
	}
	return dataBase.GetConnection(), err
}

func (app *Application) setupRepositories(db *dbProvider.DB) (*repository.Repositories, error) {
	return repository.NewRepositories(db)
}

func (app *Application) setupServices(repos *repository.Repositories) *service.Services {
	return service.NewServices(repos)
}

func (app *Application) setupServer(services *service.Services) *applicationServer.ApplicationServer {
	return applicationServer.NewServer(services)
}

func (app *Application) runUpdateOldUsersStatus() {
	repos, err := app.setupRepositories(app.dataBase)
	if err != nil {
		logger.Printf(notify.ErrorRepositoryFatal, err)
		return
	}

	services := app.setupServices(repos)
	count, err := services.User.UpdateOldUsersStatus()
	if err != nil {
		logger.Printf(notify.LogForErrorUpdateUsers, err)
	} else {
		logger.Printf(notify.LogForPartialUpdateUsers, count)
	}
}

func (app *Application) startDailyScheduler(cfg config.Config) {
	app.schedulerWg.Add(1)

	go func() {
		defer app.schedulerWg.Done()

		app.runUpdateOldUsersStatus()

		now := time.Now()
		nextMidnight := time.Date(
			now.Year(),
			now.Month(),
			now.Day()+cfg.Routine.IncrementDay,
			cfg.Routine.Hour,
			cfg.Routine.Minute,
			cfg.Routine.Second,
			cfg.Routine.Millisecond,
			now.Location())

		initialDelay := nextMidnight.Sub(now)

		logger.Printf(notify.LogStartRotineAction,
			nextMidnight.Format(constant.FormatDate),
			initialDelay)

		select {
		case <-time.After(initialDelay):
		case <-app.schedulerDone:
			logger.Println(notify.LogRotineNoStarted)
			return
		}

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				app.runUpdateOldUsersStatus()
				nextRun := time.Now().Add(24 * time.Hour)
				logger.Printf(notify.LogNextRotine, nextRun.Format(constant.FormatDate))
			case <-app.schedulerDone:
				logger.Println(notify.LogRotineOff)
				return
			}
		}
	}()
}

func (app *Application) StopScheduler() {
	close(app.schedulerDone)
	app.schedulerWg.Wait()
	logger.Println(notify.LogRotineStoped)
}

func (app *Application) Run(ctx context.Context) error {
	app.startDailyScheduler(*app.Configuration)

	return app.Server.Start(ctx)
}

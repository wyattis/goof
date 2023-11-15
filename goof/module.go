package goof

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"

	"github.com/wyattis/goof/http/middleware"
	"github.com/wyattis/goof/log"
	"github.com/wyattis/goof/migrate"
	"github.com/wyattis/goof/sql"
	"github.com/wyattis/goof/sql/driver"
)

type HttpConfig struct {
	Addr string
	Host string
	SSL  struct {
		Enabled  bool
		CertFile string
		KeyFile  string
	}
}

type RootConfig struct {
	Production bool
	DB         driver.Config
	Http       HttpConfig
	Log        log.Config
}

type ControllerConfig struct {
	Addr string
	Host string
}

type Controller interface {
	Init(config ControllerConfig) (err error)
	MountHTTP(router gin.IRouter) (err error)
}

type BaseController struct{}

func (c *BaseController) Init(config ControllerConfig) (err error) {
	return
}

func (c *BaseController) MountHTTP(router gin.IRouter) (err error) {
	return
}

type ModuleApi interface {
	AddController(controllers ...Controller)
	AddMigration(migrations ...migrate.Migration)
	GetDB() (*sqlx.DB, error)
	GetSessionStore() (sessions.Store, error)
}

type ControllersModule interface {
	Controllers(db *sqlx.DB) []Controller
}

type MigrationsModule interface {
	Migrations() []migrate.Migration
}

type Module interface {
	Id() string
	PreInit(api ModuleApi, config any) (err error)
	Init(api ModuleApi, config any) (err error)
	PostInit(api ModuleApi, config any) (err error)
	Close() (err error)
}

type ModuleBase struct{}

func (m *ModuleBase) PreInit(api ModuleApi, config interface{}) (err error) {
	return nil
}

func (m *ModuleBase) Init(api ModuleApi, config interface{}) (err error) {
	return nil
}

func (m *ModuleBase) PostInit(api ModuleApi, config interface{}) (err error) {
	return nil
}

func (m *ModuleBase) Close() (err error) {
	return nil
}

type moduleDef struct {
	module       Module
	sessionStore sessions.Store
	config       interface{}
	controllers  []Controller
	migrations   []migrate.Migration
	db           *sqlx.DB
}

func (m *moduleDef) AddMigration(migrations ...migrate.Migration) {
	m.migrations = append(m.migrations, migrations...)
}

func (m *moduleDef) AddController(controllers ...Controller) {
	m.controllers = append(m.controllers, controllers...)
}

func (m *moduleDef) GetSessionStore() (sessions.Store, error) {
	if m.sessionStore == nil {
		return nil, fmt.Errorf("Session store has not been initialized at %s", m.module.Id())
	}
	return m.sessionStore, nil
}

func (m *moduleDef) GetDB() (db *sqlx.DB, err error) {
	if m.db == nil {
		err = fmt.Errorf("DB has not been initialized at %s", m.module.Id())
		return
	}
	return m.db, nil
}

type RootModule struct {
	Config RootConfig

	hasInitialized bool
	modules        []*moduleDef
	middleware     []gin.HandlerFunc
	engine         *gin.Engine
	db             *sqlx.DB
	sessionStore   sessions.Store
}

// Add a module to the root module. Modules are initialized in the order they are added.
func (r *RootModule) Add(modules ...Module) {
	for _, module := range modules {
		r.modules = append(r.modules, &moduleDef{
			module: module,
		})
	}
}

// Add a middleware to every route in under the root module. Middlewares are applied in the order they are added.
func (r *RootModule) AddMiddleware(middleware ...gin.HandlerFunc) {
	r.middleware = append(r.middleware, middleware...)
}

// Init the root module. All modules should be added before calling Init.
func (r *RootModule) Init() (err error) {
	if r.hasInitialized {
		return
	}
	r.hasInitialized = true
	if err = r.initRoot(); err != nil {
		return fmt.Errorf("Failed to init root module:\n %w", err)
	}
	for _, m := range r.modules {
		m.sessionStore = r.sessionStore
		if err = m.module.PreInit(m, m.config); err != nil {
			return fmt.Errorf("Failed to PreInit module %s:\n %w", m.module.Id(), err)
		}
		if mm, ok := m.module.(MigrationsModule); ok {
			m.AddMigration(mm.Migrations()...)
		}
	}
	if err = r.initDatabase(); err != nil {
		return fmt.Errorf("Failed to init database:\n %w", err)
	}
	if err = r.runMigrations(); err != nil {
		return fmt.Errorf("Failed to run migrations:\n %w", err)
	}
	for _, m := range r.modules {
		m.db = r.db
		if err = m.module.Init(m, m.config); err != nil {
			return fmt.Errorf("Failed to Init module %s:\n %w", m.module.Id(), err)
		}
	}
	if err = r.initControllers(); err != nil {
		return fmt.Errorf("Failed to init controllers:\n %w", err)
	}
	for _, m := range r.modules {
		if err = m.module.PostInit(m, m.config); err != nil {
			return fmt.Errorf("Failed to PostInit module %s:\n %w", m.module.Id(), err)
		}
	}
	return
}

func (r *RootModule) initRoot() (err error) {
	if err = log.Init(r.Config.Log); err != nil {
		return fmt.Errorf("Failed to init log:\n %w", err)
	}
	migrate.SetLogger(&log.Logger)
	// TODO: make session store configurable
	r.sessionStore = sessions.NewCookieStore([]byte("secret"))
	gin.SetMode(gin.ReleaseMode)
	r.engine = gin.Default()
	r.engine.Use(r.middleware...)
	if !r.Config.Production {
		r.engine.Use(middleware.CORS())
	} else {
		r.printRoutes()
	}
	return
}

func (r *RootModule) printRoutes() {
	for _, r := range r.engine.Routes() {
		fmt.Printf("%s %s -> %v\n", r.Method, r.Path, r.HandlerFunc)
	}
}

func (r *RootModule) initControllers() (err error) {
	for _, m := range r.modules {
		if c, ok := m.module.(ControllersModule); ok {
			m.AddController(c.Controllers(r.db)...)
		}
		for _, c := range m.controllers {
			if err = c.Init(ControllerConfig{
				Addr: r.Config.Http.Addr,
				Host: r.Config.Http.Host,
			}); err != nil {
				return fmt.Errorf("Failed to Init controller from module %s:\n %w", m.module.Id(), err)
			}
		}
		for _, c := range m.controllers {
			if err = c.MountHTTP(r.engine); err != nil {
				return fmt.Errorf("Failed to Init controller from module %s:\n %w", m.module.Id(), err)
			}
		}
	}
	return
}

func (r *RootModule) initDatabase() (err error) {
	log.Debug().Interface("config", r.Config.DB).Msg("opening database")
	db, err := sql.Open(r.Config.DB)
	if err != nil {
		return fmt.Errorf("Failed to open database:\n %w", err)
	}
	r.db = sqlx.NewDb(db, r.Config.DB.DriverName.String())
	err = r.db.Ping()
	log.Debug().Err(err).Msg("pinged database")
	return
}

func (r *RootModule) runMigrations() (err error) {
	log.Debug().Msg("preparing migrations")
	migrations := []migrate.Migration{}
	var version uint = 1
	for _, m := range r.modules {
		for _, migration := range m.migrations {
			migration.Version = version
			migrations = append(migrations, migration)
			version++
		}
	}
	targetVersion := version - 1
	log.Debug().Msgf("running migrations up to version %d", targetVersion)
	return migrate.MigrateUpTo(migrations, r.db.DB, r.Config.DB.DriverName, r.Config.DB.Database, targetVersion)
}

// Start the server. This will initialize the server if it has not already been initialized.
func (r *RootModule) Run() (err error) {
	if err = r.Init(); err != nil {
		return fmt.Errorf("Failed to init server:\n %w", err)
	}
	if r.Config.Http.SSL.Enabled {
		log.Info().Msgf("Starting TLS server on '%s'", r.Config.Http.Addr)
		return r.engine.RunTLS(r.Config.Http.Addr, r.Config.Http.SSL.CertFile, r.Config.Http.SSL.KeyFile)
	}
	log.Info().Msgf("Starting server on '%s'", r.Config.Http.Addr)
	return r.engine.Run(r.Config.Http.Addr)
}

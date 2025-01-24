package database

import (
	"fmt"
	"net/url"
	mLogger "template/internal/mlogger"
	"template/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormDatabase struct {
	db *gorm.DB
}

var log = mLogger.New("DB", true, mLogger.VERBOSE)

func (obj *GormDatabase) Init(databaseURL string) error {
	// Parse the URL
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return err
	}

	switch parsedURL.Scheme {
	case "mysql":
		password, _ := parsedURL.User.Password()

		// Construct DSN for MySQL
		dsn := parsedURL.User.Username() + ":" + password +
			"@tcp(" + parsedURL.Host + ")" + parsedURL.Path

		// Open MySQL database
		obj.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.E("Failed to connect to the MySQL database: %v", err)
			return err
		}
		log.I("Successfully connected to the MySQL database")

	case "postgres":
		// Construct DSN for PostgreSQL
		password, _ := parsedURL.User.Password()
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			parsedURL.Hostname(),
			parsedURL.User.Username(),
			password,
			parsedURL.Path[1:]) // Remove the leading "/" from the path

		// Open PostgreSQL database
		obj.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.E("Failed to connect to the PostgreSQL database: %v", err)
			return err
		}
		log.I("Successfully connected to the PostgreSQL database")

	case "sqlite":
		// SQLite database uses the path as the database file
		dbFile := parsedURL.Path[1:] // Remove the leading "/" from the path

		// Open SQLite database
		obj.db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to the SQLite database: %v", err)
			return err
		}
		log.I("Successfully connected to the SQLite database")

	default:
		error_str := fmt.Sprintf("unsupported database scheme: %v", parsedURL.Scheme)
		log.E(error_str)
		return fmt.Errorf(error_str)
	}

	// Automigrate the tables
	err = obj.db.AutoMigrate(&models.User{})
	if err != nil {
		log.C("Failed to Automigrate database : %v", err)
		return err
	}
	return nil
}

func (obj *GormDatabase) DeInit(){
	
}



func (g *GormDatabase) AllUsers() ([]*models.User, error) {
	var users []*models.User
	err := g.db.Find(&users).Error // Fetch all users
	if err != nil {
		return nil, fmt.Errorf("error fetching all users: %v", err)
	}
	return users, nil
}

func (g *GormDatabase) GetUserById(id uint64) (*models.User, error) {
	var user models.User
	err := g.db.First(&user, id).Error // Fetch user by primary key (ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching user by ID: %v", err)
	}
	return &user, nil
}

func (g *GormDatabase) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := g.db.Where("email = ?", email).First(&user).Error // Fetch user by email
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching user by email: %v", err)
	}
	return &user, nil
}

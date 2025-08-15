package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/appcheck"
	"firebase.google.com/go/v4/auth"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/auth"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"musicfree.root/handlers"
)

func FirebaseInit() FireAuth {

	/*	provider, error := oidc.NewProvider(context.TODO(), "https://accounts.google.com")
		if error != nil {
			fmt.Print(error.Error())
		}

		conf := &jwt.Config{
			Email: "xxx@developer.gserviceaccount.com",

			PrivateKey: []byte("-----BEGIN RSA PRIVATE KEY-----..."),
			Scopes: []string{
				"https://www.googleapis.com/auth/bigquery",
				"https://www.googleapis.com/auth/blogger",
			},
			TokenURL: google.JWTTokenURL,
			// If you would like to impersonate a user, you can
			// create a transport with a subject. The following GET
			// request will be made on the behalf of user@example.com.
			// Optional.
			Subject: "user@example.com",
		}

		// Get these from your Auth0 Application Dashboard.
		domain := "example.us.auth0.com"
		clientID := "EXAMPLE_16L9d34h0qe4NVE6SaHxZEid"
		clientSecret := "EXAMPLE_XSQGmnt8JdXs23407hrK6XXXXXXX"

		config_oidc := &oidc.Config{
			ClientID: clientID,
		}

		oauth2Config := oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			// RedirectURL:  redirectURL,

			// Discovery returns the OAuth2 endpoints.
			Endpoint: provider.Endpoint(),

			// "openid" is a required scope for OpenID Connect flows.
			Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
		}

		return provider, oauth2Config, config_oidc*/
	domain := "example.us.auth0.com"
	clientID := "EXAMPLE_16L9d34h0qe4NVE6SaHxZEid"
	clientSecret := "EXAMPLE_XSQGmnt8JdXs23407hrK6XXXXXXX"
	oauth2Cnf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		// RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		//Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	opt := option.WithCredentialsFile("refreshToken.json")

	app, err := firebase.NewApp(context.TODO(), nil, opt)
	if err != nil {
		fmt.Print(err.Error())
	}

	app_check, err2 := app.AppCheck(context.TODO())
	if err2 != nil {
		fmt.Print(err.Error())
	}
	auth, err3 := app.Auth(context.TODO())
	if err3 != nil {
		fmt.Print(err3.Error())
	}

	return FireAuth{auth: auth, app_check: app_check}

}

func corsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		fmt.Print("here")
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	}
}

type FireAuth struct {
	auth      *auth.Client
	user_db   *sql.DB
	app_check *appcheck.Client
}

func (f_auth *FireAuth) AutenthificationHandler(w http.ResponseWriter, req *http.Request) {
	id_token := req.Header.Get("Authorization")
	token, error := f_auth.auth.VerifyIDToken(req.Context(), id_token)
	if error != nil {
		w.WriteHeader(403)
		return
	}
	is_exists_query, err4 := f_auth.user_db.QueryContext(req.Context(), `SELECT  EXISTS(SELECT id FROM users WHERE id=? LIMIT 1)`, token.UID)
	if err4 != nil {
		w.WriteHeader(500)
	}
	is_exists_query.Next()
	var is_exists int
	is_exists_query.Scan(is_exists)
	if is_exists == 0 {
		_, err5 := f_auth.user_db.ExecContext(req.Context(), `INSERT INTO users (id, email, username)(?,?,?) `, token.UID, token.Claims["email"], token.Claims["username"])
		if err5 != nil {

		}
	}
	expire := time.Hour * 24 * 7
	string_result, err1 := f_auth.auth.SessionCookie(req.Context(), id_token, expire)
	if err1 != nil {
		w.WriteHeader(500)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    string_result,
		MaxAge:   int(expire.Seconds()),
		HttpOnly: true,
		Secure:   true,
	})

}

func (f_auth *FireAuth) SessionCookieChech(next http.HandlerFunc, is_guest bool) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		cookie, error := req.Cookie("session")
		if error != nil {

			if is_guest {
				guest_cookie, err1 := req.Cookie("guest_session")
				if err1 == nil {
					ctx := context.WithValue(context.TODO(), "session_field", guest_cookie)
					next(w, req.WithContext(ctx))
				}
			}
			w.WriteHeader(403)
			fmt.Print(error.Error())
			return
		}

		decoded, err := f_auth.auth.VerifySessionCookieAndCheckRevoked(req.Context(), cookie.Value)
		if err != nil {
			w.WriteHeader(403)
			return
		}
		context := context.WithValue(context.TODO(), "id_field", decoded.Subject)
		next(w, req.WithContext(context))
	}

}

func (f_auth *FireAuth) AppCheckMiddleWare(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		_, err := f_auth.app_check.VerifyToken(req.Header.Get("X-Firebase-AppCheck"))
		if err != nil {
			w.WriteHeader(403)
			return
		}
		next(w, req)

	}

}

func CorsMiddleware(next http.HandlerFunc, request_method string) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			corsHandler(w, req)
		} else {
			if req.Method != request_method {
				w.WriteHeader(405)
			} else {
				w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
				next(w, req)
			}
		}
	}

}

func (f_auth *FireAuth) AuthenticationMiddleware(next http.HandlerFunc, method_string string, is_guest bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		CorsMiddleware(f_auth.AppCheckMiddleWare((f_auth.SessionCookieChech(next, is_guest))), method_string)
	}
}

func main() {
	fire_auth := FirebaseInit()

	var db, err5 = sql.Open("mysql", SqlInitialize().FormatDSN())

	if err5 != nil {

	}
	fire_auth.user_db = db
	env := handlers.Env{
		DB: db}

	http.HandleFunc("POST /music/create/albumn", env.AlbumnCreate)
	http.HandleFunc("/music/create/author", corsHandler)
	http.HandleFunc("/music/create/author", CorsMiddleware(env.MusicianCreate, "POST"))
	http.HandleFunc("POST /music/create/tags_genres", env.GenreTagCreation)
	http.HandleFunc("GET /music/create/search/authors/{name}/{page_size}/{page_index}", env.Musician_For_Search)
	http.HandleFunc("GET/music/create/search/{type}/{name}/{page_size}/{page_index}", env.Tags_Genres__For_Search)
	http.HandleFunc("/music/upload/albumn", CorsMiddleware(handlers.AlbumnUpload, "POST"))
	http.HandleFunc("/music/upload/id/albumn", CorsMiddleware(env.UploadedAlbumnIds, "POST"))
	http.HandleFunc(" /music/upload/id/author", CorsMiddleware(env.UploadMusicianIds, "POST"))
	http.HandleFunc("POST /music/upload/author", CorsMiddleware(handlers.MusicianCoverUpload, "POST"))
	http.HandleFunc("/music/listened/add", corsHandler)
	http.HandleFunc("/music/listened/add", fire_auth.AuthenticationMiddleware(env.AddListened, "POST", false))
	http.HandleFunc("/auth/login", CorsMiddleware(fire_auth.AppCheckMiddleWare(fire_auth.AutenthificationHandler), "POST"))

	http.ListenAndServe("localhost:8000", nil)
}

type CustomClaims struct {
	Scope string `json:"scope"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func SqlInitialize() *mysql.Config {

	cfg := mysql.NewConfig()
	//os.Getenv("DBUSER")
	//os.Getenv("DBPASS")
	cfg.User = "root"
	cfg.Passwd = "saharok2342saharok"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "recordings"
	return cfg

}

package internal

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var cookie sync.Map

const CookieName string = "forum"

type duration struct {
	expiry map[interface{}]time.Time
	mu     sync.Mutex
}

var sessionDuration = duration{expiry: make(map[interface{}]time.Time)}

func AddCookie(w http.ResponseWriter, r *http.Request, id int) {
	sessionDuration.mu.Lock()
	defer sessionDuration.mu.Unlock()

	u := uuid.NewV4().String()
	deleteExistingCookie(id, u)

	cookie.Store(u, id)
	expire := time.Now().AddDate(0, 0, 1)
	sessionDuration.expiry[u] = expire

	http.SetCookie(w,
		&http.Cookie{
			Name:     CookieName,
			Value:    u,
			Path:     "/",
			HttpOnly: true,
			Expires:  expire,
		})
}

func DeleteCookie(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(CookieName)
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
	cookie.Delete(c.Value)
}

func deleteExistingCookie(id int, uuid string) {
	cookie.Range(func(key, value interface{}) bool {
		if id == value.(int) {
			cookie.Delete(key)
		}
		return true
	})
}

func isSession(r *http.Request) bool {
	c, err := r.Cookie(CookieName)
	if err != nil {
		return false
	}
	_, ok := cookie.Load(c.Value)
	return ok
}

func GetUserIDByCookie(r *http.Request) (int, error) {
	c, err := r.Cookie(CookieName)
	if err != nil {
		return 0, err
	}
	value, ok := cookie.Load(c.Value)
	if !ok {
		return 0, fmt.Errorf("getUserIDByCookie: cannot load value from cookie store")
	}
	userID := value.(int)
	return userID, nil
}

func DeleteExpiredSession() {
	for {
		cookie.Range(func(key, value interface{}) bool {
			sessionDuration.mu.Lock()
			if time.Now().Unix() > sessionDuration.expiry[key].Unix() {
				cookie.Delete(key)
			}
			sessionDuration.mu.Unlock()
			return true
		})
		time.Sleep(time.Second * 5)
	}
}

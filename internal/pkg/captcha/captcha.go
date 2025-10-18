package captcha

import (
	"time"
	
	"github.com/mojocn/base64Captcha"
)

type Store interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Exists(key string) bool
}

type storeAdapter struct {
	store Store
}

func (s *storeAdapter) Set(id string, value string) error {
	key := "captcha:" + id
	return s.store.Set(key, value, 5*time.Minute)
}

func (s *storeAdapter) Get(id string, clear bool) string {
	key := "captcha:" + id
	
	if !s.store.Exists(key) {
		return ""
	}
	
	val, err := s.store.Get(key)
	if err != nil {
		return ""
	}
	
	if clear {
		s.store.Delete(key)
	}
	
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func (s *storeAdapter) Verify(id, answer string, clear bool) bool {
	return s.Get(id, clear) == answer
}

type Captcha struct {
	driver  base64Captcha.Driver
	store   base64Captcha.Store
}

func NewCaptcha(store Store) *Captcha {
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.8, 120)
	adapter := &storeAdapter{store: store}
	
	return &Captcha{
		driver: driver,
		store:  adapter,
	}
}

func (c *Captcha) Generate() (id string, b64s string, err error) {
	id, content, answer := c.driver.GenerateIdQuestionAnswer()
	item, err := c.driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}
	
	err = c.store.Set(id, answer)
	if err != nil {
		return "", "", err
	}
	
	b64s = item.EncodeB64string()
	return id, b64s, nil
}

func (c *Captcha) Verify(id string, answer string) bool {
	return c.store.Verify(id, answer, true)
}

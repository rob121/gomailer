package gomailer

type config map[string]interface{}

//these are the minium config values
var minimumvalid = []string{"Host", "Port", "User", "Password", "From", "Reply"}

func (c config) Get(key string) interface{} {

	if val, ok := c[key]; ok {
		return val
	}

	return nil
}

func (c config) validate() bool {

	for k, _ := range c {

		if _, ok := c[k]; !ok {

			return false
		}

	}

	return true

}

func (c config) Set(key string, value interface{}) {
	c[key] = value
}

func (c config) GetString(key string) string {
	return c.Get(key).(string)
}

func (c config) GetInt(key string) int {
	return c.Get(key).(int)
}

func (c config) GetInt64(key string) int64 {
	return c.Get(key).(int64)
}

func (c config) GetFloat64(key string) float64 {
	return c.Get(key).(float64)
}

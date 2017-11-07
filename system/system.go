package system

import (
	"net/http"
	"encoding/json"
	"reflect"
)

type Application struct {
}
type Controller struct {
}

func (application *Application) Route(controller interface{}, route string) interface{} {
	fn := func(w http.ResponseWriter, r *http.Request) {
		methodValue := reflect.ValueOf(controller).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func( w http.ResponseWriter, r *http.Request) ([]byte, error))
		result, err := method(w, r)
		w.Header().Set("Content-Type", "application/json")

		if (err != nil) {
			response := make(map[string]interface{})
			{
				response["message"] = "something went wrong"
				w.WriteHeader(http.StatusInternalServerError)
			}

			errResponse, _ := json.Marshal(response)
			w.Write(errResponse)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(result)
		}

	}
	return fn
}
/*

func PrepareApplicationContext() error {

	var err error

	_, err = InitMongoDbSession()
	*/
/*a:=c.C("ashu")
	a.Insert("aaaaaaaaaaaa")
*//*

	if err != nil {
		log.Println("Can't connect to mongodb database")
		return err
	} else {
		log.Println("Mongodb started successfully!")
	}
	return nil
}
*/


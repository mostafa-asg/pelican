package http

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mostafa-asg/pelican/store"
	"github.com/mostafa-asg/pelican/util"
	bytesUtil "github.com/mostafa-asg/pelican/util/bytes"
)

func GetHandler(kvStore *store.Store) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		value, found := kvStore.GetByteArray(key)
		if !found {
			w.Write([]byte("{}"))
		} else {
			w.Write([]byte("{ \"value\": \""))

			typ := r.Header.Get("type")
			switch typ {
			case "uint16":
				value := strconv.FormatUint(uint64(bytesUtil.ToUint16(value)), 10)
				w.Write([]byte(value))
			case "int16":
				value := strconv.FormatInt(int64(bytesUtil.ToInt16(value)), 10)
				w.Write([]byte(value))
			case "uint32":
				value := strconv.FormatUint(uint64(bytesUtil.ToUint32(value)), 10)
				w.Write([]byte(value))
			case "int32":
				value := strconv.FormatInt(int64(bytesUtil.ToInt32(value)), 10)
				w.Write([]byte(value))
			case "uint64":
				value := strconv.FormatUint(bytesUtil.ToUint64(value), 10)
				w.Write([]byte(value))
			case "int64":
				value := strconv.FormatInt(bytesUtil.ToInt64(value), 10)
				w.Write([]byte(value))
			case "bool":
				value := strconv.FormatBool(bytesUtil.ToBool(value[0]))
				w.Write([]byte(value))
			case "string":
				w.Write(value)
			default:
				var buf bytes.Buffer
				encoder := base64.NewEncoder(base64.StdEncoding, &buf)
				encoder.Write(value)
				encoder.Close()

				w.Write(buf.Bytes())
			}

			w.Write([]byte("\" }"))
		}
	}

}

func PutHandler(kvStore *store.Store) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		value, err := ioutil.ReadAll(r.Body)
		expire := r.Header.Get("expire")
		expireStrategy := r.Header.Get("strategy")
		if expireStrategy == "" {
			expireStrategy = "0"
		}

		if expire != "" {

			expireDuration, err := util.ToTimeDuration(expire)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			strategy, err := strconv.Atoi(expireStrategy)
			if err != nil {
				w.WriteHeader(400)
				return
			}

			kvStore.PutWithExpire(key, value, expireDuration, store.Strategy(strategy))
		} else {
			kvStore.Put(key, value)
		}

		if err != nil {
			w.Write([]byte("{ \"error\" : \""))
			w.Write([]byte(err.Error()))
			w.Write([]byte("\" }"))
		}
	}

}

func DelHandler(kvStore *store.Store) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		kvStore.Del(key)
	}

}

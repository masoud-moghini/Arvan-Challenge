package router

import (
	"arvan-challenge/application/rds"
	"arvan-challenge/application/router/dto"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/redis/go-redis/v9"
)

type (
	RequestHandlers struct {
		InMemoryServices rds.InMemoryServices
	}
)

func (rHandlers RequestHandlers) HandleIncommingRequest(w http.ResponseWriter, r *http.Request) {
	var remainedMinuteQuota int = rHandlers.InMemoryServices.GetAndDecreaseMinuteQuota(
		context.Background(),
		r.Header["user_id"][0],
		rHandlers.InMemoryServices.RedisClients.RedisClientForMinuteQuota,
	)

	var remainedMonthQuota int = rHandlers.InMemoryServices.GetAndDecreaseMonthlyQuota(
		context.Background(),
		r.Header["user_id"][0],
		rHandlers.InMemoryServices.RedisClients.RedisClientForMinuteQuota,
	)
	//bind to new variable
	data := &dto.UserRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, dto.ErrInvalidRequest(err))
		return
	}
	if remainedMinuteQuota > 0 && remainedMonthQuota > 0 {

		processData(
			context.Background(),
			data.DataToBeProcess,
			false,
			rHandlers.InMemoryServices.RedisClients.RedisClientForDataProcessed,
			w)
		return
	}
	processData(
		context.Background(),
		data.DataToBeProcess,
		true,
		rHandlers.InMemoryServices.RedisClients.RedisClientForDataProcessed,
		w)

}

func processData(
	ctx context.Context,
	dataToBeProcessed string,
	useOnlyCache bool,
	dataProcessingCacheDb *redis.Client,
	w http.ResponseWriter,
) {
	processedDataInString, err := dataProcessingCacheDb.Get(ctx, dataToBeProcessed).Result()
	if err == redis.Nil { //we could not find it in cache so process it if user has accessed in his quota
		if useOnlyCache { // user exceeded its quota so send it error message
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			result, _ := json.Marshal(dto.UserResponse{Data: "not authorized"})
			w.Write(result)
		} else { // user is within his quota so send processed message
			dataprocessed := strconv.Itoa(rand.Int()) //processing data as generating random value
			//put processed message in processing cachedb
			//set expire time to auto destroy unused data
			dataProcessingCacheDb.SetNX(ctx, dataToBeProcessed, dataprocessed, 10*time.Minute)
			w.WriteHeader(http.StatusOK)
			result, _ := json.Marshal(dto.UserResponse{Data: dataprocessed})
			w.Write(result)
		}
	} else if err == nil { // we could find it in cache so send it any way
		w.WriteHeader(http.StatusOK)
		result, _ := json.Marshal(dto.UserResponse{Data: processedDataInString})
		w.Write(result)
	}
}

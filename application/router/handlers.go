package router

import (
	"arvan-challenge/application/rds"
	"arvan-challenge/application/router/dto"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

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
	if remainedMinuteQuota > 0 && remainedMonthQuota > 0 {
		data := &dto.UserRquest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, dto.ErrInvalidRequest(err))
			return
		}
		processData(
			context.Background(),
			data.DataToBeProcess,
			false,
			rHandlers.InMemoryServices.RedisClients.RedisClientForDataProcessed,
			w, r)
		// article := data.Article
		// dbNewArticle(article)

		// render.Status(r, http.StatusCreated)
		// render.Render(w, r, NewArticleResponse(article))
	}
	w.Write([]byte("welcome"))
}

func processData(
	ctx context.Context,
	dataToBeProcessed string,
	useOnlyCache bool,
	dataProcessingCacheDb *redis.Client,
	w http.ResponseWriter,
	r *http.Request,
) {
	_, err := dataProcessingCacheDb.Get(ctx, dataToBeProcessed).Result()
	if err == redis.Nil {
		if useOnlyCache {
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			result, _ := json.Marshal(dto.UserResponse{Data: "not authorized"})
			w.Write(result)
		} else {
			dataprocessed := strconv.Itoa(rand.Int())
			dataProcessingCacheDb.Set(ctx, dataToBeProcessed, dataprocessed, 0)
			w.WriteHeader(http.StatusNonAuthoritativeInfo)
			result, _ := json.Marshal(dto.UserResponse{Data: dataprocessed})
			w.Write(result)
		}
	}

}

package Download

import (
	"net/http"

	"115Quick_server/internal/logic/Download"
	"115Quick_server/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClearCompletedTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Download.NewClearCompletedTasksLogic(r.Context(), svcCtx)
		resp, err := l.ClearCompletedTasks()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// middleware/recovery.go
package middleware

import (
    "log"
    "net/http"
    "runtime"

    "backend/dto"
    "backend/utils"
)

func Recovery(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("terjadi kesalahan: %v", err)
                for i := 2; ; i++ {
                    _, file, line, ok := runtime.Caller(i)
                    if !ok {
                        break
                    }
                    log.Printf("  %s:%d", file, line)
                }

                response := dto.ErrorResponse{
                    Success: false,
                    Status:  http.StatusInternalServerError,
                    Message: "Terjadi kesalahan internal pada server",
                }

                utils.JSONResponse(w, response.Status, response)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
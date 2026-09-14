package utils

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetStatusFromError(err error, w http.ResponseWriter) {
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				ResponseError(w, http.StatusBadRequest, st.Message())
				return
			case codes.Unauthenticated:
				ResponseError(w, http.StatusUnauthorized, st.Message())
				return
			}
		}

		ResponseError(w, http.StatusInternalServerError, st.Message())
		return
	}
}

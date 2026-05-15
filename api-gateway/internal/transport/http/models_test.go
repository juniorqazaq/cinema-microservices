package httpgateway

import (
	"errors"
	"net/http"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCErrorUsesGRPCGatewayHTTPMapping(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "unimplemented",
			err:        status.Error(codes.Unimplemented, "feature not available"),
			wantStatus: http.StatusNotImplemented,
			wantMsg:    "upstream service error",
		},
		{
			name:       "failed precondition stays conflict",
			err:        status.Error(codes.FailedPrecondition, "seat is already reserved"),
			wantStatus: http.StatusConflict,
			wantMsg:    "seat is already reserved",
		},
		{
			name:       "not found forwards grpc message",
			err:        status.Error(codes.NotFound, "movie not found"),
			wantStatus: http.StatusNotFound,
			wantMsg:    "movie not found",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotStatus, gotMsg := grpcError(tc.err)
			if gotStatus != tc.wantStatus {
				t.Fatalf("status = %d, want %d", gotStatus, tc.wantStatus)
			}
			if gotMsg != tc.wantMsg {
				t.Fatalf("message = %q, want %q", gotMsg, tc.wantMsg)
			}
		})
	}
}

func TestGRPCErrorHidesOpaqueFailures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "plain error",
			err:        errors.New("dial tcp: connection refused"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "internal grpc status",
			err:        status.Error(codes.Internal, "sql timeout"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotStatus, gotMsg := grpcError(tc.err)
			if gotStatus != tc.wantStatus {
				t.Fatalf("status = %d, want %d", gotStatus, tc.wantStatus)
			}
			if gotMsg != "upstream service error" {
				t.Fatalf("message = %q, want %q", gotMsg, "upstream service error")
			}
		})
	}
}

func TestGRPCErrorNil(t *testing.T) {
	t.Parallel()

	gotStatus, gotMsg := grpcError(nil)
	if gotStatus != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", gotStatus, http.StatusInternalServerError)
	}
	if gotMsg != "upstream service error" {
		t.Fatalf("message = %q, want %q", gotMsg, "upstream service error")
	}
}

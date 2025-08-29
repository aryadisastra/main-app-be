package unit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/aryadisastra/main-app-be/internal/dto"
	"github.com/aryadisastra/main-app-be/internal/router"
)

type envelope struct {
	Result  bool            `json:"result"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

const secret = "8"

func mkToken(userID, role string) string {
	claims := jwt.MapClaims{
		"sub":       userID,
		"role_code": role,
		"exp":       time.Now().Add(2 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := t.SignedString([]byte(secret))
	return signed
}

func setupDB(t *testing.T) *gorm.DB {
	dsn := "postgres://postgres:123@localhost:5432/logistic_db?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	db.Exec("DELETE FROM tb_r_shipments;")
	return db
}

func TestCreateTrackFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupDB(t)
	r := router.New(db, secret)

	userToken := mkToken("11111111-1111-1111-1111-111111111111", "user")

	crBody, _ := json.Marshal(dto.CreateShipmentRequest{
		SenderName: "A", SenderAddress: "Addr A",
		ReceiverName: "B", ReceiverAddress: "Addr B",
		ItemDescription: "Box",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/shipments", bytes.NewBuffer(crBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var env envelope
	_ = json.Unmarshal(w.Body.Bytes(), &env)
	require.True(t, env.Result)
	var created dto.ShipmentResponse
	_ = json.Unmarshal(env.Data, &created)
	require.NotEmpty(t, created.TrackingNumber)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/shipments/track/"+created.TrackingNumber, nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

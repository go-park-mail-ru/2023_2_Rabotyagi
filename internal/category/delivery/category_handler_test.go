package delivery

//
//import (
//	"errors"
//	"fmt"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/internal/category/mocks"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/middleware"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/models"
//	"github.com/go-park-mail-ru/2023_2_Rabotyagi/pkg/my_logger"
//	"github.com/stretchr/testify/mock"
//	"go.uber.org/mock/gomock"
//	"net/http"
//	"testing"
//)
//
//type MockLogger struct {
//	mock.Mock
//	RealLogger *my_logger.MyLogger
//}
//
//func (m *MockLogger) Info(message string) {
//	m.Called(message)
//}
//
//func (m *MockLogger) Error(message string, err error) {
//	m.Called(message, err)
//}
//
//func TestArtistDeliveryHTTP_Get(t *testing.T) {
//	// Init
//	type mockBehavior func(au *mock_category.MockICategoryService)
//
//	c := gomock.NewController(t)
//
//	mockService := mock_category.NewMockICategoryService(c)
//
//	h, _ := NewCategoryHandler(mockService)
//
//	// Routing
//	router := http.NewServeMux()
//	//router.Get("/api/artists/{artistID}/", h.Get)
//	router.Handle("/api/v1/category/get_full",
//		middleware.SetupCORS(h.GetFullCategories, "", ""))
//
//	// Test filling
//	const correctArtistID uint32 = 1
//	correctArtistIDPath := fmt.Sprint(correctArtistID)
//
//	expectedReturnArtist := models.Artist{
//		ID:        1,
//		Name:      "Oxxxymiron",
//		AvatarSrc: "/artists/avatars/oxxxymiron.png",
//	}
//
//	correctResponse := `{
//		"id": 1,
//		"name": "Oxxxymiron",
//		"isLiked": false,
//		"cover": "/artists/avatars/oxxxymiron.png"
//	}`
//
//	testTable := []struct {
//		name             string
//		artistIDPath     string
//		user             *models.User
//		mockBehavior     mockBehavior
//		expectedStatus   int
//		expectedResponse string
//	}{
//		{
//			name:         "Common",
//			artistIDPath: correctArtistIDPath,
//			user:         &correctUser,
//			mockBehavior: func(au *artistMocks.MockUsecase) {
//				au.EXPECT().GetByID(gomock.Any(), correctArtistID).Return(&expectedReturnArtist, nil)
//				au.EXPECT().IsLiked(gomock.Any(), correctArtistID, correctUser.ID).Return(false, nil)
//			},
//			expectedStatus:   http.StatusOK,
//			expectedResponse: correctResponse,
//		},
//		{
//			name:             "Incorrect ID In Path",
//			artistIDPath:     "0",
//			mockBehavior:     func(au *artistMocks.MockUsecase) {},
//			expectedStatus:   http.StatusBadRequest,
//			expectedResponse: commonTests.ErrorResponse(commonHttp.InvalidURLParameter),
//		},
//		{
//			name:         "No Artist To Get",
//			artistIDPath: correctArtistIDPath,
//			user:         &correctUser,
//			mockBehavior: func(au *artistMocks.MockUsecase) {
//				au.EXPECT().GetByID(gomock.Any(), correctArtistID).Return(nil, &models.NoSuchArtistError{})
//			},
//			expectedStatus:   http.StatusBadRequest,
//			expectedResponse: commonTests.ErrorResponse(artistNotFound),
//		},
//		{
//			name:         "Server Error",
//			artistIDPath: correctArtistIDPath,
//			user:         &correctUser,
//			mockBehavior: func(au *artistMocks.MockUsecase) {
//				au.EXPECT().GetByID(gomock.Any(), correctArtistID).Return(nil, errors.New(""))
//			},
//			expectedStatus:   http.StatusInternalServerError,
//			expectedResponse: commonTests.ErrorResponse(artistGetServerError),
//		},
//	}
//
//	for _, tc := range testTable {
//		t.Run(tc.name, func(t *testing.T) {
//			// Call mock
//			tc.mockBehavior(au)
//
//			commonTests.DeliveryTestGet(t, r, "/api/artists/"+tc.artistIDPath+"/",
//				tc.expectedStatus, tc.expectedResponse,
//				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
//		})
//	}
//}
